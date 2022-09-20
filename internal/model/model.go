package model

import (
	"time"

	"github.com/google/uuid"

	"github.com/carfloresf/the-sun-god/pkg/random"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title" binding:"required" fake:"??????????????????"`
	Author    string    `json:"author" fake:"t2_??????"`
	Subreddit string    `json:"subreddit" fake:"/r/???????"`
	Link      string    `json:"link,omitempty"  binding:"omitempty,required_without=Content,url" fake:"http://???????"`
	Content   string    `json:"content,omitempty"  binding:"required_without=Link" fake:"???????????? ???????????? ????????????"`
	Score     int       `json:"score" fake:"####"`
	Promoted  bool      `json:"promoted"`
	NSFW      bool      `json:"nsfw" fake:"@bool"`
	CreatedAt time.Time `json:"-"`
}

func (p *Post) GenerateAuthor() {
	suffix := random.RandSequence(8)

	p.Author = "t2_" + suffix
}

func (p *Post) GenerateID() {
	p.ID = uuid.New()
}
