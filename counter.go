/*
Package counter implements a counter which will persist the result in a local file.
*/
package counter

import (
	"os"
	"time"
)

type stat struct {
	epoch int64
	count int
	next  *stat
}

// Counter defines parameter for the counter.
type Counter struct {
	file          *os.File
	bufferTimeout time.Duration
	count         int
	incChannel    chan chan int
	statBegin     *stat
	statEnd       *stat
	history       int
}

// Init will setup a counter and loads the initial value if the file exists
// It accepts a filename and buffersize
func Init(fn string, btime time.Duration, h int) (*Counter, error) {
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	c := &Counter{}
	c.incChannel = make(chan chan int)
	c.bufferTimeout = btime
	c.file = f

	c.statBegin = &stat{}
	b := c.statBegin
	for i := 1; i < h; i++ {
		b.next = &stat{}
		b = b.next
	}
	b.next = c.statBegin
	c.statEnd = c.statBegin
	c.history = h

	go c.loop()
	return c, nil
}

// Inc will increment the counter
func (c *Counter) Inc(ret chan int) {
	c.incChannel <- ret
}

func (c *Counter) loop() {
	for {
		select {
		case ret := <-c.incChannel:
			ret <- c.inc()
		case <-time.After(c.bufferTimeout):
			c.file.Sync()
		}
	}
}

func (c *Counter) inc() int {
	now := time.Now().Unix()

	for c.statBegin.epoch <= now-int64(c.history) {
		c.count -= c.statBegin.count
		c.statBegin.count = 0
		c.statBegin.epoch = 0
		if c.statBegin == c.statEnd {
			break
		}
		c.statBegin = c.statBegin.next
	}

	if c.statEnd.epoch != now && c.statEnd.epoch != 0 {
		c.statEnd = c.statEnd.next
	}
	c.statEnd.epoch = now
	c.statEnd.count++
	c.count++
	return c.count
}

func (c *Counter) read() {

}

func (c *Counter) write() {

}
