package main
import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
  pathToDB = "./posts.sqlite3"
)

type Post struct {
  Title string
  FileName string
  Content string
  Description string
  PubDate string
  UpdateDate string
  Tags []Tag
  Thumbnail Thumbnail
}

type Thumbnail struct {
  Img string `json:"img"`
  Alt string `json:"alt"`
  Title string `json:"title"`
}

type Tag struct {
  Name string
  Category string
  Description string
}

func makeDB() {
  db, err := sql.Open("sqlite3", pathToDB)
  if err != nil {
    log.Fatal(err.Error())
  }
  defer db.Close()

  // post table
  _, err = db.Exec(`CREATE TABLE IF NOT EXISTS posts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL UNIQUE,
  file_name TEXT NOT NULL UNIQUE,
  content TEXT NOT NULL,
  description TEXT NOT NULL,
  pub_date TEXT NOT NULL CHECK(pub_date LIKE '____-__-__'),
  update_date TEXT CHECK(update_date LIKE '____-__-__'),
  thumbnail TEXT -- in json format but go engine can't handle real json
);`)
  if err != nil {
    log.Println(err.Error())
  }

  // tag table
  _, err = db.Exec(`CREATE TABLE IF NOT EXISTS tags (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  category STRING NOT NULL DEFAULT 'content', -- for medium, content, and lang
  description TEXT
);`)
  if err != nil {
    log.Println(err.Error())
  }

  // associative identity
  _, err = db.Exec(`CREATE TABLE IF NOT EXISTS posts_tags (
  post_id INTEGER,
  tag_id INTEGER,
  PRIMARY KEY (post_id, tag_id),
  FOREIGN KEY (post_id) REFERENCES posts(id),
  FOREIGN KEY (tag_id) REFERENCES tags(id)
);`)
  if err != nil {
    log.Println(err.Error())
  }

}

func aggregatePosts(postQty int, filterTag string) (posts []Post, err error){
  if postQty < 0 {
    return []Post{}, nil
  }
  
  db, err := sql.Open("sqlite3", pathToDB)
  if err != nil {
    log.Fatal(err.Error())
  }
  defer db.Close()

  var query string
  var filters []interface{}
  if filterTag == "" {
    query = `SELECT title, file_name, pub_date, update_date
  FROM posts
  ORDER BY pub_date DESC`
  } else {
    query = `SELECT title, file_name, pub_date, update_date
  FROM posts JOIN posts_tags
  ON posts.id = posts_tags.post_id JOIN tags
  ON posts_tags.tag_id = tags.id
  WHERE tags.name = ?
  ORDER BY pub_date DESC`
    filters = append(filters, filterTag)
  }

  if postQty > 0 {
    query = query + `
  LIMIT ?`
    filters = append(filters, postQty)
  }

  rows, err := db.Query(query, filters...)
  defer rows.Close()

  for rows.Next() {
    post := Post{}
    var update_date sql.NullString
    err := rows.Scan(&post.Title, &post.FileName, &post.PubDate, &update_date)
    if err != nil {
      return posts, err
    }
    if update_date.Valid {
      post.UpdateDate = update_date.String
    } else {
      post.UpdateDate = ""
    }
    posts = append(posts, post)
  }

  return posts, nil
}

func fetchPost(fileNameNoExtension string) (post Post, err error) {
  db, err := sql.Open("sqlite3", pathToDB)
  if err != nil {
    log.Fatal(err.Error())
  }
  defer db.Close()

  var id int
  var update_date, thumbnailJSON sql.NullString
  err = db.QueryRow(`SELECT id, title, file_name, content, description, pub_date, update_date, thumbnail
  FROM posts
  WHERE file_name = ?`, fileNameNoExtension).Scan(&id, &post.Title, &post.FileName, &post.Content, &post.Description, &post.PubDate, &update_date, &thumbnailJSON)
  if err != nil {
    return Post{}, err
  }

  tagRows, err := db.Query(`SELECT tags.name
  FROM tags JOIN posts_tags
  ON tags.id = posts_tags.tag_id
  WHERE posts_tags.post_id = ?
  ORDER BY name`, id)
  if err != nil {
    return Post{}, err
  }
  defer tagRows.Close()

  var tags []Tag 
  for tagRows.Next() {
    var name string
    err := tagRows.Scan(&name)
    if err != nil {
      log.Println(err.Error())
      continue
    }

    tags = append(tags, Tag{Name: name})
  }
  post.Tags = tags

  // optional stuff
  if update_date.Valid {
    post.UpdateDate = update_date.String
  }
  if thumbnailJSON.Valid {
    // e.g. { "img" : "pic.jpeg", "alt" : "cool pic", "title" : "what you see if you hover"}
    var thumbnail Thumbnail
    err := json.Unmarshal([]byte(thumbnailJSON.String), &thumbnail)
    if err != nil {
      log.Println(err.Error())
    } else {
      post.Thumbnail = thumbnail
    }
  }

  return post, nil
}

func fetchTag(tagName string) (tag Tag, err error) {
  db, err := sql.Open("sqlite3", pathToDB)
  if err != nil {
    log.Fatal(err.Error())
  }
  defer db.Close()

  err = db.QueryRow(`SELECT name, description, category
  FROM tags
  WHERE name = ?`, tagName).Scan(&tag.Name, &tag.Description, &tag.Category)
  if err != nil {
    return Tag{}, err
  }

  return tag, nil
}

func doesPostExist(fileNameNoExtension string) bool {
  db, err := sql.Open("sqlite3", pathToDB)
  if err != nil {
    log.Fatal(err.Error())
  }
  defer db.Close()

  var count int
  err = db.QueryRow(`SELECT COUNT(*)
  FROM posts
  WHERE file_name = ?`, fileNameNoExtension).Scan(&count)
  if err != nil {
    log.Println(err.Error())
    return false
  }

  return count > 0
}

func doesTagExist(tag string) bool {
  db, err := sql.Open("sqlite3", pathToDB)
  if err != nil {
    log.Fatal(err.Error())
  }
  defer db.Close()

  var count int
  err = db.QueryRow(`SELECT COUNT(*)
  FROM tags
  WHERE name = ?`, tag).Scan(&count)
  if err != nil {
    log.Println(err.Error())
    return false
  }

  return count > 0
}

func addPost(post Post) bool {
  // connect to db and add post

  return false
}

func addTag(tag Tag) bool {
  // connect to db and add tag

  return false
}
