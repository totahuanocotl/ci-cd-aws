package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/totahuanocotl/hello-world/pkg/helloworld"
)

var (
	// Provided at compiled time
	version   string
	buildTime string
)

var (
	debugLogging bool
	config       = &helloworld.Config{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hello-world",
	Short: "service to greet and farewell",
	Run:   validate,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Version = fmt.Sprintf("%s (%s)", version, buildTime)
	rootCmd.PersistentFlags().BoolVarP(&debugLogging, "debug", "X", false, "enable debug logging.")
	rootCmd.PersistentFlags().Int32VarP(&config.Port, "port", "p", 8080, "port to listen for greetings and farewells")
	rootCmd.PersistentFlags().DurationVar(&config.ShutdownTimeout, "shutdown-timeout", 5*time.Second, "maximum time to wait for a graceful shutdown")

}

func initConfig() {
	if debugLogging {
		log.SetLevel(log.DebugLevel)
	}
	log.Infof("Config:%#v", config)
}

func validate(_ *cobra.Command, _ []string) {
	if err := config.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	greeter, err := helloworld.New(config, prometheus.NewRegistry())
	if err != nil {
		log.Fatalf("Failed to create greeter service: %v", err)
	}

	addSignalHandler(greeter)
	if err := greeter.Start(); err != nil {
		log.Fatalf("Failed to start greeter service: %v", err)
	}
	log.Info("Stopped greeter service")
}

func addSignalHandler(greeter *helloworld.Greeter) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		defer close(signals)
		for sig := range signals {
			log.Infof("Signalled %v, shutting down gracefully", sig)
			err := greeter.Stop()
			if err != nil {
				log.Errorf("Error while shutting down: %v", err)
			}
		}
	}()
}
