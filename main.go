package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
  // flag address for function
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/", home)
	mux.HandleFunc("/about", about)
	mux.HandleFunc("/library", library)
	mux.HandleFunc("/posts", posts)
	mux.HandleFunc("/posts/", post)
	mux.HandleFunc("/todo", todo)
	mux.HandleFunc("/friends", friends)
	mux.HandleFunc("/favicon.ico", favicon)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting the server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
