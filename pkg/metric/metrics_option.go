package metric

import (
	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"google.golang.org/genproto/googleapis/api/distribution"
)

type OptionBuilder struct {
}

func Int64Point(n int64) *monitoringpb.TypedValue {
	return &monitoringpb.TypedValue{
		Value: &monitoringpb.TypedValue_Int64Value{
			Int64Value: n,
		},
	}
}

func StringPoint(s string) *monitoringpb.TypedValue {
	return &monitoringpb.TypedValue{
		Value: &monitoringpb.TypedValue_StringValue{
			StringValue: s,
		},
	}
}

func DoublePoint(d float64) *monitoringpb.TypedValue {
	return &monitoringpb.TypedValue{
		Value: &monitoringpb.TypedValue_DoubleValue{
			DoubleValue: d,
		},
	}
}

func BoolPoint(b bool) *monitoringpb.TypedValue {
	return &monitoringpb.TypedValue{
		Value: &monitoringpb.TypedValue_BoolValue{
			BoolValue: b,
		},
	}
}

func DistributionPoint(
	count int64,
	mean float64,
	sumOfSquaredDeviation float64,
	bucketOptions *distribution.Distribution_BucketOptions,
	bucketCounts []int64,
) *monitoringpb.TypedValue {
	return &monitoringpb.TypedValue{
		Value: &monitoringpb.TypedValue_DistributionValue{
			DistributionValue: &distribution.Distribution{
				Count:                 count,
				Mean:                  mean,
				SumOfSquaredDeviation: sumOfSquaredDeviation,
				BucketOptions:         bucketOptions,
				BucketCounts:          bucketCounts,
			},
		},
	}
}
