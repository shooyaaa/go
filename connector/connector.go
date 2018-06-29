package connector

import (
	"log"
	"net/http"
	"strconv"
	"time"
)

type HttpServer struct {
	Root string
	Addr string
}

func (hs *HttpServer) Run() {
	http.HandleFunc("/", hs.serveHome)
	err := http.ListenAndServe(hs.Addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (hs *HttpServer) serveHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("a request come")
	header := map[string]string{
		"Content-Type": "text/html",
		"Date":         "2018-20-12",
	}
	for key, s := range header {
		w.Header().Add(key, s)
	}
	w.WriteHeader(http.StatusOK)
	for i := 0; i < 100; i++ {
		time.Sleep(1000 * time.Millisecond)
		response := make([]byte, 10)
		num := []byte(strconv.Itoa(i))
		log.Printf("num %v", num)
		response = append(response, num...)
		response = append(response, []byte("<br/>")...)
		w.Write(response)
		w.(http.Flusher).Flush()
	}
}
