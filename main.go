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

  // TODO: Make cooler router
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/about", pageHandler("about"))
	mux.HandleFunc("/library", pageHandler("library"))
	mux.HandleFunc("/todo", pageHandler("todo"))
	mux.HandleFunc("/friends", pageHandler("friends"))
	mux.HandleFunc("/posts", postsPageHandler)
	mux.HandleFunc("/posts/", postHandler)
	mux.HandleFunc("/tags", tagHandler)
	mux.HandleFunc("/tags/", tagHandler)
	mux.HandleFunc("/favicon.ico", faviconHandler)



	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting the server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
