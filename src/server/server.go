package server

import (
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ClickHouse-Ninja/Proton/proto/pinba"
	"github.com/kshvakov/clickhouse"
	"github.com/kshvakov/clickhouse/lib/data"
)

func RunServer(options Options) error {
	conn, err := net.ListenPacket("udp", options.Address)
	if err != nil {
		return err
	}
	log.Printf("Proton server listen UDP [%s], Prometheus exporter [%s] concurrency: %d", options.Address, options.MetricsAddress, options.Concurrency)
	server := server{
		dsn:         options.DSN,
		backlog:     make(chan requestContainer, options.BacklogSize),
		connections: make(chan clickhouse.Clickhouse, options.Concurrency),
	}
	if err := server.prepare(); err != nil {
		return err
	}
	go server.metrics(options.MetricsAddress)
	for i := 0; i < options.Concurrency; i++ {
		go server.listen(conn)
		go server.background()
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	{
		log.Printf("stopped signal[%s]", <-signals)
	}
	return nil
}

type server struct {
	dsn         string
	block       *data.Block
	backlog     chan requestContainer
	connections chan clickhouse.Clickhouse
}

func (server *server) prepare() error {
	conn, err := server.connection()
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.Begin()
	if _, err = conn.Prepare(insertSQL); err != nil {
		return err
	}
	if server.block, err = conn.Block(); err != nil {
		return err
	}
	return nil
}

func (server *server) connection() (clickhouse.Clickhouse, error) {
	select {
	case conn := <-server.connections:
		return conn, nil
	default:
		return clickhouse.OpenDirect(server.dsn)
	}
}

func (server *server) releaseConn(conn clickhouse.Clickhouse, err error) error {
	if err == nil {
		select {
		case server.connections <- conn:
			return nil
		default:
		}
	}
	conn.Close()
	return err
}

func (server *server) listen(conn net.PacketConn) {
	var buffer [math.MaxUint16]byte
	for {
		var request pinba.Request
		if ln, _, err := conn.ReadFrom(buffer[:]); err == nil {
			if err := request.Unmarshal(buffer[:ln]); err == nil {
				container := requestContainer{
					Request:   request,
					timestamp: time.Now(),
				}
				select {
				case server.backlog <- container:
					opsBacklogSize.Add(1)
				default:
					log.Println("backlog is full")
				}
			}
		}
	}
}
