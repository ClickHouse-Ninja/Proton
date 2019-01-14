package server

import (
	"github.com/ClickHouse-Ninja/Proton/proto/pinba"
)

type Options struct {
	DSN            string
	Address        string
	MetricsAddress string
	BacklogSize    int
	Concurrency    int
}

type request struct {
	pinba.Request
	timestamp uint32
}

func (req *request) Hostname() string {
	if req.Request.Hostname != nil {
		return *req.Request.Hostname
	}
	return ""
}
func (req *request) Schema() string {
	if req.Request.Schema != nil {
		return *req.Request.Schema
	}
	return ""
}
func (req *request) Status() int16 {
	if req.Request.Status != nil {
		return int16(*req.Request.Status)
	}
	return 0
}
func (req *request) ServerName() string {
	if req.Request.ServerName != nil {
		return *req.Request.ServerName
	}
	return ""
}

func (req *request) ScriptName() string {
	if req.Request.ScriptName != nil {
		return *req.Request.ScriptName
	}
	return ""
}

func (req *request) RequestCount() uint32 {
	if req.Request.RequestCount != nil {
		return *req.Request.RequestCount
	}
	return 0
}

func (req *request) RequestTime() float32 {
	if req.Request.RequestTime != nil {
		return *req.Request.RequestTime
	}
	return 0
}
func (req *request) DocumentSize() uint32 {
	if req.Request.DocumentSize != nil {
		return *req.Request.DocumentSize
	}
	return 0
}

func (req *request) MemoryPeak() uint32 {
	if req.Request.MemoryPeak != nil {
		return *req.Request.MemoryPeak
	}
	return 0
}
func (req *request) MemoryFootprint() uint32 {
	if req.Request.MemoryFootprint != nil {
		return *req.Request.MemoryFootprint
	}
	return 0
}

func (req *request) RuUtime() float32 {
	if req.Request.RuUtime != nil {
		return *req.Request.RuUtime
	}
	return 0
}
func (req *request) RuStime() float32 {
	if req.Request.RuStime != nil {
		return *req.Request.RuStime
	}
	return 0
}

func (req *request) tags() ([]string, []string) {
	var (
		name  = make([]string, 0, len(req.TagValue))
		value = make([]string, 0, len(req.TagValue))
	)
	if len(req.TagName) == len(req.TagValue) && len(req.TagValue) <= len(req.Dictionary) {
		for _, k := range req.TagName {
			name = append(name, req.Dictionary[int(k)])
		}
		for _, k := range req.TagValue {
			value = append(value, req.Dictionary[int(k)])
		}
	}
	return name, value
}
