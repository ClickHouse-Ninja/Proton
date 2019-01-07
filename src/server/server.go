package server

import (
	"fmt"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ClickHouse-Ninja/Proton/proto/pinba"
	"github.com/kshvakov/clickhouse"
	"github.com/kshvakov/clickhouse/lib/data"
)

type Options struct {
	DSN         string
	Address     string
	BacklogSize int
}

func RunServer(options Options) error {
	conn, err := net.ListenPacket("udp", options.Address)
	if err != nil {
		return err
	}
	log.Println("listen ", options.Address)
	server := server{
		dsn:         options.DSN,
		backlog:     make(chan pinba.Request, options.BacklogSize),
		connections: make(chan clickhouse.Clickhouse, runtime.NumCPU()),
	}
	if err := server.prepare(); err != nil {
		return err
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		go server.listen(conn)
		go server.background()
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("stopped signal[%s]", <-signals)
	return nil
}

type server struct {
	dsn         string
	block       *data.Block
	backlog     chan pinba.Request
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

func (server *server) background() {
	for tick := time.Tick(time.Second); ; {
		block := server.block.Copy()
	loop:
		for block.Reserve(); ; {
			select {
			case request := <-server.backlog:
				block.NumRows++
				block.WriteString(0, *request.Hostname)
				block.WriteString(1, *request.Schema)
				block.WriteInt16(2, int16(*request.Status))
				block.WriteString(3, *request.ServerName)
				block.WriteString(4, *request.ScriptName)
				block.WriteInt64(5, int64(*request.RequestCount))
				block.WriteFloat32(6, *request.RequestTime)
				block.WriteUInt32(7, *request.DocumentSize)
				block.WriteUInt32(8, *request.MemoryPeak)
				block.WriteUInt32(9, *request.MemoryFootprint)
				block.WriteFloat32(10, *request.RuUtime)
				block.WriteFloat32(11, *request.RuStime)
				tagName, tagValue := tags(request)
				{
					block.WriteArray(12, clickhouse.Array(tagName))
					block.WriteArray(13, clickhouse.Array(tagValue))
				}
				// timer
				block.WriteArray(14, clickhouse.Array(request.TimerHitCount))
				block.WriteArray(15, clickhouse.Array(request.TimerValue))
				block.WriteArray(16, clickhouse.Array(request.TimerRuStime))
				block.WriteArray(17, clickhouse.Array(request.TimerRuUtime))
				// block.WriteArray(18, clickhouse.Array(Array(T))  TagsName @todo add support of Array(Array(T)) to the driver
				// block.WriteArray(19, clickhouse.Array(Array(T))) TagsValue
				block.WriteDateTime(18, time.Now()) // @todo
			case <-tick:
				break loop
			}
		}
		if err := server.write(block); err != nil {
			log.Println("write ", err)
		}
	}
}

func tags(req pinba.Request) ([]string, []string) {
	var (
		name  []string
		value []string
	)
	if len(req.TagName) <= len(req.Dictionary) && len(req.TagValue) <= len(req.Dictionary) {
		for _, k := range req.TagName {
			name = append(name, req.Dictionary[int(k)])
		}
		for _, k := range req.TagValue {
			value = append(value, req.Dictionary[int(k)])
		}
	}
	return name, value
}

func (server *server) write(block *data.Block) error {
	conn, err := server.connection()
	if err != nil {
		return err
	}
	conn.Begin()
	if _, err = conn.Prepare(insertSQL); err != nil {
		return server.releaseConn(conn, err)
	}
	if err = conn.WriteBlock(block); err != nil {
		return server.releaseConn(conn, err)
	}
	return server.releaseConn(conn, conn.Commit())
}

func (server *server) listen(conn net.PacketConn) {
	var buffer [math.MaxUint16]byte
	for {
		var request pinba.Request
		if ln, _, err := conn.ReadFrom(buffer[:]); err == nil {
			if err := request.Unmarshal(buffer[:ln]); err == nil {
				select {
				case server.backlog <- request:
				default:
					fmt.Println("backlog is full")
				}
			}
		}
	}
}

const (
	insertSQL = `
	INSERT INTO proton.request (
		Hostname
		, Schema
		, Status
		, ServerName
		, ScriptName
		, RequestCount
		, RequestTime
		, DocumentSize
		, MemoryPeak
		, MemoryFootprint
		, Utime
		, Stime
		, Tags.Name
		, Tags.Value
		, Timers.HitCount
		, Timers.Value
		, Timers.Utime
		, Timers.Stime
		/*, Timers.TagsName
		, Tiers.TagsValue*/
		, Timestamp
	) VALUES (
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?,
		?
	)
	`
)
