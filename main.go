package main

import (
	"flag"
	"log"
	"net/http"
)

//func homeHandler(w http.ResponseWriter, r *http.Request) {
//  w.Header().Set("Content-Type", "text/html; charset=utf-8")
//  fmt.Fprint(w, "<h1>Welcome to my epic site B-)</h1>")
//}
//
//func contactHandler(w http.ResponseWriter, r *http.Request) {
//  w.Header().Set("Content-Type", "text/html; charset=utf-8")
//  fmt.Fprint(w, "<h1>Contact Page</h1> <p>Email me at <a href=\"mailto:angel@carbajal-castaneda.com\">angel@carbajal-castaneda.com</a> to get in touch</p>")
//}

//func pathHandler(w http.ResponseWriter, r *http.Request) {
//  switch r.URL.Path {
//  case "/":
//    homeHandler(w,r)
//  case "/contact":
//    contactHandler(w,r)
//  default:
//    // TODO: handle page not found error
//    http.Error(w, "<h2>404 page not found</h2>", http.StatusNotFound)
//  }
//}

//type Router struct {}
//
//func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//  switch r.URL.Path {
//  case "/":
//    home(w,r)
//  case "/about":
//    about(w,r)
//  case "/library":
//    library(w,r)
//  default:
//    // TODO: handle page not found error
//    http.Error(w, "<h2>404 page not found</h2>", http.StatusNotFound)
//  }
//}

func main() {
//  var router Router
  addr := flag.String("addr", ":4000", "HTTP network address")
  flag.Parse()

  mux := http.NewServeMux()

  fileServer := http.FileServer(http.Dir("./static"))
  mux.Handle("/static/",http.StripPrefix("/static",fileServer))

  mux.HandleFunc("/", home)
  mux.HandleFunc("/about", about)
  mux.HandleFunc("/library", library)
  mux.HandleFunc("/posts", posts)
  mux.HandleFunc("/todo", todo)
  mux.HandleFunc("/friends", friends)
  mux.HandleFunc("/favicon.ico", favicon)


  log.Printf("Starting the server on %s", *addr)
  err := http.ListenAndServe(*addr, mux)
  log.Fatal(err)
}
