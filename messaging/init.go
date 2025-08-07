package messaging

import "time"

var startedAt time.Time

func init() {
	startedAt = time.Now()
}
