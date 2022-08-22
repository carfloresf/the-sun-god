package storage

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"

	"github.com/carfloresf/the-sun-god/config"
	"github.com/carfloresf/the-sun-god/internal/model"
)

const (
	insertPostQuery = `INSERT
						INTO post(id,title,author,subreddit,link,content,score,promoted,nsfw) 
						VALUES    ($1, $2, $3, $4, $5, $6, $7, $8, $9)`
	updatePostScoreQuery = `UPDATE post SET score = $1 WHERE id = $2`
	countTotalPostsQuery = `SELECT COUNT(*) FROM post WHERE promoted = false`
	queryPosts           = `SELECT id,title,author,subreddit,link,content,score,promoted,nsfw,created_at 
							FROM post WHERE promoted = $1 ORDER BY score DESC, created_at DESC LIMIT $2 OFFSET $3;`
	queryPromotedPosts = `SELECT id,title,author,subreddit,link,content,score,promoted,nsfw,created_at 
							FROM post WHERE promoted = true ORDER BY RANDOM() LIMIT $1;`
)

var allStatements = []string{insertPostQuery, updatePostScoreQuery, countTotalPostsQuery, queryPosts, queryPromotedPosts}

type DB struct {
	Conn               *sql.DB
	preparedStatements map[string]*sql.Stmt
}

func NewStorage(config *config.DB) (*DB, error) {
	sqliteDatabase, err := sql.Open("sqlite3", config.DBFile)
	if err != nil {
		log.Fatal("error opening Conn connection: %w", err)
		return nil, err
	}

	db := DB{
		Conn: sqliteDatabase,
	}

	return &db, nil
}

func (d *DB) PrepareAllStatements() error {
	d.preparedStatements = make(map[string]*sql.Stmt)

	err := d.prepareStatements(allStatements)
	if err != nil {
		log.Errorln("error preparing statements:", err)
		return err
	}

	return nil
}

func (d *DB) prepareStatements(queries []string) error {
	for _, query := range queries {
		stmt, err := d.Conn.Prepare(query)
		if err != nil {
			return fmt.Errorf("error while preparing statement %w", err)
		}

		d.preparedStatements[query] = stmt
	}

	return nil
}

func (d *DB) Close() error {
	for _, stmt := range d.preparedStatements {
		err := stmt.Close()
		if err != nil {
			log.Errorf("error closing statement: %s", err)
			return err
		}
	}

	return d.Conn.Close()
}

func (d *DB) UpdatePostScore(pID uuid.UUID, score int) error {
	_, err := d.preparedStatements[updatePostScoreQuery].Exec(score, pID)
	if err != nil {
		log.Errorf("error updating post score: %v", err)
		return err
	}

	return nil
}

func (d *DB) CreatePost(post model.Post) (uuid.UUID, error) {
	_, err := d.preparedStatements[insertPostQuery].Exec(
		post.ID,
		post.Title,
		post.Author,
		post.Subreddit,
		post.Link,
		post.Content,
		post.Score,
		post.Promoted,
		post.NSFW)
	if err != nil {
		log.Errorf("error executing statement: %s", err)
		return uuid.UUID{}, err
	}

	return post.ID, nil
}

func (d *DB) GetPosts(limit, offset int) ([]model.Post, int, error) {
	var posts []model.Post

	var total int

	err := d.preparedStatements[countTotalPostsQuery].QueryRow().Scan(&total)
	if err != nil {
		log.Errorf("error getting total posts: %s", err)
		return nil, 0, err
	}

	rows, err := d.preparedStatements[queryPosts].Query(false, limit, offset)
	if err != nil {
		log.Errorf("error getting posts: %s", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var post model.Post

		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Author,
			&post.Subreddit,
			&post.Link,
			&post.Content,
			&post.Score,
			&post.Promoted,
			&post.NSFW,
			&post.CreatedAt)
		if err != nil {
			log.Errorf("error scanning post: %s", err)
			return nil, total, err
		}

		posts = append(posts, post)
	}

	return posts, total, nil
}

func (d *DB) GetPromotedPosts(limit int) ([]model.Post, error) {
	var posts []model.Post

	rows, err := d.preparedStatements[queryPromotedPosts].Query(limit)
	if err != nil {
		log.Errorf("error getting promoted posts: %s", err)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var post model.Post

		err := rows.Scan(&post.ID,
			&post.Title,
			&post.Author,
			&post.Subreddit,
			&post.Link,
			&post.Content,
			&post.Score,
			&post.Promoted,
			&post.NSFW,
			&post.CreatedAt,
		)
		if err != nil {
			log.Errorf("error scanning post: %s", err)
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}
