package server

import (
	"log"
	"time"
)

func (server *server) background() {
	for tick := time.Tick(time.Second); ; {
		block := server.block.Copy()
	loop:
		for block.Reserve(); ; {
			select {
			case request := <-server.reqBacklog:
				block.NumRows++
				block.WriteString(0, request.GetHostname())
				block.WriteString(1, request.GetSchema())
				block.WriteInt16(2, int16(request.GetStatus()))
				block.WriteString(3, request.GetServerName())
				block.WriteString(4, request.GetScriptName())
				block.WriteUInt32(5, request.GetRequestCount())
				block.WriteFloat32(6, request.GetRequestTime())
				block.WriteUInt32(7, request.GetDocumentSize())
				block.WriteUInt32(8, request.GetMemoryPeak())
				block.WriteUInt32(9, request.GetMemoryFootprint())
				block.WriteFloat32(10, request.GetRuUtime())
				block.WriteFloat32(11, request.GetRuStime())
				tagName, tagValue := request.tags()
				{
					block.WriteArray(12, tagName)
					block.WriteArray(13, tagValue)
				}
				// timer
				block.WriteArray(14, request.TimerHitCount)
				block.WriteArray(15, request.GetTimerValue())
				block.WriteArray(16, request.GetTimerRuStime())
				block.WriteArray(17, request.GetTimerRuUtime())
				timerTagName, timerTagValue := request.timerTags()
				{
					block.WriteArray(18, timerTagName)
					block.WriteArray(19, timerTagValue)
				}
				block.WriteUInt32(20, request.timestamp)
			case <-tick:
				break loop
			}
		}
		opsReqProcessed.Add(float64(block.NumRows))
		if err := server.writeBlock(insertIntoRequestsSQL, block); err != nil {
			log.Println("request write error: ", err)
		}
	}
}

const (
	insertIntoRequestsSQL = `
	INSERT INTO proton.requests (
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
		, Timers.TagsName
		, Timers.TagsValue
		, Timestamp
	) VALUES (
		?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?, ?, ?, ?, ?, ?,
		?, ?, ?
	)
	`
)
