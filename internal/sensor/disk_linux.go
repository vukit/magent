//go:build linux

package sensor

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/vukit/magent/internal/metric"
)

type Disk struct{}

func (disk *Disk) Metrics(metrics, devices []string) (result metric.Metrics, err error) {
	if len(devices) == 0 {
		devices = append(devices, "/")
	}

	result = make(metric.Metrics, 0)

	for _, device := range devices {
		out, err := exec.Command("/bin/df", "--output=size,avail", device).Output()
		if err != nil {
			return nil, fmt.Errorf("sensors/linux/disk: %v", err)
		}

		data := bytes.Fields(out)

		size, err := strconv.ParseInt(string(data[2]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("sensors/linux/disk: %v", err)
		}

		avail, err := strconv.ParseInt(string(data[3]), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("sensors/linux/disk: %v", err)
		}

		for _, metricName := range metrics {
			switch metricName {
			case "size":
				result = append(result, metric.Metric{Name: "Disk Total Size " + device, IValue: 1024 * size, TValue: "IGAUGE"})
			case "avail":
				result = append(result, metric.Metric{Name: "Disk Available Size " + device, IValue: 1024 * avail, TValue: "IGAUGE"})
			}
		}
	}

	return result, nil
}
