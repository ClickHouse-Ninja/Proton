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
		Name: "proton_server_processed_requests_total",
		Help: "The total number of processed requests",
	})
	cntBacklogCap = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "proton_server_backlog_cap",
		Help: "Server backlog capacity",
	})
	cntBacklogSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "proton_server_backlog_size",
		Help: "Server backlog size",
	})
	cntConcurrency = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "proton_server_concurrency_size",
		Help: "Number of server concurrent processes",
	})
	cntTablesRows = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "proton_server_table_rows",
		Help: "Number of rows in the table",
	}, []string{"table"})
	cntDictionaryRows = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "proton_server_dictionary_rows",
		Help: "Number of rows in the dictionary",
	}, []string{"column"})
)

func (server *server) metrics(address string) {
	if len(address) == 0 {
		return
	}
	cntBacklogCap.Set(float64(cap(server.reqBacklog)))
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
		cntBacklogSize.Set(float64(len(server.reqBacklog)))
		if err := server.dictionaryMetrics(); err != nil {
			log.Println("check dictionary size error: ", err)
		}
		if err := server.tablesInfoMetrics(); err != nil {
			log.Println("check table size error: ", err)
		}
		handler.ServeHTTP(w, r)
	}))
	log.Println(http.ListenAndServe(address, nil))
}

func (server *server) dictionaryMetrics() error {
	rows, err := server.sqlConnection.Query(dictionaryInfoSQL)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			column  string
			numRows int64
		)
		if err := rows.Scan(&column, &numRows); err != nil {
			return err
		}
		cntDictionaryRows.WithLabelValues(column).Set(float64(numRows))
	}
	return nil
}

func (server *server) tablesInfoMetrics() error {
	rows, err := server.sqlConnection.Query(tablesInfoSQL)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			table   string
			numRows int64
		)
		if err := rows.Scan(&table, &numRows); err != nil {
			return err
		}
		cntTablesRows.WithLabelValues(table).Set(float64(numRows))
	}
	return nil
}

const (
	dictionaryInfoSQL = `
	SELECT
		partition
		, SUM(rows) AS rows
	FROM system.parts
	WHERE database = 'proton' AND table = 'dictionary' AND active
	GROUP BY partition
	`
	tablesInfoSQL = `
	SELECT
		table
		, SUM(rows) AS rows
	FROM system.parts
	WHERE database = 'proton' AND table <> 'dictionary' AND active
	GROUP BY table
	`
)
