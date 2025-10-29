package metric

import (
	"context"

	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
)

type NoopTable struct{}

func (n *NoopTable) NewTable(name string, labels map[string]string) Table {
	// no-op implementation
	return &NoopTable{}
}

func (n *NoopTable) Close() (err error) {
	// no-op implementation
	return nil
}

func (n *NoopTable) SendTimeSeries(ctx context.Context, timeSeries []*monitoringpb.TimeSeries) error {
	// no-op implementation
	return nil
}

func (n *NoopTable) SendMetrics(ctx context.Context, method string, metrics map[string]*monitoringpb.TypedValue) error {
	// no-op implementation
	return nil
}
