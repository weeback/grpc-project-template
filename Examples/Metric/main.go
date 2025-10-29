package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/monitoring/apiv3/v2/monitoringpb"
	"github.com/weeback/grpc-project-template/pkg/metric"
)

var (
	temp = struct {
		GCPProjectID       string `yaml:"gcp_project_id"`
		GCPCredentialsPath string `yaml:"gcp_credentials_path"`
		GCPCredentialsJSON []byte `yaml:"-"`
	}{}
)

func readMetricClient(path string) (err error) {
	if path != "" {
		temp.GCPCredentialsJSON, err = os.ReadFile(path)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("gcp-credentials-path must be set")
}

func init() {
	flag.Func("gcp-project-id", "GCP Project ID", func(id string) error {
		if id == "" {
			return fmt.Errorf("gcp-project-id cannot be empty")
		}
		temp.GCPProjectID = id
		return nil
	})
	flag.Func("gcp-credentials-path", "GCP Credentials Path", func(path string) error {
		temp.GCPCredentialsPath = path
		return readMetricClient(path)
	})
	flag.Parse()

	if temp.GCPProjectID == "" {
		fmt.Fprintf(os.Stderr, "project-id must be set, "+
			"Try '%s --help' for more information.\n",
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	if len(temp.GCPCredentialsJSON) == 0 {
		fmt.Fprintf(os.Stderr, "gcp-credentials-path must be set and valid, "+
			"Try '%s --help' for more information.\n",
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}

}

func main() {

	type MetricExample struct {
		Method  string
		Metrics map[string]*monitoringpb.TypedValue
	}

	// Initialize the metrics server
	// Enable monitoring metrics with client if initialized, then log any error
	mm, err := metric.NewMonitoringMetric(temp.GCPProjectID, temp.GCPCredentialsJSON)
	if err != nil {
		log.Fatalf("Failed to create Metric Client: %v", err)
	}
	defer func() {
		if err := mm.Close(); err != nil {
			fmt.Printf("[WARN ]Failed to close Metric Client: %v\n", err)
		}
		//
		fmt.Println("wait a bit to ensure all metrics are sent before exit ...")
		time.Sleep(30 * time.Second)

	}()
	fmt.Println("Metric Client initialized successfully")

	tb := mm.NewTable("mongodb", map[string]string{
		"app": "MetricExample",
		"env": "dev",
	})

	oneValue := metric.Int64Point(1)
	zeroValue := metric.Int64Point(0)

	// Simulate input
	metrics := []MetricExample{
		{
			Method:  "Read",
			Metrics: map[string]*monitoringpb.TypedValue{"read_operations": oneValue, "read_errors": zeroValue},
		},
		{
			Method:  "Write",
			Metrics: map[string]*monitoringpb.TypedValue{"write_operations": oneValue, "write_errors": oneValue},
		},
		{
			Method:  "Read",
			Metrics: map[string]*monitoringpb.TypedValue{"read_operations": oneValue, "read_errors": zeroValue},
		},
		{
			Method:  "Write",
			Metrics: map[string]*monitoringpb.TypedValue{"write_operations": oneValue, "write_errors": oneValue},
		},
		{
			Method:  "Read",
			Metrics: map[string]*monitoringpb.TypedValue{"read_operations": oneValue, "read_errors": oneValue},
		},
		{
			Method:  "Write",
			Metrics: map[string]*monitoringpb.TypedValue{"write_operations": oneValue, "write_errors": oneValue},
		},
		{
			Method:  "Read",
			Metrics: map[string]*monitoringpb.TypedValue{"read_operations": oneValue, "read_errors": oneValue},
		},
		{
			Method:  "Write",
			Metrics: map[string]*monitoringpb.TypedValue{"write_operations": oneValue, "write_errors": oneValue},
		},
		{
			Method:  "Read",
			Metrics: map[string]*monitoringpb.TypedValue{"read_operations": oneValue, "read_errors": zeroValue},
		},
		{
			Method:  "Write",
			Metrics: map[string]*monitoringpb.TypedValue{"write_operations": oneValue, "write_errors": oneValue},
		},
		{
			Method:  "Read",
			Metrics: map[string]*monitoringpb.TypedValue{"read_operations": oneValue, "read_errors": oneValue},
		},
		{
			Method:  "Write",
			Metrics: map[string]*monitoringpb.TypedValue{"write_operations": oneValue, "write_errors": oneValue},
		},
		{
			Method:  "Read",
			Metrics: map[string]*monitoringpb.TypedValue{"read_operations": oneValue, "read_errors": zeroValue},
		},
		{
			Method:  "Write",
			Metrics: map[string]*monitoringpb.TypedValue{"write_operations": oneValue, "write_errors": oneValue},
		},
		{
			Method:  "Read",
			Metrics: map[string]*monitoringpb.TypedValue{"read_operations": oneValue, "read_errors": zeroValue},
		},
		{
			Method:  "Write",
			Metrics: map[string]*monitoringpb.TypedValue{"write_operations": oneValue, "write_errors": oneValue},
		},
		{
			Method:  "Read",
			Metrics: map[string]*monitoringpb.TypedValue{"read_operations": oneValue, "read_errors": zeroValue},
		},
		{
			Method:  "Write",
			Metrics: map[string]*monitoringpb.TypedValue{"write_operations": oneValue, "write_errors": zeroValue},
		},
		{
			Method:  "Read",
			Metrics: map[string]*monitoringpb.TypedValue{"read_operations": oneValue, "read_errors": zeroValue},
		},
		{
			Method:  "Write",
			Metrics: map[string]*monitoringpb.TypedValue{"write_operations": oneValue, "write_errors": zeroValue},
		},
	}

	// Duplicate the metrics slice to simulate more data
	for range 8 {
		metrics = append(metrics, metrics...)
	}

	// Simulate sending some metrics
	for i, m := range metrics {

		if err := tb.SendMetrics(context.TODO(), m.Method, m.Metrics); err != nil {
			fmt.Printf("[ERROR] Failed to send metrics for method %s (i=%d): %v\n---\n", m.Method, i, err)
		} else {
			fmt.Printf("[INFO ] Successfully sent metrics for method %s (i=%d)\n---\n", m.Method, i)
		}
	}
}
