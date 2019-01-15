// +build !go1.9

package gopinba

import (
	"github.com/mkevac/monotime"
	"time"
)

func now() time.Time {
	return monotime.Now()
}
