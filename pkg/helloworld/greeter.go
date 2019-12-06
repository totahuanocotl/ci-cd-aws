package helloworld

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/totahuanocotl/hello-world/pkg/telemetry"
	"github.com/totahuanocotl/hello-world/pkg/web"

	log "github.com/sirupsen/logrus"
)

// Config holds configuration for the service
type Config struct {
	Port            int32
	ShutdownTimeout time.Duration
	Registry        *prometheus.Registry
}

// Validate will return an error if the validation is invalid
func (c *Config) Validate() error {
	if c.Port <= 0 {
		return fmt.Errorf("port [%d] should be > 0", c.Port)
	}
	return nil
}

// Greeter has the setup for the service
type Greeter struct {
	server  web.Server
	config  *Config
	metrics telemetry.Telemetry
}

// New creates a new Greeter instance based on the configuration.
func New(config *Config, registry *prometheus.Registry) (*Greeter, error) {
	server := web.NewServer(config.Port)
	metrics, err := telemetry.New(registry)
	if err != nil {
		return nil, fmt.Errorf("failed to create telemetry: %v", err)
	}
	greeter := &Greeter{
		server:  server,
		config:  config,
		metrics: metrics,
	}

	server.RegisterEndpoint(handlerReady(greeter))
	server.RegisterEndpoint(handlerHealthz(greeter))
	server.RegisterEndpoint(handlerHello(greeter))
	server.RegisterEndpoint(handlerGoodbye(greeter))
	server.RegisterEndpoint(handlerMetrics(greeter))
	return greeter, nil
}

// Start starts the Greeter server
func (v *Greeter) Start() error {
	log.Info("Starting greeter")
	return v.server.Start()
}

// Stop stops the Greeter server
func (v *Greeter) Stop() error {
	log.Info("Starting greeter")
	return v.server.Stop(v.config.ShutdownTimeout)
}

// Hello will return a greeting to the caller with the name.
func (v *Greeter) Hello(name string) (string, error) {
	return fmt.Sprintf("Hello, %s!", sanitizedName(name)), nil
}

// Goodbye will return a farewell to the caller with the name.
func (v *Greeter) Goodbye(name string) (string, error) {
	return fmt.Sprintf("Farewell, %s!", sanitizedName(name)), nil
}

func sanitizedName(name string) string {
	if name == "" {
		return "Stranger"
	}
	return name
}
