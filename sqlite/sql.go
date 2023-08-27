// Package sqlite for creating, reading, and writing to a sqlite database for my blog
package sqlite

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3" // blank import so the sql engine can do its thing in init
)

var (
	pathToDB       = "./posts.sqlite3"
	errNotComplete = errors.New("unfilled required attributes")
)

// DB custom type is gonna be my attempt at making paths to the database
// defined in the original package. Still a todo for now.
type DB struct {
	PathToDB string
}

// Post type has pretty much all of the content and metadata that I'd need
// to manipulate for a complex blogging system that I can still manage.
type Post struct {
	Title       string
	FileName    string
	Content     string
	Description string
	PubDate     string
	UpdateDate  string
	Tags        []Tag
	Thumbnail   Img
}

// Img struct exists for my thumbnail image in the front of my site, but
// I will definitely use this package for other things that I'll need for
// other friend's sites I'm designing
type Img struct {
	Src   string `json:"src"`
	Alt   string `json:"alt"`
	Title string `json:"title"`
}

// Tag struct is there to aggregate together posts with common themes.
// I might add an html related attribute like content in post to let
// tag pages be more than just a paragraph.
type Tag struct {
	Name        string
	Category    string
	Description string
}

// OpenDB just opens the connection to the db from the pathToDB package var
func OpenDB() (db *sql.DB) {
	db, err := sql.Open("sqlite3", pathToDB)
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}

// CloseDB is the mirror of open. Generally just defer it immediately after
// OpenDB
func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

// MakeDB makes the basic schema for the database if it doesn't exist.
// originally, it had a trigger for update_date whenever something changed, but
// that was too clunky to deal with for me. I might eventually add it searching
// for a backup .sql file in the same directory as the db in case there's
// backup data.
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

// AggregatePosts gets all the posts and tags metadata from the db into a slice
// of posts sorted in reverse chron order. if you give it a tag, it'll only
// return posts with that tag. giving the func a negative number returns an
// empty slice and giving it zero will return all the entries that match. If
// you give it an integer less than the total amount of posts available with
// the filter, it'll return that many posts back in the slice, still in reverse
// chron order.
func AggregatePosts(postQty int, filterTag string) (posts []Post, err error) {
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

// FetchPost brings back the struct data of a single post including a tag slice
// of all matching tags to post in associative identity.
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
		}
		post.Thumbnail = thumbnail
	}

	return post, nil
}

// FetchThumbnail is for a very niche thing I needed for my home page that
// displays my latest photos post with a valid thumbnail. Sqlite doesn't have
// structs, and I don't wanna bother with pgsql or mariadb for something this
// small, so I just marshal and unmarshal json into the db as text. This func
// fetches just that post and gives back not only the img struct but the rest
// of the post as well for easily being able to link back to the post.
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

// FetchTag is basically identical to FetchPost but way smaller and less
// complicated
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

// DoesPostExist is an evolution from my original filesystem func that I still
// use called doesFileExist that ensures a certain post in actually in the db.
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

// DoesTagExist is the exact same deal as DoesPostExist
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

// checkPost is an internal func to ensure that all the attributes in a post
// that are NOT NULL in the db are, in fact, not null before loading them into
// the db.
func checkPost(p Post) error { // don't need to check thumbnail
	if p.Title == "" ||
		p.FileName == "" ||
		p.Content == "" ||
		p.Description == "" ||
		p.PubDate == "" ||
		p.UpdateDate == "" ||
		len(p.Tags) < 1 {
		return errNotComplete
	}
	return nil
}

// checkTag is idential but smaller
func checkTag(t Tag) error {
	if t.Name == "" ||
		t.Category == "" ||
		t.Description == "" {
		return errNotComplete
	}
	return nil
}

// AddPost checks to see if a post is valid first will checkPost then also by
// ensuring every tag that is inside the post exists in the db. You don't even
// need to add anything more than the names of the tags in the tag slice of the
// post as that's how it's checked before the post is inserted into the db as
// well as being how the tag_id is filled into the posts_tags associative
// identity for linking the metadata together.
func AddPost(post Post) (err error) {
	if err = checkPost(post); err != nil {
		return err
	}

	// ensure tag existence
	for _, t := range post.Tags {
		if !DoesTagExist(t.Name) {
			return errors.New("missing tag: " + t.Name)
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

// AddTag is just a simpler version of AddPost. The only filtering done before
// is seeing if your tag struct has all the required attributes to add to the
// db.
func AddTag(tag Tag) (err error) {
	if err = checkTag(tag); err != nil {
		return err
	}

	db := OpenDB()
	defer CloseDB(db)

	_, err = db.Exec(`INSERT INTO tags (name, category, description)
  VALUES
  (?,  ?,  ?)
`, tag.Name, tag.Category, tag.Description)
	if err != nil {
		return err
	}

	return nil
}

// DeletePost is incomplete
func DeletePost(fileName string) (err error) {
	if !DoesPostExist(fileName) {
		return errors.New(fileName + " doesn't exist")
	}

	db := OpenDB()
	defer CloseDB(db)

	_, err = db.Exec(`DELETE FROM posts
  WHERE file_name = ?`, fileName)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTag is incomplete
func DeleteTag(tagName string) (err error) {
	if !DoesTagExist(tagName) {
		return errors.New(tagName + " doesn't exist")
	}

	db := OpenDB()
	defer CloseDB(db)

	_, err = db.Exec(`DELETE FROM tags
  WHERE name = ?`, tagName)
	if err != nil {
		return err
	}

	return nil
}
