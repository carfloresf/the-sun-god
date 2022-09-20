package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPost_GenerateAuthor(t *testing.T) {
	type fields struct {
		ID        uuid.UUID
		Title     string
		Author    string
		Subreddit string
		Link      string
		Content   string
		Score     int
		Promoted  bool
		NSFW      bool
		CreatedAt time.Time
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{"test",
			fields{
				Title:     "test",
				Subreddit: "/r/test",
				Link:      "https://www.google.com",
				Content:   "test",
				Score:     5,
				Promoted:  false,
				NSFW:      false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Post{
				Title:     tt.fields.Title,
				Subreddit: tt.fields.Subreddit,
				Link:      tt.fields.Link,
				Content:   tt.fields.Content,
				Score:     tt.fields.Score,
				Promoted:  tt.fields.Promoted,
				NSFW:      tt.fields.NSFW,
			}

			assert.Empty(t, p.Author)

			p.GenerateAuthor()

			assert.NotEmpty(t, p.Author)
			assert.Equal(t, "t2_", p.Author[:3])
		})
	}
}

func TestPost_GenerateID(t *testing.T) {
	type fields struct {
		ID        uuid.UUID
		Title     string
		Author    string
		Subreddit string
		Link      string
		Content   string
		Score     int
		Promoted  bool
		NSFW      bool
		CreatedAt time.Time
	}

	tests := []struct {
		name   string
		fields fields
	}{
		{"newUUID",
			fields{
				Title:     "test",
				Subreddit: "/r/test",
				Link:      "https://www.google.com",
				Content:   "test",
				Score:     5,
				Promoted:  false,
				NSFW:      false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Post{
				Title:     tt.fields.Title,
				Subreddit: tt.fields.Subreddit,
				Link:      tt.fields.Link,
				Content:   tt.fields.Content,
				Score:     tt.fields.Score,
				Promoted:  tt.fields.Promoted,
				NSFW:      tt.fields.NSFW,
			}

			assert.Empty(t, p.ID)
			p.GenerateID()

			assert.NotEmpty(t, p.ID)
			assert.True(t, IsValidUUID(p.ID.String()))
		})
	}
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
