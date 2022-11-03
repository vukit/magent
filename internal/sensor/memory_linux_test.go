//go:build linux

package sensor

import (
	"testing"

	"github.com/vukit/magent/internal/metric"
)

func TestMemory(t *testing.T) {
	var (
		metrics metric.Metrics
		err     error
	)

	metrics, err = (&Memory{}).Metrics([]string{"total", "free", "cached", "buffers", "app", "used"})
	if err != nil {
		t.Error(err)
	}

	if metrics[0].Name != "Memory total" {
		t.Errorf("wrong metric name, want 'Memory total', got '%s'", metrics[0].Name)
	}
	if metrics[0].IValue < 0 {
		t.Errorf("wrong metric value, want value >= 0, got %d", metrics[0].IValue)
	}
	if metrics[0].TValue != "IGAUGE" {
		t.Errorf("wrong metric type, want type = 'IGAUGE', got '%s", metrics[0].TValue)
	}

	if metrics[1].Name != "Free memory" {
		t.Errorf("wrong metric name, want 'Free memory', got '%s'", metrics[1].Name)
	}
	if metrics[1].IValue < 0 {
		t.Errorf("wrong metric value, want value >= 0, got %d", metrics[1].IValue)
	}
	if metrics[1].TValue != "IGAUGE" {
		t.Errorf("wrong metric type, want type = 'IGAUGE', got '%s", metrics[1].TValue)
	}

	if metrics[2].Name != "Cached memory" {
		t.Errorf("wrong metric name, want 'Cached memory', got '%s'", metrics[2].Name)
	}
	if metrics[2].IValue < 0 {
		t.Errorf("wrong metric value, want value >= 0, got %d", metrics[2].IValue)
	}
	if metrics[2].TValue != "IGAUGE" {
		t.Errorf("wrong metric type, want type = 'IGAUGE', got '%s", metrics[2].TValue)
	}

	if metrics[3].Name != "Buffers memory" {
		t.Errorf("wrong metric name, want 'Buffers memory', got '%s'", metrics[3].Name)
	}
	if metrics[3].IValue < 0 {
		t.Errorf("wrong metric value, want value >= 0, got %d", metrics[3].IValue)
	}
	if metrics[3].TValue != "IGAUGE" {
		t.Errorf("wrong metric type, want type = 'IGAUGE', got '%s", metrics[3].TValue)
	}

	if metrics[4].Name != "Application memory" {
		t.Errorf("wrong metric name, want 'Buffers memory', got '%s'", metrics[4].Name)
	}
	if metrics[4].IValue < 0 {
		t.Errorf("wrong metric value, want value >= 0, got %d", metrics[4].IValue)
	}
	if metrics[4].TValue != "IGAUGE" {
		t.Errorf("wrong metric type, want type = 'IGAUGE', got '%s", metrics[4].TValue)
	}

	if metrics[5].Name != "Used memory" {
		t.Errorf("wrong metric name, want 'Buffers memory', got '%s'", metrics[5].Name)
	}
	if metrics[5].IValue < 0 {
		t.Errorf("wrong metric value, want value >= 0, got %d", metrics[5].IValue)
	}
	if metrics[5].TValue != "IGAUGE" {
		t.Errorf("wrong metric type, want type = 'IGAUGE', got '%s", metrics[5].TValue)
	}
}
