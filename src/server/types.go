package server

import (
	"sort"

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
	sort.Sort(&tagSort{
		name:  name,
		value: value,
	})
	return name, value
}

func (req *request) timerTags() ([][]string, [][]string) {
	var (
		name  = make([][]string, len(req.TimerTagCount))
		value = make([][]string, len(req.TimerTagCount))
	)

	if len(req.TimerTagName) == len(req.TimerTagValue) && len(req.TimerTagValue) <= len(req.Dictionary) {
		var names, values []string
		for _, k := range req.TimerTagName {
			names = append(names, req.Dictionary[int(k)])
		}
		for _, k := range req.TimerTagValue {
			values = append(values, req.Dictionary[int(k)])
		}

		var start uint32
		for idx, ln := range req.TimerTagCount {
			name[idx] = names[start : start+ln]
			value[idx] = values[start : start+ln]
			sort.Sort(&tagSort{
				name:  name[idx],
				value: value[idx],
			})
			start += ln
		}
	}
	return name, value
}

func (req *request) GetTimerValue() (value []float32) {
	if value = req.Request.TimerValue; len(value) != len(req.TimerTagCount) {
		return make([]float32, len(req.TimerTagCount))
	}
	return value
}

func (req *request) GetTimerRuStime() (value []float32) {
	if value = req.Request.TimerRuStime; len(value) != len(req.TimerTagCount) {
		return make([]float32, len(req.TimerTagCount))
	}
	return value
}

func (req *request) GetTimerRuUtime() (value []float32) {
	if value = req.Request.TimerRuUtime; len(value) != len(req.TimerTagCount) {
		return make([]float32, len(req.TimerTagCount))
	}
	return value
}

type tagSort struct {
	name, value []string
}

func (t *tagSort) Len() int { return len(t.name) }

func (t *tagSort) Less(i, j int) bool { return t.name[i] < t.name[j] }

func (t *tagSort) Swap(i, j int) {
	t.name[i], t.name[j] = t.name[j], t.name[i]
	t.value[i], t.value[j] = t.value[j], t.value[i]
}

var _ sort.Interface = (*tagSort)(nil)
