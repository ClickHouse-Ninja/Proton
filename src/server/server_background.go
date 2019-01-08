package server

import (
	"log"
	"time"

	"github.com/kshvakov/clickhouse"
	"github.com/kshvakov/clickhouse/lib/data"
)

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
				block.WriteUInt32(5, *request.RequestCount)
				block.WriteFloat32(6, *request.RequestTime)
				block.WriteUInt32(7, *request.DocumentSize)
				block.WriteUInt32(8, *request.MemoryPeak)
				block.WriteUInt32(9, *request.MemoryFootprint)
				block.WriteFloat32(10, *request.RuUtime)
				block.WriteFloat32(11, *request.RuStime)
				tagName, tagValue := request.tags()
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
				block.WriteDateTime(18, request.timestamp)
			case <-tick:
				break loop
			}
		}
		if err := server.write(block); err != nil {
			log.Println("write ", err)
		}
	}
}

func (server *server) write(block *data.Block) error {
	if block.NumRows == 0 {
		return nil
	}
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
