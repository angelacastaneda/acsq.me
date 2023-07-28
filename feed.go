package main

import (
	"time"
)

// use this to check for valid feed: https://validator.w3.org/feed/
func generateFeed(posts []Post) []byte {
  domain := "angel-castaneda.com" // todo abstract more
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

  for _, post := range posts { // todo add content and summary
    entry := `
  <entry>
    <title>` + post.Title + `</title>
    <link href="https://www.` + domain + `/posts/` + post.FileName + `"/>
    <id>https://www.` + domain + `/posts/` + post.FileName + `</id>
    <published>` + post.Date + `T00:00:00.000Z</published>
    <updated>` + post.Date + `T00:00:00.000Z</updated>` // todo actually add updating system.
    for _, tag := range post.Tags {
      category := `
    <category term="` + tag + `" scheme="https://www.` + domain + `/posts"/>`
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
