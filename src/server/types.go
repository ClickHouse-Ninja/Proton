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
