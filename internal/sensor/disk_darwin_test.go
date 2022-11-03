//go:build darwin

package sensor

import (
	"testing"

	"github.com/vukit/magent/internal/metric"
)

func TestDisk(t *testing.T) {
	var (
		metrics metric.Metrics
		err     error
	)

	metrics, err = (&Disk{}).Metrics([]string{"size", "avail"}, []string{})
	if err != nil {
		t.Error(err)
	}

	if metrics[0].Name != "Disk Total Size /" {
		t.Errorf("wrong metric name, want 'Disk Total Size /', got '%s'", metrics[0].Name)
	}

	if metrics[0].IValue < 0 {
		t.Errorf("wrong metric value, want value >= 0, got %d", metrics[0].IValue)
	}

	if metrics[0].TValue != "IGAUGE" {
		t.Errorf("wrong metric type, want type = 'IGAUGE', got '%s", metrics[0].TValue)
	}

	if metrics[1].Name != "Disk Available Size /" {
		t.Errorf("wrong metric name, want 'Disk Available Size /', got '%s'", metrics[1].Name)
	}

	if metrics[1].IValue < 0 {
		t.Errorf("wrong metric value, want value >= 0, got %d", metrics[1].IValue)
	}

	if metrics[1].TValue != "IGAUGE" {
		t.Errorf("wrong metric type, want type = 'IGAUGE', got '%s", metrics[1].TValue)
	}

	metrics, err = (&Disk{}).Metrics([]string{"unknown"}, []string{})
	if err != nil {
		t.Error(err)
	}

	if len(metrics) != 0 {
		t.Errorf("non-empty result, got %v", metrics)
	}
}
