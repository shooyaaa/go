package connector

import (
	"encoding/json"
	"log"
	"net/http"
)

type HttpHandler func(http.ResponseWriter, *http.Request)

type HttpServer struct {
	Root    string
	Addr    string
	Handler map[string]HttpHandler
}

func (hs *HttpServer) Register(key string, handler HttpHandler) {
	hs.Handler[key] = handler
}

func (hs *HttpServer) Info(w http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(map[string]string{"Addr": hs.Addr})
	if err != nil {
		log.Println("Error while marshal json :", err)
	}
	w.Header().Add("Content-Type", "text/json")
	w.Write(data)
}

func (hs *HttpServer) Run() {
	http.HandleFunc("/", hs.serveHome)
	err := http.ListenAndServe(hs.Addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func (hs *HttpServer) serveHome(w http.ResponseWriter, r *http.Request) {
	header := map[string]string{}
	for key, s := range header {
		w.Header().Add(key, s)
	}
	uri := r.URL.Path
	if handler, ok := hs.Handler[uri]; ok {
		handler(w, r)
	} else {
		http.ServeFile(w, r, hs.Root+uri)
	}
}
