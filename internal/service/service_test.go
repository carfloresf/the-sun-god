package service

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"

	"github.com/carfloresf/the-sun-god/internal/model"
	"github.com/carfloresf/the-sun-god/pkg/subreddit"
)

var mockedStorage *StorageMock
var postUUID uuid.UUID

// testdata
var p10 []model.Post
var p10NSFW []model.Post
var p25 []model.Post
var p100 []model.Post
var feed11 []model.Post
var feed27 []model.Post
var feed502 []model.Post
var promotedPost2 []model.Post

// insert func to insert promoted posts
func insert(a []model.Post, index int, value model.Post) []model.Post {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}

	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value

	return a
}

func TestMain(m *testing.M) {
	postUUID = uuid.New()

	// Generated feeds with test data
	f := faker.New()

	// Generated 10 posts results
	for i := 0; i < 10; i++ {
		post := model.Post{}
		f.Struct().Fill(&post)
		post.GenerateID()
		post.Score = i
		post.Promoted = false
		post.NSFW = false
		pNSFW := post
		pNSFW.NSFW = true
		p10NSFW = append(p10NSFW, pNSFW)
		p10 = append(p10, post)
	}

	// Generated 25 posts results
	for i := 0; i < 25; i++ {
		post := model.Post{}
		f.Struct().Fill(&post)
		post.GenerateID()
		post.Score = i
		post.Promoted = false
		post.NSFW = false
		p25 = append(p25, post)
	}

	// Generated 100 posts results
	for i := 0; i < 500; i++ {
		post := model.Post{}
		f.Struct().Fill(&post)
		post.GenerateID()
		post.Score = i
		post.Promoted = false
		post.NSFW = false
		p100 = append(p100, post)
	}

	// Promoted posts results
	promotedPost2 = []model.Post{
		{
			ID:       postUUID,
			Title:    "promoted test",
			Content:  "test promoted content 2",
			Author:   "t2_1",
			Promoted: true,
		},
		{
			ID:       postUUID,
			Title:    "promoted test 2",
			Content:  "test promoted content 2",
			Author:   "t2_1123123",
			Promoted: true,
		},
	}

	feed11 = append(feed11, p10...)
	feed11 = insert(feed11, 1, promotedPost2[0])

	feed27 = append(feed27, p25...)
	feed27 = insert(feed27, 1, promotedPost2[0])
	feed27 = insert(feed27, 15, promotedPost2[1])

	feed502 = append(feed502, p100...)
	feed502 = insert(feed502, 1, promotedPost2[0])
	feed502 = insert(feed502, 15, promotedPost2[1])

	mockedStorage = &StorageMock{
		CreatePostFunc: func(post model.Post) (uuid.UUID, error) {
			if post.Title == "test" {
				return postUUID, nil
			} else {
				return uuid.UUID{}, fmt.Errorf("unable to create post")
			}
		},
		GetPostsFunc: func(limit int, offset int) ([]model.Post, int, error) {
			if limit == 10 && offset == 0 {
				return p10, 10, nil
			}

			if limit == 25 && offset == 0 {
				return p25, 25, nil
			}

			if limit == 24 && offset == 0 {
				return p10NSFW, 10, nil
			}

			return nil, 0, fmt.Errorf("unable to get posts")
		},
		GetPromotedPostsFunc: func(limit int) ([]model.Post, error) {
			return promotedPost2, nil
		},
		UpdatePostScoreFunc: func(pID uuid.UUID, score int) error {
			if pID != postUUID {
				return fmt.Errorf("UpdatePostScore() pID = %v, want %v", pID, postUUID)
			}
			if pID == postUUID && score != 100 {
				return nil
			}

			return fmt.Errorf("unable to update score for post %v", pID)
		},
	}

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestService_UpdatePostScore(t *testing.T) {
	type args struct {
		pID   uuid.UUID
		score int
	}

	tests := []struct {
		name    string
		storage *StorageMock
		args    args
		wantErr bool
		err     error
	}{
		{"success",
			mockedStorage,
			args{pID: postUUID, score: 5},
			false,
			nil},
		{"error",
			mockedStorage,
			args{pID: postUUID, score: 100},
			true,
			fmt.Errorf("unable to update score for post %v", postUUID),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Storage: tt.storage,
			}

			err := s.UpdatePostScore(tt.args.pID, tt.args.score)
			if tt.wantErr {
				assert.Error(t, err, tt.err)
			} else if err != nil {
				t.Errorf("UpdatePostScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_CreatePost(t *testing.T) {
	type args struct {
		post model.Post
	}

	tests := []struct {
		name              string
		storage           *StorageMock
		args              args
		registerSubReddit func()
		want              uuid.UUID
		wantErr           bool
		err               error
	}{
		{"success",
			mockedStorage,
			args{
				post: model.Post{
					Title:     "test",
					Author:    "t2_1",
					Subreddit: "/r/test",
					Score:     5,
					Content:   "test content",
				},
			},
			func() {
				subreddit.Set("/r/test")
			},
			postUUID,
			false,
			nil,
		},
		{"post&link-provided",
			mockedStorage,
			args{
				post: model.Post{
					Title:     "test",
					Author:    "t2_1123",
					Subreddit: "/r/test",
					Score:     5,
					Content:   "test content",
					Link:      "http//test.com",
				},
			}, func() {
			subreddit.Set("/r/test")
		},
			uuid.UUID{},
			true,
			ErrMutuallyExclusiveFields,
		},
		{"subreddit-not-found",
			mockedStorage,
			args{
				post: model.Post{
					Title:     "test",
					Author:    "t2_1123",
					Subreddit: "/r/testxxx",
					Score:     5,
					Content:   "test content",
				},
			}, func() {},
			uuid.UUID{},
			true,
			subreddit.ErrSubredditNotFound,
		},
		{"storage-error",
			mockedStorage,
			args{
				post: model.Post{
					Title:     "test-error",
					Author:    "t2_1123",
					Subreddit: "/r/testxxx",
					Score:     5,
					Content:   "test content",
				},
			}, func() {
			subreddit.Set("/r/testxxx")
		},
			uuid.UUID{},
			true,
			fmt.Errorf("unable to create post"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Storage: tt.storage,
			}

			tt.registerSubReddit()

			pID, err := s.CreatePost(tt.args.post)
			if tt.wantErr {
				assert.EqualError(t, err, tt.err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pID)
				assert.Equal(t, tt.want, pID)
			}
		})
	}
}

func TestService_GetFeed(t *testing.T) {
	type args struct {
		page  int
		limit int
	}

	tests := []struct {
		name    string
		storage Storage
		args    args
		want    GetFeedResponse
		wantErr bool
	}{
		{"success-11",
			mockedStorage,
			args{page: 1, limit: 10},
			GetFeedResponse{
				Posts: feed11,
				Pagination: Pagination{
					Page:      1,
					Limit:     10,
					Total:     10,
					PageCount: 11,
				},
			},
			false,
		},
		{"success-27",
			mockedStorage,
			args{page: 1, limit: 25},
			GetFeedResponse{
				Posts: feed27,
				Pagination: Pagination{
					Page:      1,
					Limit:     25,
					Total:     25,
					PageCount: 27,
				},
			},
			false,
		},
		{"success-NSFW",
			mockedStorage,
			args{page: 1, limit: 24},
			GetFeedResponse{
				Posts: p10NSFW,
				Pagination: Pagination{
					Page:      1,
					Limit:     24,
					Total:     10,
					PageCount: 10,
				},
			},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				Storage: tt.storage,
			}
			got, err := s.GetFeed(tt.args.page, tt.args.limit)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_adjacentNSFW(t *testing.T) {
	type args struct {
		posts []model.Post
		index int
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"adjacent-nsfw", args{p10NSFW, 0}, true},
		{"no-adjacent", args{p25, 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, adjacentNSFW(tt.args.posts, tt.args.index), "adjacentNSFW(%v, %v)", tt.args.posts, tt.args.index)
		})
	}
}

func Test_insertAllPromotedPosts(t *testing.T) {
	type args struct {
		posts         []model.Post
		promotedPosts []model.Post
	}

	tests := []struct {
		name string
		args args
		want []model.Post
	}{
		{"insert-all-promoted-posts", args{p10, promotedPost2}, feed11},
		{"500-posts", args{p100, promotedPost2}, feed502},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, insertAllPromotedPosts(tt.args.posts, tt.args.promotedPosts), "insertAllPromotedPosts(%v, %v)", tt.args.posts, tt.args.promotedPosts)
		})
	}
}
