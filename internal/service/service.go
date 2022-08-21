package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/spf13/cast"

	"github.com/carfloresf/the-sun-god/internal/constants"
	"github.com/carfloresf/the-sun-god/internal/model"
	"github.com/carfloresf/the-sun-god/pkg/subreddit"
)

var ErrMutuallyExclusiveFields = fmt.Errorf("link and content are mutually exclusive fields")
var ErrLinkOrContentNotFound = fmt.Errorf("link or content are needed")

type GetFeedResponse struct {
	Posts      []model.Post `json:"data"`
	Pagination Pagination   `json:"pagination"`
	Error      error        `json:"error,omitempty"`
}

type Pagination struct {
	PageCount int    `json:"pageCount"`
	Previous  string `json:"prev,omitempty"`
	Next      string `json:"next,omitempty"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Total     int    `json:"total"`
}

type Storage interface {
	CreatePost(post model.Post) (uuid.UUID, error)
	UpdatePostScore(pID uuid.UUID, score int) error
	GetPosts(limit, offset int) ([]model.Post, int, error)
	GetPromotedPosts(limit int) ([]model.Post, error)
}

type Service struct {
	Storage Storage
}

func (s *Service) UpdatePostScore(pID uuid.UUID, score int) error {
	return s.Storage.UpdatePostScore(pID, score)
}

func (s *Service) CreatePost(post model.Post) (uuid.UUID, error) {
	if post.Link != "" && post.Content != "" {
		return uuid.UUID{}, ErrMutuallyExclusiveFields
	}

	if !subreddit.Exists(post.Subreddit) {
		return uuid.UUID{}, subreddit.ErrSubredditNotFound
	}

	post.GenerateAuthor()
	post.GenerateID()

	return s.Storage.CreatePost(post)
}

func (s *Service) GetFeed(page, limit int) (GetFeedResponse, error) {
	posts, total, err := s.Storage.GetPosts(limit, (page-1)*limit)
	if err != nil {
		return GetFeedResponse{}, err
	}

	promotedPosts, err := s.Storage.GetPromotedPosts(constants.DefaultPromotedPosts)
	if err != nil {
		return GetFeedResponse{}, err
	}

	resultPosts := insertAllPromotedPosts(posts, promotedPosts)

	response := GetFeedResponse{
		Posts: resultPosts,
		Pagination: Pagination{
			PageCount: len(resultPosts),
			Page:      page,
			Limit:     limit,
			Total:     total,
		},
	}

	if page > 1 {
		response.Pagination.Previous = "/v1/feed?page=" + cast.ToString(page-1)
	}

	if page*limit < total {
		response.Pagination.Next = "/v1/feed?page=" + cast.ToString(page+1)
	}

	return response, nil
}

func adjacentNSFW(posts []model.Post, index int) bool {
	if index == 0 {
		return posts[index].NSFW
	}

	if index == len(posts)-1 {
		return posts[index-1].NSFW
	}

	return posts[index-1].NSFW || posts[index].NSFW
}

func insertAllPromotedPosts(posts []model.Post, promotedPosts []model.Post) []model.Post {
	multiplier := 1 // mostly used for bigger limits
	promotedIndex := 0
	destinationPosts := make([]model.Post, 0)

	for i := 0; i < len(posts); i++ {
		if i == 1 && len(posts) >= 3 && !adjacentNSFW(posts, i) {
			if len(promotedPosts)-1 >= promotedIndex {
				destinationPosts = append(destinationPosts, promotedPosts[promotedIndex])
				promotedIndex++
			}
		}

		if i == multiplier*14 && len(posts) >= multiplier*15 && !adjacentNSFW(posts, i) {
			if len(promotedPosts)-1 >= promotedIndex {
				destinationPosts = append(destinationPosts, promotedPosts[promotedIndex])
				multiplier++
				promotedIndex++
			}
		}

		destinationPosts = append(destinationPosts, posts[i])
	}

	return destinationPosts
}
