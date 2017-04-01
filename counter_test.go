package counter

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestCounter_Init(t *testing.T) {
	tmpfile, _ := ioutil.TempFile("", "test")
	defer os.Remove(tmpfile.Name())

	_, e := Init("", 0, 10)
	if e == nil {
		t.Error("Init did not fail with wrong filename")
	}

	_, e = Init(tmpfile.Name(), 0, 10)
	if e != nil {
		t.Error("Init failed with correct params")
	}

}

func TestCounter_Inc(t *testing.T) {
	tmpfile, _ := ioutil.TempFile("", "test")
	defer os.Remove(tmpfile.Name())

	c, _ := Init(tmpfile.Name(), 0, 2)
	ch := make(chan int)
	c.Inc(ch)
	if n := <-ch; n != 1 {
		t.Error("Incorrect increment", n)
	}

	c.Inc(ch)
	if n := <-ch; n != 2 {
		t.Error("Incorrect increment", n)
	}
	time.Sleep(time.Second * 3)
	c.Inc(ch)
	if n := <-ch; n != 1 {
		t.Error("Incorrect increment", n)
	}
}
