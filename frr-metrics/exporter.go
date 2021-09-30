package main

import (
	"flag"
	"fmt"
	stdlog "log"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/web"

	"go.universe.tf/metallb/frr-metrics/collector"
	"go.universe.tf/metallb/internal/logging"
)

var (
	metricsPort   = flag.String("metrics-port", ":7473", "Address to listen on for web interface.")
	metricsPath   = flag.String("metrics-path", "/metrics", "Path under which to expose metrics.")
	tlsConfigPath = flag.String("tls-config-path", "", "[EXPERIMENTAL] Path to config yaml file that can enable TLS or authentication.")
)

func metricsHandler(logger log.Logger) http.Handler {
	BGPCollector := collector.NewBGP(logger)

	registry := prometheus.NewRegistry()
	registry.MustRegister(BGPCollector)

	gatherers := prometheus.Gatherers{
		prometheus.DefaultGatherer,
		registry,
	}

	handlerOpts := promhttp.HandlerOpts{
		ErrorLog:      stdlog.New(log.NewStdlibAdapter(level.Error(logger)), "", 0),
		ErrorHandling: promhttp.ContinueOnError,
		Registry:      registry,
	}

	return promhttp.HandlerFor(gatherers, handlerOpts)
}

func main() {
	logger, err := logging.Init("error")
	if err != nil {
		fmt.Printf("failed to initialize logging: %s\n", err)
		os.Exit(1)
	}

	http.Handle(*metricsPath, metricsHandler(logger))
	srv := &http.Server{Addr: *metricsPort}
	level.Info(logger).Log("msg", "Starting exporter", "metricsPath", metricsPath, "port", metricsPort)

	if err := web.ListenAndServe(srv, *tlsConfigPath, logger); err != nil {
		level.Error(logger).Log("error", err)
		os.Exit(1)
	}
}
