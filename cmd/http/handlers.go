package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
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

func rockNRoll() (string, int) { // todo put this in a more sensible place
	awesomeTunes := []string{
		// todo use something with less ads
		"https://youtu.be/ZV_UsQPTBy4", // "Sound and Vision" - David Bowie
		"https://youtu.be/GKdl-GCsNJ0", // "Here Comes the Sun" - The Beatles (duh)
		"https://youtu.be/ZVgHPSyEIqk", // "Let Down" - Radiohead
		"https://youtu.be/AZKch8dZ61w", // "St. Elmo's Fire" - Brian Eno
		"https://youtu.be/OP63BRzKmB0", // "Blade Runner (End Titles)" - Vanegelis
		"https://youtu.be/eLlmbCkb3As", // "Fallen Angel" - King Crimson
		"https://youtu.be/Hgx267jVma0", // "A Pillow of Winds" - Pink Floyd
		"https://youtu.be/vdvnOH060Qg", // "Happiness is a Warm Gun" - The Beatles (again)
		"https://youtu.be/Eo2ZsAOlvEM", // "America" - Simon and Garfunkel
		"https://youtu.be/fWB40wYQO-w", // "Dancing in My Head" - The Raincoats
		"https://youtu.be/GIrcy12Hruo", // "The Plains / Bitter Dancer" - Fleet Foxes
		"https://youtu.be/DMEOjFm4DJw", // "Cassius, -" - Fleet Foxes again cause I just saw their concert for a second time now
		"https://youtu.be/t_tIYlzSd2c", // "Bachelorette" - Björk
		"https://youtu.be/zG-q9Jozp4o", // "A New Kind of Water" - This Heat
		"https://youtu.be/X1GH9WN92s0", // "Another Green World" - Brian Eno
		"https://youtu.be/3GE-sfEbJ7I", // "Sheep" - Pink Floyd
		"https://youtu.be/dc6huqPzerY", // "Indiscipline" - King Crimson
		"https://youtu.be/95cufW4h-gA", // "One More Cup of Coffee" - Bob Dylan
		"https://youtu.be/i6d3yVq1Xtw", // "El Condor Pasa (If I Could)" - Simon and Garfunkel
		"https://youtu.be/OYmmthTXbSA", // "Stella Maris" - Einstürzende Neubauten
		"https://youtu.be/Y_V6y1ZCg_8", // "Norwegian Wood (This Bird Has Flow)" - The Beatles
		"https://youtu.be/LQ3nAhJyE44", // "Sunblind" - Fleet Foxes
		"https://youtu.be/K63CD2pwjD0", // "Wednesday Morning, 3 A.M." - Simon and Garfunkel
		"https://youtu.be/AtGEgxaO7nI", // "Alphabet Town" - Elliott Smith
		"https://youtu.be/NHDOk7lA53w", // "Ful Stop" - Radiohead
		"https://youtu.be/5ugdrdFrhI0", // "Nosferatu Man" - Slint
		"https://youtu.be/ojF9qAQ-8n4", // "Tangram Set 2" - Tangerine Dream
		"https://youtu.be/gl4lvJmvqQU", // "Happiness Is Easy" - Talk Talk
		"https://youtu.be/Ef9zt8aCRQo", // "Here Today" - The Beach Boys
		"https://youtu.be/sDcDCZGcZj8", // "Rocky Raccoon" - The Beatles
		"https://youtu.be/CHLQs6u9wXw", // "Here There and Everywhere" - The Beatles (best cover of the song)
		"https://youtu.be/ciLNMesqPh0", // "Vincent" - Don McLean
		"https://youtu.be/oFd9OhnKqvw", // "I Nearly Married A Human" - Tubeway Army
		"https://youtu.be/vwM77SSxLp8", // "Time It's Time" - Talk Talk (love this album)
		"https://youtu.be/zG-q9Jozp4o", // "A New Kind of Water" - Deceit
	}
	trackIndex := rand.Intn(len(awesomeTunes))
	track := awesomeTunes[trackIndex]

	return track, trackIndex
}

func doesFileExist(pathToFile string) bool {
	info, err := os.Stat(filepath.Clean(pathToFile))
	if err != nil || info.IsDir() {
		return false
	}
	return true
}

func getTempFuncs() (funcMap map[string]any) {
	funcMap = template.FuncMap{
		"lastOne":          lastOne,
		"translate":        translate,
		"translateHTML":    translateHTML,
		"translateKeyword": translateKeyword,
		"translatePath":    translatePath,
		"translateHost":    translateHost,
		"translateDate":    translateDate,
	}
	return funcMap
}

func bindTMPL(files ...string) (*template.Template, error) {
	for _, checkFile := range files {
		if !doesFileExist(checkFile) {
			return nil, errors.New("Template file missing " + checkFile)
		}
	}

	tmpl, err := template.New("notSureWhatThisDoes").Funcs(getTempFuncs()).ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func sqlBindTMPL(sqlContent string, files ...string) (*template.Template, error) {
	tmpl, err := bindTMPL(files...)
	if err != nil {
		return nil, err
	}

	sqlContent = `{{ define "sql" }}
` + sqlContent + `
{{ end }}`

	_, err = tmpl.New("sql").Parse(sqlContent)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func fetchBaseData(host string, path string) map[string]any {
	lang := fetchLang(host)
	data := make(map[string]any)
	email := translate(lang, "me@angelcastaneda.org", "yo@angelcastaneda.org", "ich@angelcastaneda.org")

	data["Lang"] = lang
	data["Scheme"] = scheme
	data["Path"] = path
	data["Domain"] = host
	data["Email"] = email

	return data
}

func serveTMPL(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data map[string]any) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "base", data)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = buf.WriteTo(w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
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

func postHandler(w http.ResponseWriter, r *http.Request) {
	// first step is to clean the url
	path := strings.SplitN(r.URL.Path, "/", 4) // TODO turn this into middleware
	if len(path) != 3 {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// then see if the post exists
	post := path[2]
	post = strings.TrimSuffix(post, ".html")
	if !dblog.DoesPostExist(post) {
		fancyErrorHandler(w, r, http.StatusNotFound)
		return
	}

	// then redirect to correct lang
	lang := fetchLang(r.Host)
	// de.example.org/entradas/cool-post.html -> de.example.org/posten/cool-post.html
	// example.org/posts/cool-post -> example.org/posts/cool-posts.html
	if r.URL.Path != translatePath(lang, r.URL.Path) || !strings.HasSuffix(r.URL.Path, ".html") {
		http.Redirect(w, r, translatePath(lang, "/posts/")+post+".html", 302)
		return
	}

	// then build page
	postData, err := dblog.FetchPost(post)
	if err != nil {
		log.Println(err.Error())
		fancyErrorHandler(w, r, http.StatusInternalServerError)
		return
	}
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
