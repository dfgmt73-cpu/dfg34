package main

import (
	"errors"
	"strings"
	"time"
)

type Article struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	Rating    float64   `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
}

type Rating struct {
	ID        int64 `json:"id"`
	ArticleID int64 `json:"article_id"`
	Score     int   `json:"score"`
}

func (a *Article) Validate() error {
	trimmedTitle := strings.TrimSpace(a.Title)
	if len(trimmedTitle) < 3 {
		return errors.New("заголовок должен содержать минимум 3 символа")
	}
	if strings.TrimSpace(a.Content) == "" {
		return errors.New("содержимое статьи не может быть пустым")
	}
	return nil
}

func (r *Rating) Validate() error {
	if r.ArticleID <= 0 {
		return errors.New("некорректный ID статьи")
	}
	if r.Score < 1 || r.Score > 5 {
		return errors.New("оценка должна быть в диапазоне от 1 до 5")
	}
	return nil
}
