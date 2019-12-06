package telemetry

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "hello_world"
	subsystem = "greeter"
)

// Telemetry exposes the available metrics
type Telemetry interface {
	Registry() *prometheus.Registry
	RequestDuration(path string, start time.Duration)
}

// New creates a new instance of the telemetry provider
func New(registry *prometheus.Registry) (Telemetry, error) {
	p := &provider{
		registry: registry,
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "duration_milliseconds",
				Help:      "Duration of http request_by_path",
				Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000, 10000},
			}, []string{"path"}),
	}

	for _, c := range p.metrics() {
		err := p.registry.Register(c)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

type provider struct {
	registry        *prometheus.Registry
	requestDuration *prometheus.HistogramVec
}

func (p *provider) metrics() []prometheus.Collector {
	return []prometheus.Collector{p.requestDuration}
}

// Registry returns the registry
func (p *provider) Registry() *prometheus.Registry {
	return p.registry
}

func (p *provider) RequestDuration(path string, duration time.Duration) {
	p.requestDuration.WithLabelValues(path).Observe(1000 * duration.Seconds())
}
