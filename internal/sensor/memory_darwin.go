//go:build darwin

package sensor

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"

	"github.com/vukit/magent/internal/metric"
)

const (
	totalMetric = "total"
	freeMetric  = "free"
	usedMetric  = "used"
)

type Memory struct {
	metrics []string
	values  map[string]int64
}

func (memory *Memory) Metrics(memoryMetrics []string) (result []metric.Metric, err error) {
	memory.metrics = memoryMetrics
	memory.values = make(map[string]int64)

	for _, metric := range memory.metrics {
		if metric == totalMetric {
			value, err := getTotalValue()
			if err != nil {
				return nil, err
			}

			memory.values[totalMetric] = value
		}

		if metric == usedMetric {
			value, err := getUsedValue()
			if err != nil {
				return nil, err
			}

			memory.values[usedMetric] = value
		}

		if metric == freeMetric {
			memory.values[freeMetric] = 0
		}
	}

	if err := memory.checkMetrics(); err != nil {
		return nil, err
	}

	return memory.result(), nil
}

func (memory *Memory) checkMetrics() error {
	if _, ok := memory.values["free"]; ok {
		if _, ok := memory.values["total"]; !ok {
			return errors.New("sensors/memory_darwin: memory total value not found")
		}

		if _, ok := memory.values["used"]; !ok {
			return errors.New("sensors/memory_darwin: used memory value not found")
		}
	}

	return nil
}

func (memory *Memory) result() metric.Metrics {
	result := make([]metric.Metric, 0)

	for _, metricName := range memory.metrics {
		switch metricName {
		case totalMetric:
			result = append(result, metric.Metric{Name: "Memory total", IValue: memory.values[totalMetric], TValue: "IGAUGE"})
		case freeMetric:
			result = append(result, metric.Metric{
				Name:   "Free memory",
				IValue: memory.values[totalMetric] - memory.values[usedMetric],
				TValue: "IGAUGE",
			})
		case usedMetric:
			result = append(result, metric.Metric{Name: "Used memory", IValue: memory.values[usedMetric], TValue: "IGAUGE"})
		}
	}

	return result
}

func getUsedValue() (value int64, err error) {
	command := exec.Command("top", "-l", "1")
	filter := exec.Command("grep", "-E", "^Phys")

	pipe, err := command.StdoutPipe()
	if err != nil {
		return 0, err
	}
	defer pipe.Close()

	filter.Stdin = pipe

	err = command.Start()
	if err != nil {
		return 0, err
	}

	output, err := filter.Output()
	if err != nil {
		return 0, nil
	}

	result := bytes.Fields(output)[1]
	result = result[:len(result)-1]

	value, err = strconv.ParseInt(string(result), 10, 64)
	if err != nil {
		return 0, err
	}

	return 1024 * 1024 * value, nil
}

func getTotalValue() (value int64, err error) {
	output, err := exec.Command("sysctl", "hw.memsize").Output()
	if err != nil {
		return 0, err
	}

	result := bytes.Fields(output)[1]

	value, err = strconv.ParseInt(string(result), 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}
