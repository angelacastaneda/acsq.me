package main

import (
	"flag"
	"log"
	"net/http"

	"acsq.me/dblog"
)

var (
	langs  = []string{"en-US", "es-US", "de-DE"}
	scheme string
)

func main() {
	if err := dblog.MakeDB(); err != nil {
		log.Fatal(err)
	}

	addr := flag.String("addr", ":4000", "HTTP Network Address")
	https := flag.Bool("https", false, "TLS Encryption")
	flag.Parse() // required before flag is used

	if *https {
		scheme = "https"
	} else {
		scheme = "http"
	}

	mux := http.NewServeMux()

	// TODO: Make cooler router
	mux.HandleFunc("/", pageHandler)
	for _, lang := range langs {
		mux.HandleFunc(translatePath(lang, "/posts"), pageHandler)
		mux.HandleFunc(translatePath(lang, "/posts.html"), pageHandler)
		mux.HandleFunc(translatePath(lang, "/posts/{year}/{month}/{day}/"), postHandler)
		mux.HandleFunc(translatePath(lang, "/posts/"), postDateRedirect)
		mux.HandleFunc(translatePath(lang, "/tags/"), tagHandler)
		mux.HandleFunc(translatePath(lang, "/recommend"), recommendHandler)
		mux.HandleFunc(translatePath(lang, "/recommend.html"), recommendHandler)
	}
	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/cv.pdf", cvHandler)
	mux.HandleFunc("/angelcastaneda.asc", pgpHandler)
	mux.HandleFunc("/atom.xml", feedHandler)
	mux.HandleFunc("POST /submit", apiHandler)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting the server on %s", *addr)

	err := http.ListenAndServe(*addr, gzipHandler(redirectWWW(mux)))
	log.Fatal(err)
}
