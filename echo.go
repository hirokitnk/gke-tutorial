package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

type Action struct {
	Act      bool
	Payload  []string
	Response chan<- []string
}

func main() {
	log.Println("application server started")
	actCh := make(chan Action)
	go act(actCh)
	http.HandleFunc("/", makeHTTPHandler(actCh))
	e := http.ListenAndServe(":8080", nil)
	if e != nil {
		log.Println(e)
		panic(e)
	}
}

func makeHTTPHandler(actCh chan<- Action) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var resp = make(chan []string)
		hostname, _ := os.Hostname()
		req := []string{
			r.Method,
			hostname,
			"version :" + os.Getenv("VERSION"),
			r.Header.Get("User-Agent"),
		}
		act := Action{
			Act:      true,
			Payload:  req,
			Response: resp,
		}
		actCh <- act
		re := <-resp
		w.Write([]byte(strings.Join(re, ",\n")))
	}
}

func act(actCh <-chan Action) {
	for {
		act := <-actCh
		ret := act.Payload
		act.Response <- ret
	}
}
