/*
Package counter implements a counter which will persist the result in a local file.
*/
package counter

import (
	"encoding/gob"
	"log"
	"os"
	"time"
)

type item struct {
	Epoch int64
	Count int
}

type stat struct {
	item
	next *stat
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
// It accepts a filename and buffersize and length of history to keep
func Init(fn string, btime time.Duration, h int) (*Counter, error) {
	f, err := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	c := &Counter{}
	c.incChannel = make(chan chan int)
	c.bufferTimeout = btime
	c.file = f
	c.history = h

	c.load()
	if err != nil {
		return nil, err
	}

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
			c.write()
		}
	}
}

func (c *Counter) inc() int {
	now := time.Now().Unix()

	for c.statBegin.Epoch <= now-int64(c.history) {
		c.count -= c.statBegin.Count
		c.statBegin.Count = 0
		c.statBegin.Epoch = 0
		if c.statBegin == c.statEnd {
			break
		}
		c.statBegin = c.statBegin.next
	}

	if c.statEnd.Epoch != now && c.statEnd.Epoch != 0 {
		c.statEnd = c.statEnd.next
	}
	c.statEnd.Epoch = now
	c.statEnd.Count++
	c.count++
	return c.count
}

func (c *Counter) load() error {
	var t []item

	_, err := c.file.Seek(0, 0)
	if err == nil {
		decoder := gob.NewDecoder(c.file)
		decoder.Decode(&t)
	}

	c.statBegin = &stat{}
	c.set(c.statBegin, t, 0)
	b := c.statBegin
	for i := 1; i < c.history; i++ {
		b.next = &stat{}
		b = b.next
		c.set(b, t, i)
	}
	b.next = c.statBegin
	if c.statEnd == nil {
		c.statEnd = c.statBegin
	}

	return err
}

func (c *Counter) set(s *stat, t []item, j int) {
	if len(t) > j {
		s.Epoch = t[j].Epoch
		s.Count = t[j].Count
		c.count += t[j].Count
		c.statEnd = s
	}
}

func (c *Counter) write() error {
	var t []item
	b := c.statBegin
	for i := 0; i < c.history && b.Epoch > 0; i++ {
		t = append(t, item{b.Epoch, b.Count})
		b = b.next
	}

	e := c.file.Truncate(0)
	chk(e)
	_, e = c.file.Seek(0, 0)
	chk(e)
	encoder := gob.NewEncoder(c.file)
	e = c.file.Sync()
	chk(e)
	encoder.Encode(t)
	return nil
}

func chk(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
