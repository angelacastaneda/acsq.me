package main

import (
	"errors"
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

func postSorter(postQuant ...int) ([]Post, error) {
  files, err := os.ReadDir(postDir)
  if err != nil {
    return []Post{}, err
  }

  posts := []Post{}

  for _, file := range files {
    if !file.IsDir() && strings.HasSuffix(file.Name(),tmplFileExt) {

      // getting title
      fileName := strings.TrimSuffix(file.Name(),tmplFileExt)
      formattedName := strings.ReplaceAll(fileName, "_", " ")

      // reading file
      content, err := os.ReadFile(postDir + "/" + file.Name())
      if err != nil {
        return []Post{}, err
      }
      
      // getting date (as if)
      datePattern := regexp.MustCompile(`<time datetime="(\d{4}-\d{2}-\d{2})">`)   
      dateMatching := datePattern.FindStringSubmatch(string(content))

      var date string
      if len(dateMatching) > 1 {
        date = dateMatching[1]
      } else {
        date = ""
      }

      // getting tags
      tagsPattern := regexp.MustCompile(`{{define "keywords"}}([\w\s]+){{end}}`) 
      tagsMatching := tagsPattern.FindStringSubmatch(string(content)) 

      var tags []string
      if len(tagsMatching) > 1 {
        tags = strings.Fields(tagsMatching[1]) 
      } else {
        tags = []string{}
      }

      // putting everything together
      posts = append(posts, Post{
        FileName: fileName,
        Title: formattedName,
        Date: date,
        Tags: tags,
      })
    }
  }

  sort.Slice(posts, func(i, j int) bool {
    return posts[i].Date > posts[j].Date
  })

  if len(posts) == 0 {
    return posts, errors.New("Couldn't find " + tmplFileExt + " files in " + postDir + " :(")
  }
  
  if len(postQuant) == 0 {
    return posts, nil
  } else if len(postQuant) > 0 && postQuant[0] < len(posts) {
    return posts[:postQuant[0]], nil
  } else {
    return posts, nil
  }
}
