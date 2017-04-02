package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/kavehmz/counter"
)

var c *counter.Counter

func handler(w http.ResponseWriter, r *http.Request) {
	ch := make(chan int)
	c.Inc(ch)
	fmt.Fprintf(w, strconv.Itoa(<-ch))
}

func main() {
	if os.Getenv("PORT") == "" {
		log.Fatal("PORT must be set")
	}
	if os.Getenv("HOST") == "" {
		log.Fatal("HOST must be set")
	}
	if os.Getenv("REDISURL") == "" {
		log.Fatal("REDISURL must be set")
	}

	var e error
	c, e = counter.Init("/tmp/counter", time.Second, 60)
	if e != nil {
		log.Fatal(e)
	}

	http.HandleFunc("/", handler)

	maxServingClients := 2
	maxClientsPool := make(chan bool, maxServingClients)

	server := &http.Server{
		Addr:    ":" + os.Getenv("PORT"),
		Handler: nil,
		ConnState: func(conn net.Conn, state http.ConnState) {
			switch state {
			case http.StateNew:
				maxClientsPool <- true
			case http.StateClosed, http.StateHijacked:
				<-maxClientsPool

			}
		},
	}
	log.Fatal(server.ListenAndServe())
}
