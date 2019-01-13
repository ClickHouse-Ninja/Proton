package server

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	opsReqProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "proton_server_requests_total",
		Help: "The total number of processed requests",
	})
	opsBacklogCap = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "proton_server_backlog_cap",
		Help: "Server backlog capacity",
	})
	opsBacklogSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "proton_server_backlog_size",
		Help: "Server backlog size",
	})
	opsConcurrency = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "proton_server_concurrency_size",
		Help: "Server concurrency",
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
