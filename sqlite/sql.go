package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
  pathToDB = "./posts.sqlite3"
  ErrNotComplete = errors.New("unfilled required attributes")
)

type db struct {
  PathToDB string
}

type Post struct {
  Title string
  FileName string
  Content string
  Description string
  PubDate string
  UpdateDate string
  Tags []Tag
  Thumbnail Img
}

type Img struct {
  Src string `json:"src"`
  Alt string `json:"alt"`
  Title string `json:"title"`
}

type Tag struct {
  Name string
  Category string
  Description string
}

func OpenDB() (db *sql.DB) {
  db, err := sql.Open("sqlite3", pathToDB)
  if err != nil {
    log.Fatal(err.Error())
  }
  return db
}

func CloseDB(db *sql.DB) {
  if db != nil {
    db.Close()
  }
}

func MakeDB() (err error) {
  db := OpenDB()
  defer CloseDB(db)

  // post table
  if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS posts (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  title TEXT NOT NULL UNIQUE,
  file_name TEXT NOT NULL UNIQUE,
  content TEXT NOT NULL,
  description TEXT NOT NULL,
  pub_date TEXT NOT NULL CHECK(pub_date LIKE '____-__-__'),
  update_date TEXT NOT NULL CHECK(update_date LIKE '____-__-__'),
  thumbnail TEXT -- in json format but go engine can't handle real json
)`); err != nil {
    return err
  }

  // tag table
  if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS tags (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  category STRING NOT NULL DEFAULT 'content', -- for medium, content, and lang
  description TEXT
)`); err != nil {
    return err
  }

  // associative identity
  if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS posts_tags (
  post_id INTEGER,
  tag_id INTEGER,
  PRIMARY KEY (post_id, tag_id),
  FOREIGN KEY (post_id) REFERENCES posts(id),
  FOREIGN KEY (tag_id) REFERENCES tags(id)
)`); err != nil {
    return err
  }

  return nil
}

func AggregatePosts(postQty int, filterTag string) (posts []Post, err error){
  if postQty < 0 {
    return []Post{}, nil
  }

  db := OpenDB()
  defer CloseDB(db)

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

func FetchPost(fileName string) (post Post, err error) {
  if !DoesPostExist(fileName) {
    return Post{}, errors.New(fileName + " doesn't exist")
  }

  db := OpenDB()
  defer CloseDB(db)

  var id int
  var thumbnailJSON sql.NullString
  err = db.QueryRow(`SELECT id, title, file_name, content, description, pub_date, update_date, thumbnail
  FROM posts
  WHERE file_name = ?`, fileName).Scan(&id, &post.Title, &post.FileName, &post.Content, &post.Description, &post.PubDate, &post.UpdateDate, &thumbnailJSON)
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
  if thumbnailJSON.Valid && len(thumbnailJSON.String) > 0 {
    // e.g. { "src" : "pic.jpeg", "alt" : "cool pic", "title" : "what you see if you hover"}
    var thumbnail Img
    err := json.Unmarshal([]byte(thumbnailJSON.String), &thumbnail)
    if err != nil {
      log.Println(err.Error())
      return Post{}, err
    } else {
      post.Thumbnail = thumbnail
    }
  }

  return post, nil
}

func FetchThumbnail() (post Post, err error) {
  db := OpenDB()
  defer CloseDB(db)

  var thumbnailJSON sql.NullString
  err = db.QueryRow(`SELECT title, file_name, content, posts.description, pub_date, update_date, thumbnail
  FROM posts JOIN posts_tags
  ON posts.id = posts_tags.post_id JOIN tags
  ON posts_tags.tag_id = tags.id
  WHERE tags.name = 'photos'
  AND posts.thumbnail IS NOT NULL
  AND posts.thumbnail <> ''
  ORDER BY posts.pub_date DESC
  LIMIT 1`).Scan(&post.Title, &post.FileName, &post.Content, &post.Description, &post.PubDate, &post.UpdateDate, &thumbnailJSON)
  if err != nil {
    if err == sql.ErrNoRows {
      return Post{}, errors.New("no valid thumbnails exist")
    }
    return Post{}, err
  }

  err = json.Unmarshal([]byte(thumbnailJSON.String), &post.Thumbnail)
  if err != nil {
    return Post{}, err
  }
  return post, nil
}

func FetchTag(tagName string) (tag Tag, err error) {
  db := OpenDB()
  defer CloseDB(db)

  err = db.QueryRow(`SELECT name, description, category
  FROM tags
  WHERE name = ?`, tagName).Scan(&tag.Name, &tag.Description, &tag.Category)
  if err != nil {
    return Tag{}, err
  }

  return tag, nil
}

func DoesPostExist(fileName string) bool {
  db := OpenDB()
  defer CloseDB(db)

  var count int
  err := db.QueryRow(`SELECT COUNT(*)
  FROM posts
  WHERE file_name = ?`, fileName).Scan(&count)
  if err != nil {
    log.Println(err.Error())
    return false
  }

  return count > 0
}

func DoesTagExist(tag string) bool {
  db := OpenDB()
  defer CloseDB(db)

  var count int
  err := db.QueryRow(`SELECT COUNT(*)
  FROM tags
  WHERE name = ?`, tag).Scan(&count)
  if err != nil {
    log.Println(err.Error())
    return false
  }

  return count > 0
}

func checkPost(p Post) error { // don't need to check thumbnail
  if p.Title == "" ||
    p.FileName == "" ||
    p.Content == "" ||
    p.Description == "" ||
    p.PubDate == "" ||
    p.UpdateDate == "" ||
    len(p.Tags) < 1 {
      return ErrNotComplete
  }
  return nil
}

func checkTag(t Tag) error {
  if t.Name == "" ||
    t.Category == "" ||
    t.Description == "" {
      return ErrNotComplete
  }
  return nil
}

func AddPost(post Post) (err error) {
  if err = checkPost(post); err != nil {
    return err
  }

  // ensure tag existence
  for _, t := range post.Tags {
    if !DoesTagExist(t.Name) {
      return errors.New("missing tag: "+ t.Name)
    }
  }

  db := OpenDB()
  defer CloseDB(db)

  var jsonThumbnail []byte
  if post.Thumbnail.Src != "" {
    jsonThumbnail, err = json.Marshal(post.Thumbnail)
    if err != nil {
      return err
    }
  }

  _, err = db.Exec(`INSERT INTO posts (title, file_name, content, description, pub_date, update_date, thumbnail)
  VALUES
  (?,  ?,  ?,  ?,  ?,  ?, ?)
`, post.Title, post.FileName, post.Content, post.Description, post.PubDate, post.UpdateDate, string(jsonThumbnail))
  if err != nil {
    return err
  }

  for _, t := range post.Tags {
    _, err = db.Exec(`INSERT INTO posts_tags (post_id, tag_id)
    VALUES
    ((SELECT id
    FROM posts
    WHERE file_name = ?),
    (SELECT id
    FROM tags
    WHERE name = ?))`, post.FileName, t.Name)
  }

  return nil
}

func AddTag(tag Tag) bool {
  // connect to db and add tag

  return false
}
