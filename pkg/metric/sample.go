package metric

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// putRead sends read operation metrics to GCP Monitoring
func (m *metrics) putRead(method string, issue error) error {
	// validate
	if m.projectID == "" {
		return fmt.Errorf("projectID is required to send metrics to GCP")
	}
	if m.client == nil {
		return fmt.Errorf("GCP MetricClient is not initialized")
	}

	now := time.Now()
	timestamp := timestamppb.New(now)

	readPoint := Int64Point(1)
	readErrorPoint := Int64Point(0)
	// Check if there was an error
	if issue != nil && issue != mongo.ErrNoDocuments {
		readErrorPoint = Int64Point(1)
	}

	labels := map[string]string{
		"id":           m.generateId(), // unique ID for each metric, required
		"service_name": m.name,
		"method":       string(method),
		"date":         now.Format("2006-01-02"),
	}

	// Create time series data for different metrics
	timeSeries := []*monitoringpb.TimeSeries{
		createTimeSeries(m.name, "read_operations", labels, timestamp, readPoint),
		createTimeSeries(m.name, "read_errors", labels, timestamp, readErrorPoint),
		createTimeSeries(m.name, "write_operations", labels, timestamp, Int64Point(0)), // on method read, put write point is 0 or not set this
		createTimeSeries(m.name, "write_errors", labels, timestamp, Int64Point(0)),     // on method read, put write point is 0 or not set this
	}

	// Send metrics to GCP if we have data
	if len(timeSeries) > 0 {
		return m.sendToGCP(context.Background(), timeSeries)
	}

	return nil
}

// putWrite sends write operation metrics to GCP Monitoring
func (m *metrics) putWrite(method string, issue error) error {
	// validate
	if m.projectID == "" {
		return fmt.Errorf("projectID is required to send metrics to GCP")
	}
	if m.client == nil {
		return fmt.Errorf("GCP MetricClient is not initialized")
	}

	now := time.Now()
	timestamp := timestamppb.New(now)

	writePoint := Int64Point(1)
	writeErrorPoint := Int64Point(0)
	//
	if issue != nil && issue != mongo.ErrNoDocuments {
		writeErrorPoint = Int64Point(1)
	}

	labels := map[string]string{
		"id":           m.generateId(), // unique ID for each metric, required
		"service_name": m.name,
		"method":       string(method),
		"date":         now.Format("2006-01-02"),
	}

	// Create time series data for different metrics
	timeSeries := []*monitoringpb.TimeSeries{
		createTimeSeries(m.name, "read_operations", labels, timestamp, Int64Point(0)), // on method write, put read point is 0 or not set this
		createTimeSeries(m.name, "read_errors", labels, timestamp, Int64Point(0)),     // on method write, put read point is 0 or not set this
		createTimeSeries(m.name, "write_operations", labels, timestamp, writePoint),
		createTimeSeries(m.name, "write_errors", labels, timestamp, writeErrorPoint),
	}

	// Send metrics to GCP if we have data
	if len(timeSeries) > 0 {
		return m.sendToGCP(context.Background(), timeSeries)
	}

	return nil
}
