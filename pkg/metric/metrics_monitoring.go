package metric

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/timestamppb"

	monitoring "cloud.google.com/go/monitoring/apiv3/v2"
)

func NewMonitoringMetric(projectID string, credentialsJSON []byte, opts ...OptionBuilder) (Monitoring, error) {

	// parse projectID from credentialsJSON if not provided
	if projectID == "" {
		var creds struct {
			ProjectID string `json:"project_id"`
		}
		if err := json.Unmarshal(credentialsJSON, &creds); err != nil ||
			creds.ProjectID == "" {
			return nil, fmt.Errorf("project_id from credentials JSON empty or"+
				" failed to read: %v", err)
		}
		projectID = creds.ProjectID
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Create the monitoring client
	client, err := monitoring.NewMetricClient(ctx, option.WithCredentialsJSON(credentialsJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create monitoring client: %v", err)
	}

	return &metrics{
		name:                  defaultServiceName,
		labels:                make(map[string]string),
		projectID:             projectID,
		client:                client,
		minimumSamplingPeriod: minimumSamplingPeriod,
	}, nil
}

func (m *metrics) NewTable(name string, labels map[string]string) Table {
	children := &metrics{
		name:                  name,
		labels:                labels,
		projectID:             m.projectID,
		client:                m.client,
		minimumSamplingPeriod: minimumSamplingPeriod,
	}
	m.pendingFinalizers = append(m.pendingFinalizers, children.Close)
	return children
}

func (m *metrics) Close() (err error) {
	for _, f := range m.pendingFinalizers {
		if ne := f(); ne != nil {
			err = errors.Join(err, ne)
		}
	}
	if ne := m.sync(); ne != nil {
		return errors.Join(err, ne)
	}
	return err
}

func (m *metrics) SendTimeSeries(ctx context.Context, timeSeries []*monitoringpb.TimeSeries) error {
	if len(timeSeries) == 0 {
		return nil
	}
	return m.sendToGCP(ctx, timeSeries)
}

// SendMetrics sends a batch of metrics to GCP Monitoring
// metrics map keys are metric names, values are metric points
//
// metrics is a map where keys are metric names (e.g., "read_operations",  "read_errors", "write_operations","write_errors")
// and values are the corresponding metric points (int64, simple is 1).
func (m *metrics) SendMetrics(ctx context.Context, method string, metrics map[string]*monitoringpb.TypedValue) error {
	// build time series from metrics map
	timeSeries := make([]*monitoringpb.TimeSeries, 0, len(metrics))
	now := time.Now()
	timestamp := timestamppb.New(now)

	// create labels with default and method specific labels
	labels := m.getMetricLabels(m.labels)
	labels["method"] = method

	// build time series from metrics map
	for key, value := range metrics {
		timeSeries = append(timeSeries, createTimeSeries(m.name, key, labels, timestamp, value))
	}
	return m.SendTimeSeries(ctx, timeSeries)
}
