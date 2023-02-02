package sink

import (
	"context"
	"time"
)

type Metric struct {
	data map[string][]int64
	avg  map[string]int64
	ctx  context.Context
}

func NewMetric(ctx context.Context) *Metric {
	return &Metric{
		data: make(map[string][]int64),
		avg:  make(map[string]int64),
		ctx:  ctx,
	}
}

func (m *Metric) Add(key string, metric uint64) {
	// Calculate latency until now
	currentTime := time.Now()
	latency := currentTime.UnixMilli() - int64(metric)
	m.data[key] = append(m.data[key], latency)
}

func (m *Metric) Aggregate() {
	for host, metrics := range m.data {
		// Dont aggregate if we dont have metrics
		if len(metrics) > 0 {
			sum := int64(0)
			for _, metric := range metrics {
				sum += metric
			}
			avg := sum / int64(len(metrics))
			m.avg[host] = avg
			//log.Printf("%s -> avg %dms size: %d", host, avg, len(metrics))
			// Erase values
			m.data[host] = nil
		}

	}
}
