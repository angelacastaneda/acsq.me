package main

import (
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
      http.Redirect(w, r, "http://www." + r.Host + r.RequestURI, 302)
      return
    }

    next.ServeHTTP(w, r)
  })
}

func fancyErrorHandler(httpCode int, w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  w.WriteHeader(httpCode)

  tmpl, err := tmplBinder(
    filepath.Join(htmlDir, "errors", strconv.Itoa(httpCode) + tmplFileExt),
    filepath.Join(htmlDir, "partials", "error_meta" + tmplFileExt),
    filepath.Join(htmlDir, "partials", "error_header" + tmplFileExt),
  )
  if err != nil {
    log.Println(err.Error())
    http.Error(w,"Internal Server Error", http.StatusInternalServerError)
    return
  }

  data, err := dataFetcher(r, -1, "")
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

func tmplBinder(files ...string) (*template.Template, error) {
  files = append(files, 
    filepath.Join(htmlDir, "base" + tmplFileExt),
  )

  funcMap := template.FuncMap{
    "translate": translate,
    "lastOne": lastOne,
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
  // required in all
  data["Lang"] = langFetcher(r.Host)
  data["Domain"] = testDomain 
  data["Scheme"] = scheme 
  // data["Path"] = r.URL.Path

  // posts doesn't do anything crazy with a negative number
  data["Posts"], err = postsSorter(postQuant, tagFilter)
  if err != nil {
    return data, err
  }

  // for individual posts
  if strings.HasPrefix(r.URL.Path, "/posts/") && len(r.URL.Path) > len("/posts/") {
    data["Post"], err = postFetcher(strings.TrimPrefix(r.URL.Path, "/posts/"))
    if err != nil {
      return data, err
    }
  }

  // for cool jukebox
  if r.URL.Path == "/about" {
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  if r.URL.Path != "/" {
    fancyErrorHandler(http.StatusNotFound, w, r)
    return
  }

  tmpl, err := tmplBinder(
    filepath.Join(htmlDir, "pages", "index" + tmplFileExt),
  )
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(http.StatusInternalServerError, w, r)
    return
  }

  tmplServer(w, r, tmpl, 3, "articles")
}

func pageHandler(page string) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type","text/html; charset=utf-8")

    tmpl, err := tmplBinder(
      filepath.Join(htmlDir, "pages", page + tmplFileExt),
    )
    if err != nil {
      log.Println(err.Error())
      fancyErrorHandler(http.StatusInternalServerError, w, r)
      return
    }

    tmplServer(w, r, tmpl, -1, "")
  }
}

func postsPageHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  tmpl, err := tmplBinder(
    filepath.Join(htmlDir, "pages", "posts" + tmplFileExt),
  )
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(http.StatusInternalServerError, w, r)
    return
  }

  tmplServer(w, r, tmpl, 0, "")
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  if r.URL.Path == "/tags/" || r.URL.Path =="/tags" {
    http.Redirect(w, r, "/posts", 302)
    return
  }

  tagURL := strings.TrimPrefix(r.URL.Path,"/tags/")

  // TODO need to check if the file exists in the first place

  tmpl, err := tmplBinder(
    filepath.Join(htmlDir, "tags", tagURL + tmplFileExt), 
  )
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(http.StatusNotFound, w, r)// this is a very scuffed method
    return
  }

  tmplServer(w, r, tmpl, 0, tagURL)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  if r.URL.Path == "/posts/" || r.URL.Path == "/posts"{
    http.Redirect(w, r, "/posts", 302)
    return
  }

  // check if file exists
  fileName := strings.TrimPrefix(filepath.Clean(r.URL.Path), "/posts/")
  postPath := filepath.Join(htmlDir, "posts", fileName + tmplFileExt)

  info, err := os.Stat(postPath)
  if err != nil || info.IsDir() {
    // fancyErrorHandler(http.StatusNotFound, w, r)
    http.Error(w,"Page Not Found", http.StatusNotFound)
    return
  }

  tmpl, err := tmplBinder(
    filepath.Join(htmlDir, "partials", "post_header" + tmplFileExt),
    postPath,
  )
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(http.StatusInternalServerError, w, r)
    return
  }

  tmplServer(w, r, tmpl, -1, "")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(staticDir,"favicon.ico"))
}
