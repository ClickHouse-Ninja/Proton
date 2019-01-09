package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ClickHouse-Ninja/Proton/src/server"
)

var options server.Options

var (
	BuildDate            string
	GitBranch, GitCommit string
)

func init() {
	flag.StringVar(&options.DSN, "dsn", "native://127.0.0.1:9000", "ClickHouse DSN")
	flag.StringVar(&options.Address, "addr", ":30002", "listen address")
	flag.StringVar(&options.MetricsAddress, "metrics_addr", ":2112", "address on which to expose metrics")
	flag.IntVar(&options.BacklogSize, "backlog", 100000, "backlog size")
	flag.IntVar(&options.Concurrency, "concurrency", 2, "number of the background processes")
}

func main() {
	flag.Usage = func() {
		fmt.Println("NAME:")
		fmt.Println("  Proton - high performance Pinba storage server.")
		fmt.Println("VERSION:")
		fmt.Printf("  0.1 rev[%s] %s (%s UTC).\n", GitCommit, GitBranch, BuildDate)
		fmt.Println("USAGE:")
		flag.PrintDefaults()
	}
	flag.Parse()
	if err := server.RunServer(options); err != nil {
		log.Fatal(err)
	}
}
