package main

import (
	"bytes"
	"compress/gzip"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
  htmlDir = filepath.Join(".", "html") // routes to dirs
  staticDir = filepath.Join(".", "static")
)

type gzipResponseWriter struct {
  io.Writer
  http.ResponseWriter
}

func (grw gzipResponseWriter) Write(data []byte) (int, error) {
  return grw.Writer.Write(data)
}

func gzipHandler(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
      next.ServeHTTP(w, r)
      return
    }

    w.Header().Set("Content-Encoding", "gzip")
    gzipWriter := gzip.NewWriter(w)
    defer gzipWriter.Close()
    gzippedResponseWriter := gzipResponseWriter{Writer: gzipWriter, ResponseWriter: w}
    next.ServeHTTP(gzippedResponseWriter, r)
  })
}

func redirectHTTPS(w http.ResponseWriter, r *http.Request) {
  if r.TLS != nil {
    http.Error(w, "HTTPS already working", http.StatusBadRequest)
  }
  target := "https://" + r.Host + r.RequestURI
  http.Redirect(w, r, target, http.StatusMovedPermanently)
}

func redirectWWW(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    if !strings.HasPrefix(r.Host,"www.") && !strings.HasPrefix(r.Host,"en.") && !strings.HasPrefix(r.Host,"es.") && !strings.HasPrefix(r.Host,"de.") {
      http.Redirect(w, r, scheme + "://www." + r.Host + r.RequestURI, 302)
      return
    }

    next.ServeHTTP(w, r)
  })
}

func fancyErrorHandler(w http.ResponseWriter, r *http.Request, httpCode int) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  w.WriteHeader(httpCode)

  tmpl, err := bindTMPL(
    filepath.Join(htmlDir, "base" + tmplFileExt),
    filepath.Join(htmlDir, "partials", "error_meta" + tmplFileExt),
    filepath.Join(htmlDir, "partials", "error_header" + tmplFileExt),
    filepath.Join(htmlDir, "errors", strconv.Itoa(httpCode) + tmplFileExt),
  )
  if err != nil {
    log.Println(err.Error())
    http.Error(w,"Internal Server Error", http.StatusInternalServerError)
    return
  }

  data, err := fetchData(r.Host, "/", -1, "")
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

func doesFileExist(pathToFile string) bool {
  info, err := os.Stat(filepath.Clean(pathToFile))
  if err != nil || info.IsDir() {
    return false
  }
  return true
}

func bindTMPL(files ...string) (*template.Template, error) {
  for _, checkFile := range files {
    if !doesFileExist(checkFile) {
      return nil, errors.New("Template file missing " + checkFile)
    }
  }

  funcMap := template.FuncMap{
    "translate": translate,
    "lastOne": lastOne,
    "translateKeyword": translateKeyword,
    "translatePath": translatePath,
    "translateHost": translateHost,
    "translateDate": translateDate,
  }

  tmpl, err := template.New("notSureWhatThisDoes").Funcs(funcMap).ParseFiles(files...)
  if err != nil {
    return nil, err
  }

  return tmpl, nil
}

func fetchData(host string, path string, postQuant int, tagFilter string) (map[string]interface{}, error) {
  var err error
  lang := fetchLang(host)
  data := make(map[string]interface{})

  data["Lang"] = lang
  data["Domain"] = host
  data["Scheme"] = scheme 
  data["Path"] = path
  data["Posts"], err = aggregatePosts(postQuant, tagFilter)
  if err != nil {
    return data, err
  }

  if strings.HasPrefix(path, translatePath(lang, "/posts/")) && len(path) > len(translatePath(lang, "/posts/")) {
    data["Post"], err = fetchPost(strings.TrimPrefix(path, translatePath(lang, "/posts/")))
    if err != nil {
      return data, err
    }
  }

  if path == translatePath(lang, "/about") {
    data["Song"], data["TrackIndex"] = rockNRoll()
  }

  return data, nil
}

