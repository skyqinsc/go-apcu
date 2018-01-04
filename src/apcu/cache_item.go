package apcu

import (
	"time"
)

type Item struct {
	Object     interface{}
	Expiration int64
}

func (item *Item) Expired() bool {
	if 0 == item.Expiration {
		return false
	}

	return time.Now().UnixNano() >= item.Expiration
}
