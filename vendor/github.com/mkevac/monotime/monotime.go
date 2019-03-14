// Based on code from: https://github.com/aristanetworks/goarista/blob/46272bfb1c042fc2825d312fe33d494e9d13dd6b/atime/nanotime.go

// Copyright (C) 2016  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package monotime provides a fast monotonic clock source.
package monotime

import (
	"time"
	"unsafe"
)

// Make goimports import the unsafe package, which is required to be able
// to use //go:linkname
var _ = unsafe.Sizeof(0)

//go:noescape
//go:linkname nanotime runtime.nanotime
func nanotime() int64

// Nanoseconds returns the current time in nanoseconds from a monotonic clock.
// The time returned is based on some arbitrary platform-specific point in the
// past.  The time returned is guaranteed to increase monotonically at a
// constant rate, unlike time.Now() from the Go standard library, which may
// slow down, speed up, jump forward or backward, due to NTP activity or leap
// seconds.
func Nanoseconds() uint64 {
	return uint64(nanotime())
}

// Now returns current time from a monotonic clock
// i.e. the result of Nanoseconds() as time.Time
func Now() time.Time {
	t := time.Time{}
	return t.Add(time.Duration(Nanoseconds()) * time.Nanosecond)
}
