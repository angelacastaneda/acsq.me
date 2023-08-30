package main

import (
	"log"
	"time"

	"angel-castaneda.com/sqlite"
)

// use this to check for valid feed: https://validator.w3.org/feed/
func generateFeed(domain string, posts []sqlite.Post) []byte {
	feed := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>angel's site</title>
  <subtitle>my personal site's feed</subtitle>
  <link rel="self" href="https://` + domain + `/atom.xml"/>
  <link href="https://` + domain + `"/>
  <author>
    <name>Angel Castaneda</name>
  </author>
  <id>https://` + domain + `</id>
  <updated>` + time.Now().UTC().Format("2006-01-02T15:04:05.000Z") + `</updated>`

	for _, post := range posts {
		entry := `
  <entry>
    <title>` + post.Title + `</title>
    <link href="https://` + domain + `/posts/` + post.FileName + `"/>
    <id>https://` + domain + `/posts/` + post.FileName + `</id>
    <published>` + post.PubDate + `T00:00:00.000Z</published>
    <updated>` + post.UpdateDate + `T00:00:00.000Z</updated>
    <summary>` + post.Description + `</summary>`
		p, err := sqlite.FetchPost(post.FileName)
		if err != nil {
			log.Println(err.Error())
			break // todo undo this awful hack solution
		}
		for _, t := range p.Tags {
			category := `
    <category term="` + t.Name + `" scheme="https://` + domain + `/tags/` + t.Name + `"/>`
			entry = entry + category
		}
		entry = entry + `
  </entry>`
		feed = feed + entry
	}
	feed = feed + `
</feed>`

	return []byte(feed)
}
