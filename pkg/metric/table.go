package metric

import (
	"context"

	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
)

type Monitoring interface {
	NewTable(name string, labels map[string]string) Table
	Close() (err error)
}

type Table interface {
	NewTable(name string, labels map[string]string) Table
	Close() (err error)
	SendTimeSeries(ctx context.Context, timeSeries []*monitoringpb.TimeSeries) error
	SendMetrics(ctx context.Context, method string, metrics map[string]*monitoringpb.TypedValue) error
}
