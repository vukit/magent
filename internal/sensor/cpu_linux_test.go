//go:build linux

package sensor

import (
	"os"
	"testing"

	"github.com/vukit/magent/internal/metric"
)

func TestCPU(t *testing.T) {
	var (
		metrics metric.Metrics
		err     error
	)

	metrics, err = (&CPU{}).Metrics([]string{"usage"})
	if err != nil {
		t.Error(err)
	}

	if metrics[0].Name != "CPU Usage" {
		t.Errorf("wrong metric name, want 'CPU Usage', got '%s'", metrics[0].Name)
	}

	if metrics[0].DValue < 0 || metrics[0].DValue > 100 {
		t.Errorf("wrong metric value, want value >= 0 and value <= 100, got %.2f", metrics[0].DValue)
	}

	if metrics[0].TValue != "DGAUGE" {
		t.Errorf("wrong metric type, want type = 'DGAUGE', got '%s", metrics[0].TValue)
	}

	metrics, err = (&CPU{}).Metrics([]string{"unknown"})
	if err != nil {
		t.Error(err)
	}

	if len(metrics) != 0 {
		t.Errorf("non-empty result, got %v", metrics)
	}

	metrics, err = (&CPU{source: "/proc/nonexistent"}).Metrics([]string{"usage"})
	if err == nil {
		t.Error("reading /proc/nonexistent gave a non-empty result")
	}

	if len(metrics) != 0 {
		t.Error("reading /proc/nonexistent gave a non-empty result")
	}

	os.WriteFile("/tmp/cpudata.txt", []byte("One Two Three"), os.ModePerm)
	values, err := (&CPU{source: "/tmp/cpudata.txt"}).getProcStatValues()
	if err == nil {
		t.Errorf("non-empty result, got %v, %s", values, err)
	}
	os.Remove("/tmp/cpudata.txt")
}
