package server

import (
	"sync/atomic"
	"time"

	"github.com/kshvakov/clickhouse/lib/cityhash102"
)

var timestamp = uint32(time.Now().Unix())

func init() {
	go func() {
		for tick := time.Tick(time.Second); ; {
			select {
			case t := <-tick:
				atomic.StoreUint32(&timestamp, uint32(t.Unix()))
			}
		}
	}()
}

func now() uint32 {
	return atomic.LoadUint32(&timestamp)
}

func cityHash64(str string) uint64 {
	cityhash := cityhash102.New64()
	cityhash.Write([]byte(str))
	return cityhash.Sum64()
}
