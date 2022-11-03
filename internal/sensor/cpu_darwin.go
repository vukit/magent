//go:build darwin

package sensor

import (
	"bytes"
	"math"
	"os/exec"
	"strconv"

	"github.com/vukit/magent/internal/metric"
)

type CPU struct{}

func (cpu *CPU) Metrics(metrics []string) (result metric.Metrics, err error) {
	result = make(metric.Metrics, 0)

	for _, metricName := range metrics {
		switch metricName {
		case "usage":
			value, err := getUsageValue()
			if err != nil {
				return nil, err
			}

			result = append(result, metric.Metric{Name: "CPU Usage", DValue: value, TValue: "DGAUGE"})
		default:
			continue
		}
	}

	return result, nil
}

func getUsageValue() (value float64, err error) {
	command := exec.Command("top", "-l", "1")
	filter := exec.Command("grep", "-E", "^CPU")

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

	result := bytes.Fields(output)[6]
	result = result[:len(result)-1]

	value, err = strconv.ParseFloat(string(result), 64)
	if err != nil {
		return 0, err
	}

	return math.Ceil(100.0*(100-value)) / 100, nil
}
