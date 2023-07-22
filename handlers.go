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
  domain = "localhost:4000" // todo fix this 
  scheme = "http"
)

func redirectHTTPS(w http.ResponseWriter, r *http.Request) {
  target := "https://" + r.Host + r.URL.Path // todo get actual raw path too
  http.Redirect(w, r, target, 302)
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

func fancyErrorHandler(httpCode int, w http.ResponseWriter, r *http.Request) {
  // w.Header().Set("Content-Type","text/html; charset=utf-8")
  w.WriteHeader(httpCode)

  tmpl, err := bindTMPL(
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

  data, err := fetchData(r, -404, "")
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

func bindTMPL(files ...string) (*template.Template, error) {
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
    "translateKeyword": translateKeyword,
    "translateURL": translateURL,
    "translateDate": translateDate,
  }

  tmpl, err := template.New("noIdeaWhyThisExists").Funcs(funcMap).ParseFiles(files...)
  if err != nil {
    return nil, err
  }

  return tmpl, nil
}

func fetchData(r *http.Request, postQuant int, tagFilter string) (map[string]interface{}, error) {
  var err error
  lang := fetchLang(r.Host)
  data := make(map[string]interface{})

  data["Lang"] = lang
  data["Domain"] = domain 
  data["Scheme"] = scheme 
  data["Path"] = r.URL.Path

  if postQuant == -404 { // todo undo this hack error filter
    data["Path"] = "/"
  } 

  data["Posts"], err = aggregatePosts(postQuant, tagFilter)
  if err != nil {
    return data, err
  }

  if strings.HasPrefix(r.URL.Path, translateURL(lang, "/posts/")) && len(r.URL.Path) > len(translateURL(lang, "/posts/")) && postQuant != -404 {
    data["Post"], err = fetchPost(strings.TrimPrefix(r.URL.Path, translateURL(lang, "/posts/")))
    if err != nil {
      return data, err
    }
  }

  if r.URL.Path == translateURL(lang, "/about") {
    data["Song"], data["TrackIndex"] = rockNRoll()
  }

  return data, nil
}

func serveTMPL(w http.ResponseWriter, r *http.Request, tmpl *template.Template, postQuant int, tagFilter string) {

  data, err := fetchData(r, postQuant, tagFilter)
  if err != nil {
    log.Println(err.Error())
    // fancyErrorHandler(http.StatusInternalServerError, w, r)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }

  err = tmpl.ExecuteTemplate(w, "base", data)
  if err != nil {
    log.Println(err.Error())
    // fancyErrorHandler(http.StatusInternalServerError, w, r)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  translatedURL := translateURL(fetchLang(r.Host), r.URL.Path)

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
    fancyErrorHandler(http.StatusNotFound, w, r)
    // http.Error(w,"Page Not Found", http.StatusNotFound)
    return
  }

  if !doesFileExist(filepath.Join(htmlDir, "pages", page + tmplFileExt)) {
    fancyErrorHandler(http.StatusNotFound, w, r)
    // http.Error(w,"Page Not Found", http.StatusNotFound)
    return
  }

  tmpl, err := bindTMPL(
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
    serveTMPL(w, r, tmpl, 3, "articles")
    return
  case translateURL("en-US", "/posts"), translateURL("es-US", "/posts"), translateURL("de-DE", "/posts"): // todo make this less ugly
    serveTMPL(w, r, tmpl, 0, "")
    return
  default:
    serveTMPL(w, r, tmpl, -1, "")
  } 
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")
  
  lang := fetchLang(r.Host)

  // example.org/tags/ -> example.org/posts
  if r.URL.Path == translateURL("en-US", "/tags/") ||
  r.URL.Path == translateURL("es-US", "/tags/") ||
  r.URL.Path == translateURL("de-DE", "/tags/") {
    http.Redirect(w, r, translateURL(lang, "/posts"), 302)
    return
  }

  urlPath := strings.Split(r.URL.Path, "/")

  // de.example.org/tags/photos -> de.example.org/stichwoerter/fotos
  // example.org/tags/tag1/nonsense -> example.org/tags/tag1
  if r.URL.Path != translateURL(lang, r.URL.Path) || len(urlPath) > 3 {
    http.Redirect(w, r, translateURL(lang, "/tags/" + urlPath[2]), 302)
    return
  }
  
  tag := translateKeyword("en-US", urlPath[2])
  if !doesFileExist(filepath.Join(htmlDir, "tags", tag + tmplFileExt)) {
    fancyErrorHandler(http.StatusNotFound, w, r)
    // http.Error(w,"Page Not Found", http.StatusNotFound)
    return
  }

  tmpl, err := bindTMPL(
    filepath.Join(htmlDir, "partials", "meta" + tmplFileExt),
    filepath.Join(htmlDir, "partials", "header" + tmplFileExt),
    filepath.Join(htmlDir, "tags", tag + tmplFileExt), 
    filepath.Join(htmlDir, "partials", "footer" + tmplFileExt),
  )
  if err != nil {
    log.Println(err.Error())
    // fancyErrorHandler(http.StatusInternalServerError, w, r)
    http.Error(w,"Internal Server Error", http.StatusInternalServerError)
    return
  }

  serveTMPL(w, r, tmpl, 0, tag)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type","text/html; charset=utf-8")

  lang := fetchLang(r.Host)

  // example.org/posts/ -> example.org/posts
  if r.URL.Path == translateURL("en-US", "/posts/") ||
  r.URL.Path == translateURL("es-US", "/posts/") ||
  r.URL.Path == translateURL("de-DE", "/posts/") {
    http.Redirect(w, r, translateURL(lang, "/posts"), 302)
    return
  }

  urlPath := strings.Split(r.URL.Path, "/")
  postsRoot := urlPath[1]
  post := urlPath[2]

  // de.example.org/entradas/post1 -> de.example.org/posten/post1
  // example.org/posts/post1/nonsense -> example.org/posts/post1
  if postsRoot != translateKeyword(lang, "posts") || len(urlPath) > 3 {
    http.Redirect(w, r, translateURL(lang, "/posts/") + post, 302)
    return
  }


  if !doesFileExist(filepath.Join(htmlDir, "posts", post + tmplFileExt)) {
    fancyErrorHandler(http.StatusNotFound, w, r)
    // http.Error(w,"Page Not Found", http.StatusNotFound)
    return
  }

  tmpl, err := bindTMPL(
    filepath.Join(htmlDir, "partials", "meta" + tmplFileExt),
    filepath.Join(htmlDir, "partials", "post_header" + tmplFileExt),
    filepath.Join(htmlDir, "posts", post + tmplFileExt),
    filepath.Join(htmlDir, "partials", "footer" + tmplFileExt),
  )
  if err != nil {
    log.Println(err.Error())
    // fancyErrorHandler(http.StatusInternalServerError, w, r)
    http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    return
  }

  serveTMPL(w, r, tmpl, -1, "")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(staticDir,"favicon.ico"))
}
