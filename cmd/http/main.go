package main

import (
	"flag"
	"log"
	"net/http"

	"angel-castaneda.com/dblog"
)

var (
	fullchain = "/etc/letsencrypt/live/angel-castaneda.com/fullchain.pem"
	privkey   = "/etc/letsencrypt/live/angel-castaneda.com/privkey.pem"
	langs     = []string{"en-US", "es-US", "de-DE"}
	scheme    string
)

func main() {
	if err := dblog.MakeDB(); err != nil {
		log.Fatal(err)
	}

	addr := flag.String("addr", ":4000", "HTTP Network Address")
	flag.Parse() // required before flag is used

	if *addr == ":443" {
		scheme = "https"
	} else {
		scheme = "http"
	}

	mux := http.NewServeMux()

	// TODO: Make cooler router
	mux.HandleFunc("/", pageHandler)
	for _, lang := range langs {
		mux.HandleFunc(translatePath(lang, "/posts"), pageHandler)
		mux.HandleFunc(translatePath(lang, "/posts/"), postHandler)
		mux.HandleFunc(translatePath(lang, "/tags/"), tagHandler)
	}
	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/cv.pdf", cvHandler)
	mux.HandleFunc("/angelcastaneda.asc", pgpHandler)
	mux.HandleFunc("/pgp", pgpHandler)
	mux.HandleFunc("/atom.xml", feedHandler)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting the server on %s", *addr)

	if *addr == ":443" {
		go http.ListenAndServe(":80", http.HandlerFunc(redirectHTTPS))
		err := http.ListenAndServeTLS(*addr, fullchain, privkey, gzipHandler(redirectWWW(mux)))
		log.Fatal(err)
	} else {
		err := http.ListenAndServe(*addr, gzipHandler(redirectWWW(mux)))
		log.Fatal(err)
	}
}
