package main

import (
	"log"
	"net/http"
	"time"

	"fmt"

	"strconv"

	"github.com/kavehmz/counter"
)

var c *counter.Counter

func handler(w http.ResponseWriter, r *http.Request) {
	ch := make(chan int)
	c.Inc(ch)
	fmt.Fprintf(w, strconv.Itoa(<-ch))
}

func main() {
	var e error
	c, e = counter.Init("/tmp/counter", time.Second, 60)
	if e != nil {
		log.Fatal(e)
	}
	http.HandleFunc("/count", handler)
	http.ListenAndServe(":8080", nil)
}
