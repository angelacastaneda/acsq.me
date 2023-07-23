package main

import (
	"encoding/xml"
	"time"
)

type Feed struct {
  XMLName xml.Name `xml:"feed"`
  Title string `xml:"title"`
  Subtitle string `xml:"subtitle"`
  Link []Link `xml:"link"`
  Updated string `xml:"updated"`
  Author Author `xml:"author"`
  ID string `xml:"id"`
  Entries []Entry `xml:"entry"`
}

type Author struct {
  Name string `xml:"name"`
}

type Link struct {
  Rel string `xml:"rel,attr"`
  Href string `xml:"href,attr"`
}

type Entry struct { // todo add content
  Title string `xml:"title"`
  Link Link `xml:"link"`
  ID string `xml:"id"`
  Published string `xml:"published"`
  Updated string `xml:"updated"`
  // Summary string `xml:"summary"`
  // Content Content `xml:"content"`
  Categories []Category `xml:"category"`
}

// type Content struct {

// }

type Category struct {
  Term string `xml:"term,attr"`
  Scheme string `xml:"scheme,attr"`
}

func generateFeed(posts []Post) (string, error) {
  myFeed := Feed{
    Title: "www.angel-castaneda.com",
    Subtitle: "my personal site's feed",
    Link: []Link{
      Link{Rel:"self", Href:"https://www.angel-castaneda.com/atom.xml"},
      Link{Href:"https://www.angel-castaneda.com"},
    },
    Updated: time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
    Author: Author{Name: "Angel Castaneda"},
    ID: "www.angel-castaneda.com",
  }

  for _, post := range posts {
    entry := Entry {
      Title: post.Title,
      Link: Link{Href:"www.angel-castaneda.com/posts/" + post.FileName},
      ID: "www.angel-castaneda.com/posts/" + post.FileName,
      Published: post.Date + "T00:00:00.000Z",
      Updated: post.Date + "T00:00:00.000Z", // todo actually impletement update feature
    }
    for _, tag := range post.Tags {
      entry.Categories = append(entry.Categories, Category {
        Term: tag,
        Scheme: "https://www.angel-castaneda.com/posts",
      })
    }
    myFeed.Entries = append(myFeed.Entries, entry)
  }

  xmlFeed, err := xml.MarshalIndent(myFeed, "", "  ")
  if err != nil {
    return "", err
  }

  return xml.Header + string(xmlFeed), nil
} 
