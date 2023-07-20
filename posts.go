package main

import (
	// "errors"
	"html/template"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
  postDir = filepath.Join(".","html","posts")
  tmplFileExt = ".tmpl.html"
)

type Post struct {
  FileName string
  Title string
  Date string
  Tags []string
}

func (p Post) containsTag(filterTag string) bool {
  if filterTag == "" {
    return true
  }

  for _, tag := range p.Tags {
    if filterTag == tag {
      return true
    }
  }

  return false
}

func postFetcher(postNameNoExt string) (Post, error) {

  // reading file
  content, err := os.ReadFile(filepath.Join(postDir, postNameNoExt + tmplFileExt))
  if err != nil {
    return Post{}, err
  }

  // getting title
  formattedFileName := strings.ReplaceAll(postNameNoExt, "_", " ")

  // getting tags
  tagsPattern := regexp.MustCompile(`{{define "keywords"}}([\w\s]+){{end}}`) 
  tagsMatching := tagsPattern.FindStringSubmatch(string(content)) 
  
  var tags []string
  if len(tagsMatching) > 1 {
    tags = strings.Fields(tagsMatching[1]) 
  } else {
    tags = []string{}
  }

  // sort tags alphabetically
  sort.Slice(tags, func(i, j int) bool {
    return tags[i] < tags[j]
  })

  // getting date (as if)
  datePattern := regexp.MustCompile(`<time datetime="(\d{4}-\d{2}-\d{2})">`)   
  dateMatching := datePattern.FindStringSubmatch(string(content))

  var date string
  if len(dateMatching) > 1 {
    date = dateMatching[1]
  } else {
    date = ""
  }

  // putting everything together
  return Post{
    FileName: postNameNoExt,
    Title: formattedFileName,
    Date: date,
    Tags: tags,
  }, nil
}

func postsSorter(postQuant int, filterTag string) ([]Post, error) {
  posts := []Post{}

  if postQuant < 0 {
    return posts, nil
  }

  // read dir
  files, err := os.ReadDir(postDir)
  if err != nil {
    return []Post{}, err
  }

  for _, file := range files {
    if !file.IsDir() && strings.HasSuffix(file.Name(),tmplFileExt) {
      fileNameNoExt := strings.TrimSuffix(file.Name(),tmplFileExt)
      newPost, err := postFetcher(fileNameNoExt)
      if err != nil {
        log.Println(err.Error())
        continue // todo make better error handling to know something went wrong
      }

      // filtering by tag
      if !newPost.containsTag(filterTag) {
        continue
      }
    
      posts = append(posts, newPost)
    }
  }

  sort.Slice(posts, func(i, j int) bool {
    return posts[i].Date > posts[j].Date
  })

  // if len(posts) == 0 { // todo make error for no html when there's enough content
  //   return posts, errors.New("Couldn't find " + tmplFileExt + " files in " + postDir + " :(")
  // }
  
  if postQuant < len(posts) && postQuant > 0 {
    return posts[:postQuant], nil
  } else {
    return posts, nil
  }
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
  }
  trackIndex := rand.Intn(len(awesomeTunes))
  track := awesomeTunes[trackIndex]

  return track, trackIndex
}

func langFetcher(url string) string {
  if strings.HasPrefix(url, "es.") {
    return "es-US"
  }
  
  if strings.HasPrefix(url, "de.") {
    return "de-DE"
  }

  return "en-US"
}

func lastOne(index int, size int) bool{
  return index == size - 1
}

func translate(lang, en, es, de string) template.HTML {
  switch lang {
  case "es-US":
    return template.HTML(es)
  case "de-DE":
    return template.HTML(de)
  default:
    return template.HTML(en)
  }
}

func translateTag(lang, tagName string) string {
  
  tagDictionary := map[string]map[string]string{
    // medium
    "articles": {
      "en-US": "articles",
      "es-US": "artículos",
      "de-DE": "artikel",
    },
    "photos": {
      "en-US": "photos",
      "es-US": "fotos",
      "de-DE": "fotos",
    },
    // lang
    "english": {
      "en-US": "english",
      "es-US": "inglés",
      "de-DE": "englisch",
    },
    "spanish": {
      "en-US": "spanish",
      "es-US": "español",
      "de-DE": "spanisch",
    },
    "german": {
      "en-US": "german",
      "es-US": "alemán",
      "de-DE": "deutsch",
    },
    // tags
    "math": {
      "en-US": "math",
      "es-US": "matemáticas",
      "de-DE": "mathe",
    },
    "milwaukee": {
      "en-US": "milwaukee",
      "es-US": "milwaukee",
      "de-DE": "milwaukee",
    },
    "history": {
      "en-US": "history",
      "es-US": "historia",
      "de-DE": "geschichte",
    },
    "technology": {
      "en-US": "technology",
      "es-US": "tecnologia",
      "de-DE": "technologie",
    },
    "personal": {
      "en-US": "personal",
      "es-US": "personal",
      "de-DE": "persönliches",
    },
  } 

  translation, ok := tagDictionary[tagName][lang] 

  if ok {
    return translation
  }

  return tagName
}

// type Translations struct {
//   EnUS string
//   EsUS string
//   DeDE string
// } // todo make this better

func translateURL(lang, originalURL string) string {
  // langURLs := map[string]Translations{
  //   "about": {
  //     EnUS: "about",
  //     EsUS: "conoceme",
  //     DeDE: "uber",
  //   },
  //   "posts": {
  //     EnUS: "posts",
  //     EsUS: "entradas",
  //     DeDE: "posten",
  //   },
  //   "friends": {
  //     EnUS: "friends",
  //     EsUS: "amigos",
  //     DeDE: "freunde",
  //   },
  //   "library": {
  //     EnUS: "library",
  //     EsUS: "biblioteca",
  //     DeDE: "bibliotek",
  //   },
  // }
  langURLs := map[string]map[string]string {
    "about": {
      "en-US": "about",
      "es-US": "conoceme",
      "de-DE": "uber",
    },
    "posts": {
      "en-US": "posts",
      "es-US": "entradas",
      "de-DE": "posten",
    },
    "friends": {
      "en-US": "friends",
      "es-US": "amigos",
      "de-DE": "freunde",
    },
    "library": {
      "en-US": "library",
      "es-US": "biblioteca",
      "de-DE": "bibliotek",
    },
  }

  switch originalURL {
  case "about","conoceme","uber":
    return langURLs["about"][lang]
  case "posts","entradas","posten":
    return langURLs["posts"][lang]
  case "friends","amigos","freunde":
    return langURLs["friends"][lang]
  case "library","biblioteca","bibliotek":
    return langURLs["library"][lang]
  }

  if langURLs[originalURL][lang] != "" {
    return langURLs[originalURL][lang]
  }

  return originalURL
}
