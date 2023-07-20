package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
  htmlDir = filepath.Join(".", "html") // routes to dirs
  staticDir = filepath.Join(".", "static")
  testDomain = "localhost:4000" // todo fix this 
  scheme = "http"
)

func httpsRedirect(w http.ResponseWriter, r *http.Request) {
  target := "https://" + r.Host + r.URL.Path // todo get actual raw path too
  http.Redirect(w, r, target, 302)
}

func wwwRedirect(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if !strings.HasPrefix(r.Host,"www.") && !strings.HasPrefix(r.Host,"en.") && !strings.HasPrefix(r.Host,"es.") && !strings.HasPrefix(r.Host,"de.") {
      http.Redirect(w, r, scheme + "://www." + r.Host + r.RequestURI, 302)
      return
    }

    next.ServeHTTP(w, r)
  })
}

func fancyErrorHandler(httpCode int, w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  w.WriteHeader(httpCode)

  tmpl, err := tmplBinder(
    filepath.Join(htmlDir, "partials", "error_meta" + tmplFileExt),
    filepath.Join(htmlDir, "partials", "error_header" + tmplFileExt),
    filepath.Join(htmlDir, "errors", strconv.Itoa(httpCode) + tmplFileExt),
    filepath.Join(htmlDir, "partials", "footer" + tmplFileExt),
  )
  if err != nil {
    log.Println(err.Error())
    http.Error(w,"Internal Server Error", http.StatusInternalServerError)
    return
  }

  data, err := dataFetcher(r, -404, "")
  if err != nil {
    log.Println(err.Error())
    http.Error(w,"Internal Server Error", http.StatusInternalServerError)
    return
  }

  err = tmpl.ExecuteTemplate(w, "base", data)
  if err != nil {
    log.Println(err.Error())
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }
}

func doesFileExist(filePath string) bool {
  info, err := os.Stat(filePath)
  if err != nil || info.IsDir() {
    return false
  }
  return true
}

func tmplBinder(files ...string) (*template.Template, error) {
  for _, checkFile := range files {
    if !doesFileExist(checkFile) {
      return nil, errors.New("Template file missing " + checkFile)
    }
  }
  
  files = append(files, 
    filepath.Join(htmlDir, "base" + tmplFileExt),
  )

  funcMap := template.FuncMap{
    "translate": translate,
    "lastOne": lastOne,
    "translateTag": translateTag,
    "translateURL": translateURL,
  }

  tmpl, err := template.New("noIdeaWhyThisExists").Funcs(funcMap).ParseFiles(files...)
  if err != nil {
    return nil, err
  }

  return tmpl, nil
}

func dataFetcher(r *http.Request, postQuant int, tagFilter string) (map[string]interface{}, error) {
  data := make(map[string]interface{})
  var err error

  data["Lang"] = langFetcher(r.Host)
  data["Domain"] = testDomain 
  data["Scheme"] = scheme 

  if postQuant != -404 { // todo undo this hack error filter
    data["Path"] = r.URL.Path
  }

  data["Posts"], err = postsSorter(postQuant, tagFilter)
  if err != nil {
    return data, err
  }

  if strings.HasPrefix(r.URL.Path, "/posts/") && len(r.URL.Path) > len("/posts/") {
    data["Post"], err = postFetcher(strings.TrimPrefix(r.URL.Path, "/posts/"))
    if err != nil {
      return data, err
    }
  }

  if r.URL.Path == "/about" || r.URL.Path == "/conoceme" || r.URL.Path == "/uber" {
    data["Song"], data["TrackIndex"] = rockNRoll()
  }

  return data, nil
}

func tmplServer(w http.ResponseWriter, r *http.Request, tmpl *template.Template, postQuant int, tagFilter string) {

  data, err := dataFetcher(r, postQuant, tagFilter)
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(http.StatusInternalServerError, w, r)
    return
  }

  err = tmpl.ExecuteTemplate(w, "base", data)
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(http.StatusInternalServerError, w, r)
    return
  }
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  page := strings.Trim(r.URL.Path,"/")

  if r.URL.Path == "/" {
    page = "index"
  }

  if page != translateURL(langFetcher(r.Host), page) {
  // if langFetcher(r.Host) != "en-US" {
    http.Redirect(w, r, "/" + translateURL(langFetcher(r.Host), page), 302)
    return
    // page = translateURL("en-US", page)
  }

  if page != translateURL("en-US", page) {
    page = translateURL("en-US", page)
  }
    // switch langFetcher(r.Host){
    // case "es-US":
    //   http.Redirect(w, r, "/" + langURLs[page][es], 302)
    // }

  if !doesFileExist(filepath.Join(htmlDir, "pages", page + tmplFileExt)) {
    fancyErrorHandler(http.StatusNotFound, w, r)
    return
  }

  tmpl, err := tmplBinder(
    filepath.Join(htmlDir, "partials", "meta" + tmplFileExt),
    filepath.Join(htmlDir, "partials", "header" + tmplFileExt),
    filepath.Join(htmlDir, "pages", page + tmplFileExt),
    filepath.Join(htmlDir, "partials", "footer" + tmplFileExt),
  )
  if err != nil {
    log.Println(err.Error())
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
  }

  switch r.URL.Path  {
  case "/":
    tmplServer(w, r, tmpl, 3, "articles")
    return
  case "/posts":
    tmplServer(w, r, tmpl, 0, "")
    return
  default:
    tmplServer(w, r, tmpl, -1, "")
  } 
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  if r.URL.Path == "/tags/" || r.URL.Path =="/tags" {
    http.Redirect(w, r, "/posts", 302)
    return
  }

  if !doesFileExist(filepath.Join(htmlDir, r.URL.Path + tmplFileExt)) {
    fancyErrorHandler(http.StatusNotFound, w, r)
    // http.Error(w,"Page Not Found", http.StatusNotFound)
    return
  }

  tmpl, err := tmplBinder(
    filepath.Join(htmlDir, "partials", "meta" + tmplFileExt),
    filepath.Join(htmlDir, "partials", "header" + tmplFileExt),
    filepath.Join(htmlDir, r.URL.Path + tmplFileExt), 
    filepath.Join(htmlDir, "partials", "footer" + tmplFileExt),
  )
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(http.StatusInternalServerError, w, r)
    return
  }

  tagFilter := strings.TrimPrefix(r.URL.Path,"/tags/")
  tmplServer(w, r, tmpl, 0, tagFilter)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  if r.URL.Path == "/posts/" || r.URL.Path == "/posts"{
    http.Redirect(w, r, "/posts", 302)
    return
  }

  if !doesFileExist(filepath.Join(htmlDir, r.URL.Path + tmplFileExt)) {
    fancyErrorHandler(http.StatusNotFound, w, r)
    // http.Error(w,"Page Not Found", http.StatusNotFound)
    return
  }

  tmpl, err := tmplBinder(
    filepath.Join(htmlDir, "partials", "meta" + tmplFileExt),
    filepath.Join(htmlDir, "partials", "post_header" + tmplFileExt),
    filepath.Join(htmlDir, r.URL.Path + tmplFileExt),
    filepath.Join(htmlDir, "partials", "footer" + tmplFileExt),
  )
  if err != nil {
    log.Println(err.Error())
    // fancyErrorHandler(http.StatusInternalServerError, w, r)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }

  tmplServer(w, r, tmpl, -1, "")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(staticDir,"favicon.ico"))
}
