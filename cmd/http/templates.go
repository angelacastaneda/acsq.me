package main

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
)

func rockNRoll() (string, int) {
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
