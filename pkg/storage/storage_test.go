package storage

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/carfloresf/the-sun-god/internal/model"
)

func TestDB_prepareAllStatements(t *testing.T) {
	tests := []struct {
		name        string
		mockClosure func(mock sqlmock.Sqlmock)
		wantErr     bool
		err         error
	}{
		{name: "success",
			mockClosure: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(insertPostQuery)))
				mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(updatePostScoreQuery)))
				mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(countTotalPostsQuery)))
				mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(queryPosts)))
				mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(queryPromotedPosts)))
			},
			wantErr: false,
			err:     nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			d := &DB{
				Conn: db,
			}

			tt.mockClosure(mock)

			err := d.PrepareAllStatements()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDB_UpdatePostScore(t *testing.T) {
	type args struct {
		pID   uuid.UUID
		score int
	}

	testUUID := uuid.New()
	testScore := 1977

	tests := []struct {
		name        string
		mockClosure func(mock sqlmock.Sqlmock)
		args        args
		wantErr     bool
		err         error
	}{
		{"success",
			func(mock sqlmock.Sqlmock) {
				expectedPrepare := mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(updatePostScoreQuery)))
				expectedPrepare.ExpectExec().WithArgs(testScore, testUUID).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			args{testUUID, testScore},
			false,
			nil,
		},
		{"error",
			func(mock sqlmock.Sqlmock) {
				expectedPrepare := mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(updatePostScoreQuery)))
				expectedPrepare.ExpectExec().WithArgs(testScore, testUUID).WillReturnError(fmt.Errorf("error updating post score"))
			},
			args{testUUID, testScore},
			true,
			fmt.Errorf("error updating post score"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			preparedStatements := make(map[string]*sql.Stmt)
			d := &DB{
				Conn:               db,
				preparedStatements: preparedStatements,
			}

			tt.mockClosure(mock)

			err := d.prepareStatements([]string{updatePostScoreQuery})
			assert.NoError(t, err)
			assert.Equal(t, 1, len(d.preparedStatements))

			err = d.UpdatePostScore(tt.args.pID, tt.args.score)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDB_CreatePost(t *testing.T) {
	type args struct {
		post model.Post
	}

	successPost := model.Post{
		ID:        uuid.New(),
		Score:     1977,
		Author:    "t2_12345",
		Subreddit: "r/test",
		Link:      "https://www.google.com",
		Promoted:  false,
		NSFW:      false,
	}

	tests := []struct {
		name        string
		mockClosure func(mock sqlmock.Sqlmock)
		args        args
		want        uuid.UUID
		wantErr     bool
		err         error
	}{
		{
			name: "success",
			mockClosure: func(mock sqlmock.Sqlmock) {
				expectedPrepare := mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(insertPostQuery)))
				expectedPrepare.ExpectExec().
					WithArgs(
						successPost.ID,
						successPost.Title,
						successPost.Author,
						successPost.Subreddit,
						successPost.Link,
						successPost.Content,
						successPost.Score,
						successPost.Promoted,
						successPost.NSFW).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			args: args{
				post: successPost,
			},
			want:    successPost.ID,
			wantErr: false,
			err:     nil,
		}, {
			name: "error",
			mockClosure: func(mock sqlmock.Sqlmock) {
				expectedPrepare := mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(insertPostQuery)))
				expectedPrepare.ExpectExec().WithArgs(
					successPost.ID,
					successPost.Title,
					successPost.Author,
					successPost.Subreddit,
					successPost.Link,
					successPost.Content,
					successPost.Score,
					successPost.Promoted,
					successPost.NSFW).WillReturnError(fmt.Errorf("error inserting post"))
			},
			args: args{
				post: successPost,
			},
			want:    uuid.Nil,
			wantErr: true,
			err:     fmt.Errorf("error inserting post"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			preparedStatements := make(map[string]*sql.Stmt)
			d := &DB{
				Conn:               db,
				preparedStatements: preparedStatements,
			}

			tt.mockClosure(mock)

			err := d.prepareStatements([]string{insertPostQuery})
			assert.NoError(t, err)
			assert.Equal(t, 1, len(d.preparedStatements))

			pID, err := d.CreatePost(tt.args.post)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, pID)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDB_GetPosts(t *testing.T) {
	type args struct {
		Limit  int
		Offset int
	}

	ps := []model.Post{
		{
			ID:        uuid.New(),
			Score:     1977,
			Title:     "test title",
			Author:    "t2_12345",
			Subreddit: "/r/test",
			Link:      "https://www.google.com",
			Promoted:  false,
			NSFW:      false,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Score:     1,
			Title:     "test title 2",
			Author:    "t2_123456",
			Subreddit: "/r/test",
			Content:   "test content",
			Promoted:  false,
			NSFW:      false,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Score:     2,
			Title:     "test title 3",
			Author:    "t2_1234567",
			Subreddit: "/r/test",
			Content:   "test content",
			Promoted:  false,
			NSFW:      false,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Score:     3,
			Title:     "test title 4",
			Author:    "t2_12345678",
			Subreddit: "/r/test",
			Content:   "test content XXX",
			Promoted:  false,
			NSFW:      true,
			CreatedAt: time.Now(),
		},
	}

	recommendedArgs := args{Limit: 4, Offset: 0}

	tests := []struct {
		name        string
		mockClosure func(mock sqlmock.Sqlmock)
		args        args
		want        []model.Post
		wantTotal   int
		wantErr     bool
		err         error
	}{
		{name: "success",
			mockClosure: func(mock sqlmock.Sqlmock) {
				expectedPrepareCount := mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(countTotalPostsQuery)))
				expectedPrepare := mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(queryPosts)))
				expectedPrepareCount.ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(4))
				expectedPrepare.ExpectQuery().
					WithArgs(false, recommendedArgs.Limit, recommendedArgs.Offset).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "subreddit", "link",
						"content", "score", "promoted", "nsfw", "created_at"}).
						AddRow(ps[0].ID, ps[0].Title, ps[0].Author, ps[0].Subreddit, ps[0].Link,
							ps[0].Content, ps[0].Score, ps[0].Promoted, ps[0].NSFW, ps[0].CreatedAt).
						AddRow(ps[1].ID, ps[1].Title, ps[1].Author, ps[1].Subreddit, ps[1].Link,
							ps[1].Content, ps[1].Score, ps[1].Promoted, ps[1].NSFW, ps[1].CreatedAt).
						AddRow(ps[2].ID, ps[2].Title, ps[2].Author, ps[2].Subreddit, ps[2].Link,
							ps[2].Content, ps[2].Score, ps[2].Promoted, ps[2].NSFW, ps[2].CreatedAt).
						AddRow(ps[3].ID, ps[3].Title, ps[3].Author, ps[3].Subreddit, ps[3].Link,
							ps[3].Content, ps[3].Score, ps[3].Promoted, ps[3].NSFW, ps[3].CreatedAt))
			},
			args:      recommendedArgs,
			want:      ps,
			wantTotal: 4,
			wantErr:   false,
			err:       nil,
		},
		{name: "error",
			mockClosure: func(mock sqlmock.Sqlmock) {
				expectedPrepareCount := mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(countTotalPostsQuery)))
				expectedPrepare := mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(queryPosts)))
				expectedPrepareCount.ExpectQuery().WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(4))
				expectedPrepare.ExpectQuery().
					WithArgs(false, recommendedArgs.Limit, recommendedArgs.Offset).
					WillReturnError(fmt.Errorf("error getting posts"))
			},
			args:      recommendedArgs,
			want:      nil,
			wantTotal: 0,
			wantErr:   true,
			err:       fmt.Errorf("error getting posts"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			preparedStatements := make(map[string]*sql.Stmt)
			d := &DB{
				Conn:               db,
				preparedStatements: preparedStatements,
			}

			tt.mockClosure(mock)

			err := d.prepareStatements([]string{countTotalPostsQuery, queryPosts})
			assert.NoError(t, err)
			assert.Equal(t, 2, len(d.preparedStatements))

			posts, total, err := d.GetPosts(tt.args.Limit, tt.args.Offset)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, posts)
			assert.Equal(t, tt.wantTotal, total)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDB_GetPromotedPosts(t *testing.T) {
	type args struct {
		Limit int
	}

	ps := []model.Post{
		{
			ID:        uuid.New(),
			Score:     1977,
			Title:     "test title",
			Author:    "t2_12345",
			Subreddit: "/r/test",
			Link:      "https://www.google.com",
			Promoted:  true,
			NSFW:      false,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Score:     1,
			Title:     "test title 2",
			Author:    "t2_123456",
			Subreddit: "/r/test",
			Content:   "test content",
			Promoted:  true,
			NSFW:      false,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Score:     2,
			Title:     "test title 3",
			Author:    "t2_1234567",
			Subreddit: "/r/test",
			Content:   "test content",
			Promoted:  true,
			NSFW:      false,
			CreatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			Score:     3,
			Title:     "test title 4",
			Author:    "t2_12345678",
			Subreddit: "/r/test",
			Content:   "test content XYZ",
			Promoted:  true,
			NSFW:      false,
			CreatedAt: time.Now(),
		},
	}

	recommendedArgs := args{Limit: 4}

	tests := []struct {
		name        string
		mockClosure func(mock sqlmock.Sqlmock)
		args        args
		want        []model.Post
		wantErr     bool
		err         error
	}{
		{name: "success",
			mockClosure: func(mock sqlmock.Sqlmock) {
				expectedPrepare := mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(queryPromotedPosts)))
				expectedPrepare.ExpectQuery().
					WithArgs(recommendedArgs.Limit).
					WillReturnRows(sqlmock.NewRows([]string{"id", "title", "author", "subreddit", "link",
						"content", "score", "promoted", "nsfw", "created_at"}).
						AddRow(ps[0].ID, ps[0].Title, ps[0].Author, ps[0].Subreddit, ps[0].Link,
							ps[0].Content, ps[0].Score, ps[0].Promoted, ps[0].NSFW, ps[0].CreatedAt).
						AddRow(ps[1].ID, ps[1].Title, ps[1].Author, ps[1].Subreddit, ps[1].Link,
							ps[1].Content, ps[1].Score, ps[1].Promoted, ps[1].NSFW, ps[1].CreatedAt).
						AddRow(ps[2].ID, ps[2].Title, ps[2].Author, ps[2].Subreddit, ps[2].Link,
							ps[2].Content, ps[2].Score, ps[2].Promoted, ps[2].NSFW, ps[2].CreatedAt).
						AddRow(ps[3].ID, ps[3].Title, ps[3].Author, ps[3].Subreddit, ps[3].Link,
							ps[3].Content, ps[3].Score, ps[3].Promoted, ps[3].NSFW, ps[3].CreatedAt))
			},
			args:    recommendedArgs,
			want:    ps,
			wantErr: false,
			err:     nil,
		},
		{name: "error",
			mockClosure: func(mock sqlmock.Sqlmock) {
				expectedPrepare := mock.ExpectPrepare(fmt.Sprintf("^%s$", regexp.QuoteMeta(queryPromotedPosts)))
				expectedPrepare.ExpectQuery().
					WithArgs(recommendedArgs.Limit).
					WillReturnError(fmt.Errorf("error getting promoted posts"))
			},
			args:    recommendedArgs,
			want:    nil,
			wantErr: true,
			err:     fmt.Errorf("error getting promoted posts"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, _ := sqlmock.New()
			defer db.Close()

			preparedStatements := make(map[string]*sql.Stmt)
			d := &DB{
				Conn:               db,
				preparedStatements: preparedStatements,
			}

			tt.mockClosure(mock)

			err := d.prepareStatements([]string{queryPromotedPosts})
			assert.NoError(t, err)
			assert.Equal(t, 1, len(d.preparedStatements))

			posts, err := d.GetPromotedPosts(tt.args.Limit)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			for _, post := range posts {
				assert.Equal(t, true, post.Promoted)
			}

			assert.Equal(t, tt.want, posts)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
