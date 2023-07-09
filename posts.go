package main

import (
	// "errors"
	"log"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strings"
)

const (
  postDir = "./html/posts"
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
  content, err := os.ReadFile(postDir + "/" + postNameNoExt + tmplFileExt)
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
  // read dir
  files, err := os.ReadDir(postDir)
  if err != nil {
    return []Post{}, err
  }

  posts := []Post{}

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
  
  if postQuant == 0 {
    return posts, nil
  } else if  postQuant < len(posts) {
    return posts[:postQuant], nil
  } else {
    return posts, nil
  }
}

func rockNRoll() (string, int) { // todo put this in a more sensible place
  awesomeTunes := []string{
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
  }
  rando := rand.Intn(len(awesomeTunes))
  shuffle := awesomeTunes[rando]

  return shuffle, rando
}
