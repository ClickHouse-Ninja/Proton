package server

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "proton_processed_events_total",
		Help: "The total number of processed events",
	})
	opsBacklogCap = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "proton_backlog_cap",
		Help: "The backlog capacity",
	})
	opsBacklogSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "proton_backlog_size",
		Help: "The backlog size",
	})
)

func (server *server) metrics(address string) {
	if len(address) == 0 {
		return
	}
	opsBacklogCap.Set(float64(cap(server.backlog)))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Proton Exporter</title></head>
			<body>
			<h1>Proton Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})
	handler := promhttp.Handler()
	http.Handle("/metrics", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opsBacklogSize.Set(float64(len(server.backlog)))
		handler.ServeHTTP(w, r)
	}))
	log.Println(http.ListenAndServe(address, nil))
}
