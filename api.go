package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// API инкапсулирует хранилище и обработчики HTTP-запросов
type API struct {
	storage *Storage
}

// NewAPI создает новый экземпляр API
func NewAPI(storage *Storage) *API {
	return &API{storage: storage}
}

// CreateArticleHandler обрабатывает POST /articles
func (a *API) CreateArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var article Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	saved, err := a.storage.CreateArticle(article)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(saved)
}

// GetArticleHandler обрабатывает GET /articles?id=1
func (a *API) GetArticleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Некорректный или отсутствующий параметр id", http.StatusBadRequest)
		return
	}

	article, err := a.storage.GetArticleByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(article)
}

// AddRatingHandler обрабатывает POST /articles/rate
func (a *API) AddRatingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var rating Rating
	if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	avgRating, err := a.storage.AddRating(rating)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"article_id": rating.ArticleID,
		"new_rating": avgRating,
	})
}
