package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	version          = "v2.0.0"
	bgColor          = "white"
	healthyStatus    = true
	mutexHealthState = &sync.RWMutex{}
	up               = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "up",
		Help: "Service up with version information",
		ConstLabels: map[string]string{
			"version": version,
		},
	})
	httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Count of all HTTP requests",
	}, []string{"code", "method"})
)

func main() {
	bind := ""
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.StringVar(&bind, "bind", ":8080", "The socket to bind to.")
	flagset.Parse(os.Args[1:])

	bgColorVal, bgColorExists := os.LookupEnv("BACKGROUND_COLOR")
	if bgColorExists {
		bgColor = bgColorVal
	}

	// prometheus registries
	r := prometheus.NewRegistry()
	r.MustRegister(httpRequestsTotal)
	r.MustRegister(up)
	r.MustRegister(prometheus.NewProcessCollector(
		//uses current process' PID from os.Getpid()
		prometheus.ProcessCollectorOpts{}))

	// http handlers
	http.Handle("/", promhttp.InstrumentHandlerCounter(httpRequestsTotal, helloHandler()))
	http.Handle("/api", promhttp.InstrumentHandlerCounter(httpRequestsTotal, apiHandler()))
	http.Handle("/err", promhttp.InstrumentHandlerCounter(httpRequestsTotal, errorHandler()))
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	http.HandleFunc("/health", healthyHandler)
	http.HandleFunc("/down", downHandler)

	log.Println(fmt.Sprintf("HTTP server listening on %v", bind))
	up.Set(1)
	log.Fatal(http.ListenAndServe(bind, nil))
}

func healthyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		mutexHealthState.RLock()
		defer mutexHealthState.RUnlock()
		if healthyStatus {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Healthy"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Unhealthy"))
		}
	}
}

func downHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		mutexHealthState.Lock()
		defer mutexHealthState.Unlock()
		healthyStatus = false
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func errorHandler() http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Unexpected error occurred")
		w.WriteHeader(http.StatusInternalServerError)
	})
	return handler
}

func helloHandler() http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.WriteHeader(http.StatusOK)
			hostname, _ := os.Hostname()
			w.Write([]byte(fmt.Sprintf("<html><head><title>Hello World</title></head><body style=\"background-color: %v\"><h1>Hello from %v version %v</h1></body></html>", bgColor, hostname, version)))
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	return handler
}

// Hello contains version and hostname
type Hello struct {
	Version  string
	Hostname string
}

func apiHandler() http.HandlerFunc {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			hostname, _ := os.Hostname()
			hello := Hello{Version: version, Hostname: hostname}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(hello)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	})
	return handler
}
