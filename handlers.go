package main

import (
	"log"
	"net/http"
	"text/template"
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

  // parses templates in ui/html
  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
    return
  }

  // tries to serve template 
  err = ts.ExecuteTemplate(w,"base",nil)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
  }
}

func home(w http.ResponseWriter, r *http.Request) {

  // template files
  files := []string{
    "./html/base.tmpl.html",
    "./html/pages/home.tmpl.html",
  }

  // TODO: fancy 404 handling
  if r.URL.Path != "/" {
    //http.NotFound(w, r)
    pageNotFound(w, r)
    return
  }

  // parses templates in ui/html
  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
    return
  }

  // tries to serve template 
  err = ts.ExecuteTemplate(w,"base",nil)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
  }
}

func about(w http.ResponseWriter, r *http.Request) {

  // template files
  files := []string{
    "./html/base.tmpl.html",
    "./html/pages/about.tmpl.html",
  }

  // parses templates in ui/html
  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
    return
  }

  // tries to serve template 
  err = ts.ExecuteTemplate(w,"base",nil)
  if err != nil {
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  }
}

func library(w http.ResponseWriter, r *http.Request) {

  // template files
  files := []string{
    "./html/base.tmpl.html",
    "./html/pages/library.tmpl.html",
  }

  // parses templates in ui/html
  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
    return
  }

  // tries to serve template 
  err = ts.ExecuteTemplate(w,"base",nil)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
  }
}

func posts(w http.ResponseWriter, r *http.Request) {

  // template files
  files := []string{
    "./html/base.tmpl.html",
    "./html/pages/posts.tmpl.html",
  }

  // parses templates in ui/html
  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
    return
  }

  // tries to serve template 
  err = ts.ExecuteTemplate(w,"base",nil)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
  }
}

func todo(w http.ResponseWriter, r *http.Request) {

  // template files
  files := []string{
    "./html/base.tmpl.html",
    "./html/pages/todo.tmpl.html",
  }

  // parses templates in ui/html
  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
    return
  }

  // tries to serve template 
  err = ts.ExecuteTemplate(w,"base",nil)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
  }
}

func friends(w http.ResponseWriter, r *http.Request) {

  // template files
  files := []string{
    "./html/base.tmpl.html",
    "./html/pages/friends.tmpl.html",
  }

  // parses templates in ui/html
  ts, err := template.ParseFiles(files...)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
    return
  }

  // tries to serve template 
  err = ts.ExecuteTemplate(w,"base",nil)
  if err != nil {
    log.Println(err.Error())
    internalServerError(w, r)
    //http.Error(w,"Internal Server Error",http.StatusInternalServerError)
  }
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/favicon.ico")
}
