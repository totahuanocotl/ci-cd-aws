package helloworld

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

func handlerReady(greeter *Greeter) (string, http.Handler) {
	return timedHandler(greeter, "/ready", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
}

func handlerHealthz(greeter *Greeter) (string, http.Handler) {
	return timedHandler(greeter, "/healthz", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
}

func handlerHello(greeter *Greeter) (string, http.Handler) {
	return timedHandler(greeter, "/hello", func(rw http.ResponseWriter, req *http.Request) {
		greeting, err := greeter.Hello(getName(req))
		handleResponse(greeting, err, rw)
	})
}

func handlerGoodbye(greeter *Greeter) (string, http.Handler) {
	return timedHandler(greeter, "/goodbye", func(rw http.ResponseWriter, req *http.Request) {
		farewell, err := greeter.Goodbye(getName(req))
		handleResponse(farewell, err, rw)
	})
}

func handlerMetrics(greeter *Greeter) (string, http.Handler) {
	return "/metrics", promhttp.InstrumentMetricHandler(
		greeter.metrics.Registry(), promhttp.HandlerFor(greeter.metrics.Registry(), promhttp.HandlerOpts{}),
	)
}

func timedHandler(greeter *Greeter, mapping string, handler func(rw http.ResponseWriter, req *http.Request)) (string, http.Handler) {
	return mapping, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		start := time.Now()
		handler(rw, r)
		greeter.metrics.RequestDuration(mapping, time.Since(start))
	})
}

func getName(req *http.Request) string {
	keys, ok := req.URL.Query()["name"]
	if ok {
		return keys[0]
	}
	return ""
}

func handleResponse(message string, err error, rw http.ResponseWriter) {
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.WriteHeader(http.StatusOK)
	_, err = fmt.Fprintln(rw, message)
	if err != nil {
		log.Errorf("failed to send message %q back: %v", message, err)
	}

}
