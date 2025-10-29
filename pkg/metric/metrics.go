package metric

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
)

const (
	defaultResource    = "global"
	defaultServiceName = "default"

	defaultCustomPath = "custom.googleapis.com"
)

var (
	minimumSamplingPeriod = 10 * time.Millisecond
)

type metrics struct {
	projectID string
	name      string
	labels    map[string]string

	client                *monitoring.MetricClient
	minimumSamplingPeriod time.Duration

	// internal logic fields can be added here
	forceFreeze       bool
	pendingFinalizers []func() error

	mu       sync.RWMutex // to protect concurrent access
	pools    []*monitoringpb.TimeSeries
	lastTime time.Time
}

// getProjectID returns the GCP project ID
func (m *metrics) getProjectID() string {
	return m.projectID
}

// getResourceType returns the GCP resource type
// default is "global" if not set
func (m *metrics) getResourceType() string {
	return defaultResource
}

// getServiceName returns the service name for metrics
// default is "default" if not set
func (m *metrics) getServiceName() string {
	if m.name == "" {
		return defaultServiceName
	}
	return m.name
}

func (m *metrics) generateId() string {
	name := m.getServiceName()
	return fmt.Sprintf("%s-%s", name, uuid.NewString())
}

func (m *metrics) getMetricLabels(custom map[string]string) map[string]string {
	// default labels
	labels := map[string]string{
		"id":           m.generateId(), // unique ID for each metric, required
		"project_id":   m.getProjectID(),
		"service_name": m.getServiceName(),
		"resource":     m.getResourceType(),
		"date":         time.Now().Format("2006-01-02"),
	}
	// add custom labels
	for key, val := range custom {
		switch strings.ToLower(key) {
		case "id", "project_id", "service_name", "resource", "date":
			// skip reserved labels
		default:
			labels[key] = val
		}
	}
	return labels
}

// sendToGCP sends MongoDB metrics to Google Cloud Monitoring
// This function creates custom metrics for monitoring MongoDB operations
// including read/write operations and their error rates.
//
// Parameters:
//   - ctx: context.Context for the request
//   - timeSeries: slice of TimeSeries objects representing the metrics to be sent
//
// Returns:
//   - error: error object if any issues occur during metric sending
//
// Example:
//
//	err := metric.sendToGCP(ctx, timeSeries)
//	if err != nil {
//	   // handle error
//	}
//
// Note: Ensure the GCP client has proper authentication and permissions
func (m *metrics) sendToGCP(ctx context.Context, timeSeries []*monitoringpb.TimeSeries) error {
	// validate project ID
	if m.projectID == "" {
		return fmt.Errorf("projectID is required to send metrics to GCP")
	}
	// validate client
	if m.client == nil {
		return fmt.Errorf("GCP MetricClient is not initialized")
	}
	// validate time-series data
	for _, ts := range timeSeries {
		if ts == nil {
			return fmt.Errorf("object TimeSeries cannot be nil")
		}
		if ts.Metric == nil {
			return fmt.Errorf("field Metric in TimeSeries cannot be nil")
		}
		if ts.Resource == nil {
			return fmt.Errorf("field Resource in TimeSeries cannot be nil")
		}
		if len(ts.Points) != 1 {
			return fmt.Errorf("each TimeSeries must have exactly one data point")
		}
	}

	// check minimum sampling period
	if !m.forceFreeze &&
		time.Since(m.lastTime) <= m.minimumSamplingPeriod {

		m.mu.Lock()
		// cache the time series data
		m.pools = append(m.pools, timeSeries...)
		// skip sending to GCP to respect minimum sampling period
		m.mu.Unlock()

		// sleep for a while to avoid busy looping
		time.Sleep(5 * time.Millisecond)
		return nil
	}

	var (
		c = make(chan error, 1)
	)
	go func() {

		// respect minimum sampling period
		if delta := time.Since(m.lastTime); delta <= m.minimumSamplingPeriod {
			fmt.Printf("[DEBUG] waiting for remaining time: %s\n", delta.String())
			<-time.After(delta)
		}

		maxLen := 200 // GCP max time-series per request
		if n := len(timeSeries) + len(m.pools); n < maxLen {
			maxLen = n
		}
		// aggregate metrics
		aggregateMetrics := make([]*monitoringpb.TimeSeries, 0, maxLen)
		aggregateMetrics = append(aggregateMetrics, timeSeries...)
		// prepare to send cached metrics
		m.mu.Lock()
		// combine cached metrics with current time series
		if len(m.pools) > 0 {
			aggregateMetrics = append(aggregateMetrics, m.pools...)
			m.pools = make([]*monitoringpb.TimeSeries, 0) // clear cache
		}
		if n := len(aggregateMetrics) - maxLen; n > 0 {
			// truncate to maxLen
			aggregateMetrics = aggregateMetrics[:maxLen]
			// keep the remaining in cache
			m.pools = aggregateMetrics[maxLen:]
		}
		m.mu.Unlock()

		// validate time-series data before sending
		if len(aggregateMetrics) == 0 {
			// nothing to send
			c <- nil
			return
		}

		// create a context with timeout for sending metrics
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer func() {
			m.lastTime = time.Now()
			cancel()
		}()
		// send all metrics
		fmt.Printf("[DEBUG] Sending %d time-series to GCP\n", len(aggregateMetrics))
		if err := m.client.CreateTimeSeries(ctx, &monitoringpb.CreateTimeSeriesRequest{
			Name:       fmt.Sprintf("projects/%s", m.projectID),
			TimeSeries: aggregateMetrics,
		}); err != nil {
			fmt.Printf("[DEBUG] failed to send cached metrics to GCP: %v\n", err)
			c <- err
			return
		} else {
			fmt.Printf("[DEBUG] sent %d metrics to GCP successfully\n", len(aggregateMetrics))
		}
		c <- nil
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		// goroutine finished
		if err != nil {
			fmt.Printf("[DEBUG] metrics sent to GCP, err=%v\n", err)
			return err
		}
		fmt.Printf("[DEBUG] metrics sent to GCP successfully\n")
		return nil
	case <-time.After(999 * time.Millisecond):
		fmt.Println("[DEBUG] sending metrics to longtime, this will continue to run in background")
		return nil
	}
}

