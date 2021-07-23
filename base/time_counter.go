package base

import "time"

type TimeCounter struct {
	begin	time.Time
}

func (c *TimeCounter) Reset() {
	c.begin = time.Now()
}

func (c *TimeCounter) Count() int64 {
	return time.Now().Sub(c.begin).Microseconds()
}
