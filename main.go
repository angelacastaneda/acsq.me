package main

import (
	"flag"
	"log"
	"net/http"
)

const (
  fullchain = "/etc/letsencrypt/live/angel-castaneda.com/fullchain.pem"
  privkey = "/etc/letsencrypt/live/angel-castaneda.com/privkey.pem"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	flag.Parse() // required before flag is used

  if *addr == ":443" {
    scheme = "https"
  } 

	mux := http.NewServeMux()

  // TODO: Make cooler router
	mux.HandleFunc("/", pageHandler)
	mux.HandleFunc("/posts", pageHandler)
	mux.HandleFunc("/posts/", postHandler)
	mux.HandleFunc("/tags", tagHandler)
	mux.HandleFunc("/tags/", tagHandler)
	mux.HandleFunc("/favicon.ico", faviconHandler)

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

  www := wwwRedirect(mux)

  log.Printf("Starting the server on %s", *addr)

  if *addr == ":443" {
    go http.ListenAndServe(":80", http.HandlerFunc(httpsRedirect))
    err := http.ListenAndServeTLS(*addr, fullchain, privkey, www) 
    log.Fatal(err)
  } else {
    err := http.ListenAndServe(*addr, www)
    log.Fatal(err)
  }
}
