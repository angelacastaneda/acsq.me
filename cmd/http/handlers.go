package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"acsq.me/dblog"
)

var (
	htmlDir   = filepath.Join(".", "html") // routes to dirs
	staticDir = filepath.Join(".", "static")
)

const (
	tmplFileExt = ".tmpl.html"
)

func fancyErrorHandler(w http.ResponseWriter, r *http.Request, httpCode int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(httpCode)

	tmpl, err := bindTMPL(
		filepath.Join(htmlDir, "base"+tmplFileExt),
		filepath.Join(htmlDir, "partials", "error_meta"+tmplFileExt),
		filepath.Join(htmlDir, "partials", "error_header"+tmplFileExt),
		filepath.Join(htmlDir, "errors", strconv.Itoa(httpCode)+tmplFileExt),
	)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := fetchBaseData(r.Host, "/")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	return
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	// first step is to clean the url
	path := strings.SplitN(r.URL.Path, "/", 3) // TODO turn this into middleware
	if len(path) == 3 {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// then see if the page exists
	page := path[1]
	page = strings.TrimSuffix(page, ".html")
	page = translateKeyword("en-US", page)
	if r.URL.Path == "/" {
		page = "index"
	}
	if !doesFileExist(filepath.Join(htmlDir, "pages", page+tmplFileExt)) {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// then redirect to correct lang
	lang := fetchLang(r.Host)
	translatedURL := translatePath(lang, r.URL.Path)
	if r.URL.Path != translatedURL {
		http.Redirect(w, r, translatedURL, 302)
		return
	}

	// then redirect to correct ending
	if !strings.HasSuffix(r.URL.Path, ".html") && r.URL.Path != "/" {
		http.Redirect(w, r, r.URL.Path+".html", 302)
		return
	}

	// now start building page
	tmpl, err := bindTMPL(
		filepath.Join(htmlDir, "base"+tmplFileExt),
		filepath.Join(htmlDir, "pages", page+tmplFileExt),
	)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	data := fetchBaseData(r.Host, r.URL.Path)
	switch page {
	case "index":
		data["Posts"], err = dblog.AggregatePosts(3, "articles")
		if err != nil { // TODO consider goto for error handling
			log.Println(err.Error())
			fancyErrorHandler(w, r, http.StatusInternalServerError)
			return
		}
		data["Post"], err = dblog.FetchThumbnail()
	case "posts":
		data["Posts"], err = dblog.AggregatePosts(0, "")
	default:
		data["Song"], data["TrackIndex"] = rockNRoll()
	}
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	serveTMPL(w, r, tmpl, data)
	return
}

func tagHandler(w http.ResponseWriter, r *http.Request) {
	// first step is to clean the url
	path := strings.SplitN(r.URL.Path, "/", 4) // TODO turn this into middleware
	if len(path) != 3 {
		log.Println("too long path", path)
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// then see if tag exists
	lang := fetchLang(r.Host)
	tag := strings.TrimSuffix(path[2], ".html")
	tag = translateKeyword("en-US", tag)

	if !dblog.DoesTagExist(tag) {
		log.Println("not in db", tag)
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// then redirect to correct lang
	// de.example.org/tags/photos.html -> de.example.org/stichwoerter/fotos.html
	// example.org/tags/photos -> example.org/tags/photos.html
	if r.URL.Path != translatePath(lang, r.URL.Path) || !strings.HasSuffix(r.URL.Path, ".html") {
		http.Redirect(w, r, translatePath(lang, "/tags/"+tag+".html"), 302)
		return
	}

	// then build page
	tagData, err := dblog.FetchTag(tag)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	tmpl, err := sqlBindTMPL(tagData.Description,
		filepath.Join(htmlDir, "base"+tmplFileExt),
		filepath.Join(htmlDir, "blog", "tag"+tmplFileExt),
	)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	data := fetchBaseData(r.Host, r.URL.Path)
	data["Posts"], err = dblog.AggregatePosts(0, tag)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	data["Tag"], err = dblog.FetchTag(tag)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	serveTMPL(w, r, tmpl, data)
	return
}

func postDateRedirect(w http.ResponseWriter, r *http.Request) {
	// first step is to clean the url
	path := strings.SplitN(r.URL.Path, "/", 4) // TODO turn this into middleware
	if len(path) != 3 {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// then see if the post exists
	post := path[len(path)-1]
	post = strings.TrimSuffix(post, ".html")
	if !dblog.DoesPostExist(post) {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// then redirect to the correct date
	postData, err := dblog.FetchPost(post)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	date := postData.PubDate
	year := date[:4]
	month := date[5:7]
	day := date[8:]
	http.Redirect(w, r, translatePath(fetchLang(r.URL.Host), "/posts/")+year+"/"+month+"/"+day+"/"+post+".html", 302)
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	// first step is to clean the url
	path := strings.SplitN(r.URL.Path, "/", 7) // TODO turn this into middleware
	if len(path) != 6 {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// then see if the post exists
	post := path[len(path)-1]
	post = strings.TrimSuffix(post, ".html")
	if !dblog.DoesPostExist(post) {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// and it has the correct date
	year := r.PathValue("year")
	month := r.PathValue("month")
	day := r.PathValue("day")
	postData, err := dblog.FetchPost(post)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	if year+"-"+month+"-"+day != postData.PubDate {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// then redirect to correct lang
	lang := fetchLang(r.Host)
	// de.example.org/entradas/cool-post.html -> de.example.org/posten/cool-post.html
	// example.org/posts/cool-post -> example.org/posts/cool-posts.html
	if r.URL.Path != translatePath(lang, r.URL.Path) || !strings.HasSuffix(r.URL.Path, ".html") {
		http.Redirect(w, r, translatePath(lang, "/posts/")+year+"/"+month+"/"+day+"/"+post+".html", 302)
		return
	}

	// then build page
	tmpl, err := sqlBindTMPL(postData.Content,
		filepath.Join(htmlDir, "base"+tmplFileExt),
		filepath.Join(htmlDir, "partials", "post_header"+tmplFileExt),
		filepath.Join(htmlDir, "partials", "katex"+tmplFileExt), // todo make it check if it's a math article first
		filepath.Join(htmlDir, "blog", "post"+tmplFileExt),
	)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	data := fetchBaseData(r.Host, r.URL.Path)
	data["Post"] = postData
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	serveTMPL(w, r, tmpl, data)
	return
}

func feedHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := dblog.AggregatePosts(0, "")
	if err != nil {
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/atom+xml")
	feed := bytes.NewReader(generateFeed(r.Host, posts))
	http.ServeContent(w, r, "atom.xml", time.Now(), feed)
}

func cvHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(staticDir, "files", "cv.pdf"))
}

func pgpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	http.ServeFile(w, r, filepath.Join(staticDir, "files", "angelcastaneda.asc"))
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(staticDir, "favicon.ico"))
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Println(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		log.Println(http.StatusBadRequest)
		return
	}

	recommendation := struct {
		Name   string `json:"name"`
		Title  string `json:"title"`
		Author string `json:"author"`
		Note   string `json:"note"`
	}{
		Name:   r.FormValue("recommender"),
		Title:  r.FormValue("title"),
		Author: r.FormValue("author"),
		Note:   r.FormValue("note"),
	}

	jsonBytes, err := json.Marshal(recommendation)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(string(jsonBytes))

	params := url.Values{}
	params.Add("title", recommendation.Title)
	params.Add("recommender", recommendation.Name)
	http.Redirect(w, r, "/recommend?"+params.Encode(), http.StatusSeeOther)
}

func redirectWithParams(params url.Values, w http.ResponseWriter, r *http.Request, url string, code int) {
	http.Redirect(w, r, url+"?"+params.Encode(), code)
}

func recommendHandler(w http.ResponseWriter, r *http.Request) {
	// first step is to clean the url
	path := strings.SplitN(r.URL.Path, "/", 3) // TODO turn this into middleware
	if len(path) == 3 {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// then check if page exists
	if !doesFileExist(filepath.Join(htmlDir, "recommend"+tmplFileExt)) {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// then finally you can translate url itself
	lang := fetchLang(r.Host)
	translatedURL := translatePath(lang, r.URL.Path)
	if !strings.HasSuffix(translatedURL, ".html") {
		translatedURL += ".html"
	}
	if r.URL.Path != translatedURL {
		redirectWithParams(r.URL.Query(), w, r, translatedURL, 302)
		return
	}

	tmpl, err := bindTMPL(
		filepath.Join(htmlDir, "base"+tmplFileExt),
		filepath.Join(htmlDir, "recommend"+tmplFileExt),
	)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	data := fetchBaseData(r.Host, r.URL.RequestURI())
	Rec := struct {
		Title       string
		Recommender string
	}{
		Title:       r.URL.Query().Get("title"),
		Recommender: r.URL.Query().Get("recommender"),
	}
	data["Rec"] = Rec
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}

	serveTMPL(w, r, tmpl, data)
	return
}