func (m *metrics) sync() error {
	fmt.Println("[DEBUG] freezing cached ...")
	// Set last-time is in the past to force sending cached metrics
	m.forceFreeze = true
	// Trigger sending cached metrics
	return m.sendToGCP(context.Background(), []*monitoringpb.TimeSeries{})
}

// createTimeSeries constructs a TimeSeries object for GCP Monitoring
// based on the provided parameters.
//
// This is templated to create custom metrics for monitoring MongoDB operations.
// It sets up the metric type, labels, and data points for the time series.
//
// May be points is one of type Int64Value, DoubleValue, etc., depending on the metric being recorded.
// Example: &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_Int64Value{Int64Value: 1}}
func createTimeSeries(name string, path string, labels map[string]string,
	timestamp *timestamppb.Timestamp, point *monitoringpb.TypedValue,
) *monitoringpb.TimeSeries {

	if !strings.HasPrefix(path, defaultCustomPath) {
		// use provided metricType
		if path == "" {
			path = fmt.Sprintf("%s/%s/%s", defaultCustomPath, name, "default_metric") // path=custom.googleapis.com/mongodb/default_metric
		} else {
			path = fmt.Sprintf("%s/%s/%s", defaultCustomPath, name, path) // path=custom.googleapis.com/mongodb/read_operations
		}
	}

	// ensure labels contain "id"
	if _, exist := labels["id"]; !exist {
		labels["id"] = uuid.New().String()
	}
	// ensure point is not nil
	if point == nil {
		point = &monitoringpb.TypedValue{
			Value: &monitoringpb.TypedValue_Int64Value{
				Int64Value: 0,
			},
		}
	}

	return &monitoringpb.TimeSeries{
		Metric: &metricpb.Metric{
			Type:   path,
			Labels: labels,
		},
		Resource: &monitoredrespb.MonitoredResource{
			Type: defaultResource,
		},
		Points: []*monitoringpb.Point{
			{
				Interval: &monitoringpb.TimeInterval{
					EndTime: timestamp,
				},
				Value: point,
			},
		},
	}
}
