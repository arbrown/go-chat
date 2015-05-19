// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"text/template"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("home.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
}

func main() {
	flag.Parse()
	go h.run()

	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	if port == "" {
		port = "2015"
	}
	httpMux, wsMux := http.NewServeMux(), http.NewServeMux()
	httpMux.HandleFunc("/", serveHome)
	wsMux.HandleFunc("/ws", serveWs)

	// open shift requires web sockets to be on this port
	wsPort := "8000"
	bind := fmt.Sprintf("%s:%s", host, port)
	wsBind := fmt.Sprintf("%s:%s", host, wsPort)

	go func() {
		fmt.Printf("Listening on %s\n", wsBind)
		http.ListenAndServe(wsBind, wsMux)
	}()

	err := http.ListenAndServe(bind, httpMux)
	if err != nil {
		panic("ListenAndServe:" + err.Error())
	}
}
