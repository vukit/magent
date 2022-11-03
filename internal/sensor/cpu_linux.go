//go:build linux

package sensor

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/vukit/magent/internal/metric"
)

type CPU struct {
	source string
}

func (cpu *CPU) Metrics(metrics []string) (result metric.Metrics, err error) {
	var value float64

	if cpu.source == "" {
		cpu.source = "/proc/stat"
	}

	startValues, err := cpu.getProcStatValues()

	if len(startValues) == 0 {
		return nil, err
	}

	time.Sleep(time.Second)

	endValues, _ := cpu.getProcStatValues()

	deltaIdle := endValues[3] - startValues[3]

	deltaTotal := 0.0
	for i, startValue := range startValues {
		deltaTotal += endValues[i] - startValue
	}

	if deltaTotal != 0 {
		value = math.Ceil(10000.0*(1-deltaIdle/deltaTotal)) / 100
	}

	result = make(metric.Metrics, 0)

	for _, metricName := range metrics {
		switch metricName {
		case "usage":
			result = append(result, metric.Metric{Name: "CPU Usage", DValue: value, TValue: "DGAUGE"})
		default:
			continue
		}
	}

	return result, nil
}

func (cpu *CPU) getProcStatValues() (values []float64, err error) {
	data, err := os.ReadFile(cpu.source)
	if err != nil {
		return nil, fmt.Errorf("sensors/linux/cpu: %v", err)
	}

	firstLine := bytes.Split(data, []byte{'\n'})[0]
	for _, byteValues := range bytes.Fields(firstLine)[1:] {
		value, err := strconv.ParseFloat(string(byteValues), 64)
		if err != nil {
			return nil, fmt.Errorf("sensors/linux/cpu: %v", err)
		}

		values = append(values, value)
	}

	return
}
