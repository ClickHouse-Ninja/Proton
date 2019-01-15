package gopinba

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/mkevac/gopinba/Pinba"
)

type Client struct {
	initialized bool
	address     string
	conn        net.Conn
}

type Timer struct {
	Tags map[string]string

	// private stuff for simpler api
	stopped  bool
	started  time.Time
	duration time.Duration
}

type Request struct {
	Hostname     string
	ServerName   string
	ScriptName   string
	Schema       string
	RequestCount uint32
	RequestTime  time.Duration
	DocumentSize uint32
	MemoryPeak   uint32
	Utime        float32
	Stime        float32
	timers       []Timer
	Status       uint32
	Tags         map[string]string
	lk           sync.Mutex
}

func NewClient(address string) (*Client, error) {
	pc := &Client{address: address}
	conn, err := net.Dial("udp", address)
	if err != nil {
		return nil, err
	}

	pc.conn = conn
	pc.initialized = true

	return pc, nil
}

func iN(haystack []string, needle string) (int, bool) {
	for i, s := range haystack {
		if s == needle {
			return i, true
		}
	}
	return -1, false
}

func mergeTags(req *Pinba.Request, tags map[string]string) {
	for k, v := range tags {
		{
			pos, exists := iN(req.Dictionary, k)
			if !exists {
				req.Dictionary = append(req.Dictionary, k)
				pos = len(req.Dictionary) - 1
			}

			req.TagName = append(req.TagName, uint32(pos))
		}

		{
			pos, exists := iN(req.Dictionary, v)
			if !exists {
				req.Dictionary = append(req.Dictionary, v)
				pos = len(req.Dictionary) - 1
			}

			req.TagValue = append(req.TagValue, uint32(pos))
		}
	}
}

func mergeTimerTags(req *Pinba.Request, tags map[string]string) {
	req.TimerTagCount = append(req.TimerTagCount, uint32(len(tags)))

	for k, v := range tags {
		{
			pos, exists := iN(req.Dictionary, k)
			if !exists {
				req.Dictionary = append(req.Dictionary, k)
				pos = len(req.Dictionary) - 1
			}

			req.TimerTagName = append(req.TimerTagName, uint32(pos))
		}

		{
			pos, exists := iN(req.Dictionary, v)
			if !exists {
				req.Dictionary = append(req.Dictionary, v)
				pos = len(req.Dictionary) - 1
			}

			req.TimerTagValue = append(req.TimerTagValue, uint32(pos))
		}
	}
}

func preallocateArrays(req *Pinba.Request, timers []Timer) {

	// calculate (max) final lengths for all arrays
	nTimers := 0
	nTags := 0
	for _, timer := range timers {
		nTimers++
		nTags += len(timer.Tags)
	}

	// construct arrays capable of holding all possible values to reduce allocations
	req.TimerHitCount = make([]uint32, 0, nTimers) // number of hits for each timer
	req.TimerValue = make([]float32, 0, nTimers)   // timer value for each timer
	req.Dictionary = make([]string, 0, nTags)      // all strings used in timer tag names/values
	req.TimerTagCount = make([]uint32, 0, nTimers) // number of tags for each timer
	req.TimerTagName = make([]uint32, 0, nTags)    // flat array of all tag names (as offsets into dictionary) laid out sequentially for all timers
	req.TimerTagValue = make([]uint32, 0, nTags)   // flat array of all tag values (as offsets into dictionary) laid out sequentially for all timers
}

func (pc *Client) SendRequest(request *Request) error {

	if !pc.initialized {
		return fmt.Errorf("Client not initialized")
	}

	pbreq := Pinba.Request{
		Hostname:     request.Hostname,
		ServerName:   request.ServerName,
		ScriptName:   request.ScriptName,
		RequestCount: request.RequestCount,
		RequestTime:  float32(request.RequestTime.Seconds()),
		DocumentSize: request.DocumentSize,
		MemoryPeak:   request.MemoryPeak,
		RuUtime:      request.Utime,
		RuStime:      request.Stime,
		Status:       request.Status,
		Schema:       request.Schema,
	}

	preallocateArrays(&pbreq, request.timers)

	pbreq.TagValue = make([]uint32, 0, len(request.Tags))
	pbreq.TagName = make([]uint32, 0, len(request.Tags))

	mergeTags(&pbreq, request.Tags)

	for _, timer := range request.timers {
		pbreq.TimerHitCount = append(pbreq.TimerHitCount, 1)
		pbreq.TimerValue = append(pbreq.TimerValue, float32(timer.duration.Seconds()))
		mergeTimerTags(&pbreq, timer.Tags)
	}

	buf := make([]byte, pbreq.Size())

	n, err := pbreq.MarshalTo(buf)
	if err != nil {
		return err
	}

	_, err = pc.conn.Write(buf[:n])
	if err != nil {
		return err
	}

	return nil
}

func (req *Request) AddTimer(timer *Timer) {
	req.lk.Lock()
	defer req.lk.Unlock()

	req.timers = append(req.timers, *timer)
}

// this is exactly the same as AddTimer
//  exists only to have api naming similar to pinba php extension
func (req *Request) TimerAdd(timer *Timer) {
	timer.Stop()
	req.AddTimer(timer)
}

func TimerStart(tags map[string]string) *Timer {
	return &Timer{
		duration: 0,
		Tags:     tags,
		stopped:  false,
		started:  now(),
	}
}

func NewTimer(tags map[string]string, duration time.Duration) *Timer {
	return &Timer{
		duration: duration,
		Tags:     tags,
		stopped:  true,
		started:  now().Add(-duration),
	}
}

func (t *Timer) Stop() {
	if !t.stopped {
		t.stopped = true
		t.duration = now().Sub(t.started)
	}
}

func (t *Timer) GetDuration() time.Duration {
	return t.duration
}
