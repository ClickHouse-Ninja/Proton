package server

import (
	"time"

	"github.com/ClickHouse-Ninja/Proton/proto/pinba"
)

type Options struct {
	DSN         string
	Address     string
	BacklogSize int
	Concurrency int
}

type requestContainer struct {
	pinba.Request
	timestamp time.Time
}

func (req *requestContainer) tags() ([]string, []string) {
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
