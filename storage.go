package main

import (
	"errors"
	"sync"
	"time"
)

type Storage struct {
	mu         sync.RWMutex
	articles   map[int64]Article
	ratings    map[int64][]Rating
	nextArtID  int64
	nextRateID int64
}

func NewStorage() *Storage {
	return &Storage{
		articles:   make(map[int64]Article),
		ratings:    make(map[int64][]Rating),
		nextArtID:  1,
		nextRateID: 1,
	}
}

func (s *Storage) CreateArticle(article Article) (Article, error) {
	if err := article.Validate(); err != nil {
		return Article{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	article.ID = s.nextArtID
	article.CreatedAt = time.Now()
	article.Rating = 0

	s.articles[article.ID] = article
	s.nextArtID++

	return article, nil
}

func (s *Storage) GetArticleByID(id int64) (Article, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	article, exists := s.articles[id]
	if !exists {
		return Article{}, errors.New("статья не найдена")
	}
	return article, nil
}

func (s *Storage) AddRating(rating Rating) (float64, error) {
	if err := rating.Validate(); err != nil {
		return 0, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	article, exists := s.articles[rating.ArticleID]
	if !exists {
		return 0, errors.New("статья для оценки не найдена")
	}

	rating.ID = s.nextRateID
	s.ratings[rating.ArticleID] = append(s.ratings[rating.ArticleID], rating)
	s.nextRateID++

	var sum float64
	allRatings := s.ratings[rating.ArticleID]
	for _, r := range allRatings {
		sum += float64(r.Score)
	}
	avgRating := sum / float64(len(allRatings))

	article.Rating = avgRating
	s.articles[article.ID] = article

	return avgRating, nil
}
