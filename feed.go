package main

import (
	"time"

	"angel-castaneda.com/sqlite"
)

// use this to check for valid feed: https://validator.w3.org/feed/
func generateFeed(posts []sqlite.Post) []byte {
  domain := "angel-castaneda.com"
  feed := `<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom">
  <title>angel's site</title>
  <subtitle>my personal site's feed</subtitle>
  <link rel="self" href="https://www.` + domain + `/atom.xml"/>
  <link href="https://www.` + domain + `"/>
  <author>
    <name>Angel Castaneda</name>
  </author>
  <id>https://www.` + domain + `</id>
  <updated>` + time.Now().UTC().Format("2006-01-02T15:04:05.000Z") + `</updated>`

  for _, post := range posts {
    entry := `
  <entry>
    <title>` + post.Title + `</title>
    <link href="https://www.` + domain + `/posts/` + post.FileName + `"/>
    <id>https://www.` + domain + `/posts/` + post.FileName + `</id>
    <published>` + post.PubDate + `T00:00:00.000Z</published>
    <updated>` + post.UpdateDate + `T00:00:00.000Z</updated>
    <summary>` + post.Description + `</summary>`
    for _, tag := range post.Tags {
      category := `
    <category term="` + tag.Name + `" scheme="https://www.` + domain + `/posts"/>`
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
