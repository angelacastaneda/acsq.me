package main

import (
	"log"
	"net/http"
	"html/template"
  "strings"
)

func internalServerError(w http.ResponseWriter, r *http.Request) {
  // fancy 500 page
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  w.WriteHeader(http.StatusInternalServerError)
  files := []string{
    "./html/base.tmpl.html",
    "./html/errors/500.tmpl.html",
    "./html/partials/error_meta.tmpl.html",
    "./html/partials/error_header.tmpl.html",
  }

  // parses templates in ui/html
  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    http.Error(w,"Internal Server Error",http.StatusInternalServerError)
    return
  }

  // tries to serve template 
  err = ts.ExecuteTemplate(w,"base",nil)
  if err != nil {
    log.Println(err.Error())
    http.Error(w,"Internal Server Error",http.StatusInternalServerError)
  }
}

func pageNotFound(w http.ResponseWriter, r *http.Request) {
  // fancy 404 page
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  w.WriteHeader(http.StatusNotFound)
  files := []string{
    "./html/base.tmpl.html",
    "./html/errors/404.tmpl.html",
    "./html/partials/error_meta.tmpl.html",
    "./html/partials/error_header.tmpl.html",
  }

  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    return
  }

  err = ts.ExecuteTemplate(w,"base",nil)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
  }
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  files := []string{
    "./html/base.tmpl.html",
    "./html/pages/index.tmpl.html",
  }

  if r.URL.Path != "/" {
    pageNotFound(w, r)
    return
  }

  if !strings.HasPrefix(r.Host,"www.") && !strings.HasPrefix(r.Host,"en.") && !strings.HasPrefix(r.Host,"es.") && !strings.HasPrefix(r.Host,"de.") {
    http.Redirect(w, r, "http://www."+r.Host+r.RequestURI, 302)
  }

  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    return
  }

  fileNames := []string{"cool","epic"}

  data := struct {
    FileNames []string
  }{
    FileNames: fileNames,
  }


  err = ts.ExecuteTemplate(w, "base", data)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
  }
}

func pageHandler(page string) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    if !strings.HasPrefix(r.Host,"www.") && !strings.HasPrefix(r.Host,"en.") && !strings.HasPrefix(r.Host,"es.") && !strings.HasPrefix(r.Host,"de.") {
      http.Redirect(w, r, "http://www."+r.Host+r.RequestURI, 302)
    }

    w.Header().Set("Content-Type","text/html; charset=utf-8")
    files := []string{
      "./html/base.tmpl.html",
      "./html/pages/"+page+".tmpl.html",
    }

    ts, err := template.ParseFiles(files...)
    if err != nil {
      log.Println(err.Error())
      internalServerError(w, r)
      return
    }

    err = ts.ExecuteTemplate(w,"base",nil)
    if err != nil {
      internalServerError(w, r)
      http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
  }
}

func postsPageHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  files := []string{
    "./html/base.tmpl.html",
    "./html/pages/posts.tmpl.html",
  }

  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    return
  }

  fileNames := []string{"cool","epic"}

  data := struct {
    FileNames []string
  }{
    FileNames: fileNames,
  }


  err = ts.ExecuteTemplate(w, "base", data)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
  }
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}

func postHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  if r.URL.Path == "/posts/" {
    http.Redirect(w, r, "/posts", http.StatusMovedPermanently)
    return
  }

  if r.URL.Path == "/posts" {
    postsPageHandler(w, r)
    return
  }

  url := strings.TrimPrefix(r.URL.Path,"/posts/")

  // TODO need to check if the file exists in the first place

  files := []string{
    "./html/base.tmpl.html",
    "./html/partials/post_header.tmpl.html",
    "./html/posts/"+url+".tmpl.html",
  }

  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    pageNotFound(w, r) // this is a very scuffed method
    return
  }

  err = ts.ExecuteTemplate(w,"base",nil)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
  }
}
