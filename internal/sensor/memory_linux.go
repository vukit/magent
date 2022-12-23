//go:build linux

package sensor

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/vukit/magent/internal/metric"
)

type Memory struct {
	metrics []string
	values  map[string]int64
}

func (memory *Memory) Metrics(memoryMetrics []string) (result metric.Metrics, err error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, fmt.Errorf("sensors/linux/memory: %v", err)
	}

	memory.metrics = memoryMetrics
	memory.values = make(map[string]int64)

	re := make([]*regexp.Regexp, 4)
	re[0] = regexp.MustCompile(`^MemTotal:\s*(\d+)\s*kB$`)
	re[1] = regexp.MustCompile(`^MemFree:\s+(\d+)\s*kB$`)
	re[2] = regexp.MustCompile(`^Cached:\s+(\d+)\s*kB$`)
	re[3] = regexp.MustCompile(`^Buffers:\s+(\d+)\s*kB$`)

	lines := bytes.Split(data, []byte{'\n'})
	for _, line := range lines {
		res := re[0].FindAllSubmatch(line, -1)
		if len(res) > 0 {
			value, err := strconv.ParseInt(string(res[0][1]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("sensors/linux/memory: memory total value: %v", err)
			}

			memory.values["total"] = 1024 * value
		}

		res = re[1].FindAllSubmatch(line, -1)
		if len(res) > 0 {
			value, err := strconv.ParseInt(string(res[0][1]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("sensors/linux/memory: free memory value: %v", err)
			}

			memory.values["free"] = 1024 * value
		}

		res = re[2].FindAllSubmatch(line, -1)
		if len(res) > 0 {
			value, err := strconv.ParseInt(string(res[0][1]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("sensors/linux/memory: cached memory value: %v", err)
			}

			memory.values["cached"] = 1024 * value
		}

		res = re[3].FindAllSubmatch(line, -1)
		if len(res) > 0 {
			value, err := strconv.ParseInt(string(res[0][1]), 10, 64)
			if err != nil {
				return nil, fmt.Errorf("sensors/linux/memory: buffers memory value: %v", err)
			}

			memory.values["buffers"] = 1024 * value
		}
	}

	if err := memory.checkMetrics(); err != nil {
		return nil, err
	}

	return memory.result(), nil
}

func (memory *Memory) checkMetrics() error {
	if _, ok := memory.values["total"]; !ok {
		return errors.New("sensors/memory_linux: memory total value not found")
	}

	if _, ok := memory.values["free"]; !ok {
		return errors.New("sensors/memory_linux: free memory value not found")
	}

	if _, ok := memory.values["cached"]; !ok {
		return errors.New("sensors/memory_linux: cached memory value not found")
	}

	if _, ok := memory.values["buffers"]; !ok {
		return errors.New("sensors/memory_linux: buffers memory value not found")
	}

	return nil
}

func (memory *Memory) result() metric.Metrics {
	result := make([]metric.Metric, 0)

	for _, metricName := range memory.metrics {
		switch metricName {
		case "total":
			result = append(result, metric.Metric{Name: "Memory total", IValue: memory.values["total"], TValue: "IGAUGE"})
		case "free":
			result = append(result, metric.Metric{Name: "Free memory", IValue: memory.values["free"], TValue: "IGAUGE"})
		case "cached":
			result = append(result, metric.Metric{Name: "Cached memory", IValue: memory.values["cached"], TValue: "IGAUGE"})
		case "buffers":
			result = append(result, metric.Metric{Name: "Buffers memory", IValue: memory.values["buffers"], TValue: "IGAUGE"})
		case "used":
			result = append(result, metric.Metric{Name: "Used memory", IValue: memory.values["total"] - memory.values["free"], TValue: "IGAUGE"})
		case "app":
			result = append(result, metric.Metric{
				Name:   "Application memory",
				IValue: memory.values["total"] - memory.values["free"] - memory.values["cached"] - memory.values["buffers"],
				TValue: "IGAUGE",
			})
		}
	}

	return result
}
