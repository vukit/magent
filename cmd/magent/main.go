package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/vukit/magent/internal/collector"
	"github.com/vukit/magent/internal/config"
	"github.com/vukit/magent/internal/logger"
	"github.com/vukit/magent/internal/metric"
	"github.com/vukit/magent/internal/sensor"
)

func main() {
	filename := flag.String("c", "configs/local.json", "specifies the path to the configuration file")
	flag.Parse()

	mConfig := &config.Config{}

	err := mConfig.Read(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mLogger := logger.NewLogger(mConfig, nil)

	metrics := getMetrics(mConfig, mLogger)

	sendMetrics(metrics, mConfig, mLogger)
}

func getMetrics(mConfig *config.Config, mLogger *logger.Logger) (metrics metric.Metrics) {
	var wgMetrics, wgWorkers sync.WaitGroup

	metrics = make(metric.Metrics, 0)
	channel := make(chan metric.Metrics)

	wgMetrics.Add(1)

	go func(channel <-chan metric.Metrics) {
		for {
			data, isOpen := <-channel
			if !isOpen {
				break
			}

			metrics = append(metrics, data...)
		}
		wgMetrics.Done()
	}(channel)

	for _, s := range mConfig.Sensors {
		if !s.Enable {
			continue
		}

		switch s.Name {
		case "memory":
			s := s
			worker := func(channel chan<- metric.Metrics) {
				defer wgWorkers.Done()

				result, err := (&sensor.Memory{}).Metrics(s.Metrics)
				if err != nil {
					mLogger.Warning(err)

					return
				}
				channel <- result
			}

			wgWorkers.Add(1)

			go worker(channel)
		case "cpu":
			s := s
			worker := func(channel chan<- metric.Metrics) {
				defer wgWorkers.Done()

				result, err := (&sensor.CPU{}).Metrics(s.Metrics)
				if err != nil {
					mLogger.Warning(err)

					return
				}
				channel <- result
			}

			wgWorkers.Add(1)

			go worker(channel)
		case "disk":
			s := s
			worker := func(channel chan<- metric.Metrics) {
				defer wgWorkers.Done()

				result, err := (&sensor.Disk{}).Metrics(s.Metrics, s.Devices)
				if err != nil {
					mLogger.Warning(err)

					return
				}
				channel <- result
			}

			wgWorkers.Add(1)

			go worker(channel)
		}
	}

	wgWorkers.Wait()

	close(channel)

	wgMetrics.Wait()

	return metrics
}

func sendMetrics(metrics metric.Metrics, mConfig *config.Config, mLogger *logger.Logger) {
	var wgWorkers sync.WaitGroup

	for _, c := range mConfig.Collectors {
		if !c.Enable {
			continue
		}

		switch c.Name {
		case "console":
			worker := func() {
				defer wgWorkers.Done()

				err := (&collector.Console{Config: &mConfig.Common, Logger: mLogger}).Send(metrics)
				if err != nil {
					mLogger.Warning(err)
				}
			}

			wgWorkers.Add(1)

			go worker()
		case "yandex":
			worker := func() {
				defer wgWorkers.Done()

				err := (&collector.Yandex{Config: &mConfig.Common, Logger: mLogger}).Send(metrics, c.Parameters)
				if err != nil {
					mLogger.Warning(err)
				}
			}

			wgWorkers.Add(1)

			go worker()
		default:
			continue
		}
	}

	wgWorkers.Wait()
}
