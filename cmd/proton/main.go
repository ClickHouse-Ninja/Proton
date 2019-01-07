package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ClickHouse-Ninja/Proton/src/server"
)

var options server.Options

func init() {
	flag.StringVar(&options.DSN, "dsn", "native://127.0.0.1:9000", "ClickHouse DSN")
	flag.StringVar(&options.Address, "addr", "127.0.0.1:30002", "listen address")
	flag.IntVar(&options.BacklogSize, "backlog", 100000, "backlog size")
}
func main() {
	flag.Usage = func() {
		fmt.Println("Proton - high performance Pinba storage server.")
		flag.PrintDefaults()
	}
	flag.Parse()
	if err := server.RunServer(options); err != nil {
		log.Fatal(err)
	}
}