func serveTMPL(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data map[string]interface{}) {
  var buf bytes.Buffer
  err := tmpl.ExecuteTemplate(&buf, "base", data)
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(w, r, http.StatusInternalServerError)
    return
  }
  buf.WriteTo(w)
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  translatedURL := translatePath(fetchLang(r.Host), r.URL.Path)
  if r.URL.Path != translatedURL {
    http.Redirect(w, r, translatedURL, 302)
    return
  }

  path := strings.Split(r.URL.Path, "/")
  page := translateKeyword("en-US", path[1])
  if r.URL.Path == "/" {
    page = "index"
  } else if len(path) == 3 && path[2] == "" {
    http.Redirect(w, r, "/" + page, 302)
  } else if len(path) > 2 {
    fancyErrorHandler(w, r, http.StatusNotFound)
    return
  }

  if !doesFileExist(filepath.Join(htmlDir, "pages", page + tmplFileExt)) {
    fancyErrorHandler(w, r, http.StatusNotFound)
    return
  }

  tmpl, err := bindTMPL(
    filepath.Join(htmlDir, "base" + tmplFileExt),
    filepath.Join(htmlDir, "pages", page + tmplFileExt),
  )
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(w, r, http.StatusInternalServerError)
    return
  }

  var data map[string]interface{}
  switch translatePath("en-US", r.URL.Path)  {
  case "/":
    data, err = fetchData(r.Host, r.URL.Path, 3, "articles")
  case "/posts": 
    data, err = fetchData(r.Host, r.URL.Path, 0, "")
  default:
    data, err = fetchData(r.Host, r.URL.Path, -1, "")
  } 
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(w, r, http.StatusInternalServerError)
    return
  }

  serveTMPL(w, r, tmpl, data)
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  
  path := strings.Split(r.URL.Path, "/")
  lang := fetchLang(r.Host)

  // example.org/tags/ -> example.org/posts
  if len(path) == 3 && path[2] == "" {
    http.Redirect(w, r, translatePath(lang, "/posts"), 302)
    return
  }
  tag := translateKeyword("en-US", path[2])

  // de.example.org/tags/photos -> de.example.org/stichwoerter/fotos
  // example.org/tags/tag1/nonsense -> example.org/tags/tag1
  if r.URL.Path != translatePath(lang, r.URL.Path) || len(path) > 3 {
    http.Redirect(w, r, translatePath(lang, "/tags/" + tag), 302)
    return
  }
  
  if !doesFileExist(filepath.Join(htmlDir, "tags", tag + tmplFileExt)) {
    fancyErrorHandler(w, r, http.StatusNotFound)
    return
  }

  tmpl, err := bindTMPL(
    filepath.Join(htmlDir, "base" + tmplFileExt),
    filepath.Join(htmlDir, "tags", tag + tmplFileExt), 
  )
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(w, r, http.StatusInternalServerError)
    return
  }

  data, err := fetchData(r.Host, r.URL.Path, 0, tag)
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(w, r, http.StatusInternalServerError)
    return
  }

  serveTMPL(w, r, tmpl, data)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  lang := fetchLang(r.Host)

  path := strings.Split(r.URL.Path, "/")
  post := path[2]
  // example.org/posts/ -> example.org/posts
  if len(path) == 3 && path[2] == "" {
    http.Redirect(w, r, translatePath(lang, "/posts"), 302)
    return
  }

  // de.example.org/entradas/post1 -> de.example.org/posten/post1
  // example.org/posts/post1/nonsense -> example.org/posts/post1
  if r.URL.Path != translatePath(lang, r.URL.Path) || len(path) > 3 {
    http.Redirect(w, r, translatePath(lang, "/posts/") + post, 302)
    return
  }

  if !doesFileExist(filepath.Join(htmlDir, "posts", post + tmplFileExt)) {
    fancyErrorHandler(w, r, http.StatusNotFound)
    return
  }

  tmpl, err := bindTMPL(
    filepath.Join(htmlDir, "base" + tmplFileExt),
    filepath.Join(htmlDir, "partials", "post_header" + tmplFileExt),
    filepath.Join(htmlDir, "posts", post + tmplFileExt),
  )
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(w, r, http.StatusInternalServerError)
    return
  }

  data, err := fetchData(r.Host, r.URL.Path, -1, "")
  if err != nil {
    log.Println(err.Error())
    fancyErrorHandler(w, r, http.StatusInternalServerError)
    return
  }

  serveTMPL(w, r, tmpl, data)
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
  posts, err := aggregatePosts(0, "")
  if err != nil {
    fancyErrorHandler(w, r, http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/atom+xml")
  feed := bytes.NewReader(generateFeed(posts))
  http.ServeContent(w, r, "atom.xml", time.Now(), feed)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(staticDir,"favicon.ico"))
}
