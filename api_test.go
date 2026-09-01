package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_CreateArticleHandler(t *testing.T) {
	store := NewStorage()
	api := NewAPI(store)

	// Happy Path: валидный JSON
	body := []byte(`{"title":"HTTP Тесты","content":"Проверяем работу API эндпоинтов"}`)
	req := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	api.CreateArticleHandler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Ожидался статус 201 Created, получено: %d", rec.Code)
	}

	// Negative Path: некорректный JSON
	badReq := httptest.NewRequest(http.MethodPost, "/articles", bytes.NewBuffer([]byte(`{bad_json`)))
	badRec := httptest.NewRecorder()

	api.CreateArticleHandler(badRec, badReq)

	if badRec.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус 400 Bad Request, получено: %d", badRec.Code)
	}
}

func TestAPI_GetArticleHandler(t *testing.T) {
	store := NewStorage()
	api := NewAPI(store)

	saved, _ := store.CreateArticle(Article{
		Title:   "Поиск статьи",
		Content: "Текст статьи для GET запроса",
	})

	// Успешный запрос существующей статьи
	req := httptest.NewRequest(http.MethodGet, "/articles?id=1", nil)
	rec := httptest.NewRecorder()
	api.GetArticleHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200 OK, получено: %d", rec.Code)
	}

	// Запрос несуществующей статьи -> 404
	notFoundReq := httptest.NewRequest(http.MethodGet, "/articles?id=999", nil)
	notFoundRec := httptest.NewRecorder()
	api.GetArticleHandler(notFoundRec, notFoundReq)

	if notFoundRec.Code != http.StatusNotFound {
		t.Errorf("Ожидался статус 404 Not Found, получено: %d", notFoundRec.Code)
	}

	// Запрос с невалидным параметром ID -> 400
	invalidReq := httptest.NewRequest(http.MethodGet, "/articles?id=abc", nil)
	invalidRec := httptest.NewRecorder()
	api.GetArticleHandler(invalidRec, invalidReq)

	if invalidRec.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус 400 Bad Request, получено: %d", invalidRec.Code)
	}
	_ = saved
}

func TestAPI_AddRatingHandler(t *testing.T) {
	store := NewStorage()
	api := NewAPI(store)

	saved, _ := store.CreateArticle(Article{
		Title:   "Статья для рейтинга через API",
		Content: "Тестируем добавление оценки через HTTP POST",
	})

	// 1. Позитивный сценарий (Happy Path): валидная оценка
	bodyValid := []byte(`{"article_id": 1, "score": 5}`)
	reqValid := httptest.NewRequest(http.MethodPost, "/articles/rate", bytes.NewBuffer(bodyValid))
	recValid := httptest.NewRecorder()

	api.AddRatingHandler(recValid, reqValid)

	if recValid.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200 OK, получено: %d", recValid.Code)
	}

	// 2. Негативный сценарий: невалидная оценка (score > 5) -> 400 Bad Request
	bodyInvalidScore := []byte(`{"article_id": 1, "score": 10}`)
	reqInvalidScore := httptest.NewRequest(http.MethodPost, "/articles/rate", bytes.NewBuffer(bodyInvalidScore))
	recInvalidScore := httptest.NewRecorder()

	api.AddRatingHandler(recInvalidScore, reqInvalidScore)

	if recInvalidScore.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус 400 Bad Request при невалидном score, получено: %d", recInvalidScore.Code)
	}

	// 3. Негативный сценарий: битый JSON -> 400 Bad Request
	bodyBadJSON := []byte(`{"article_id": 1, "score": `)
	reqBadJSON := httptest.NewRequest(http.MethodPost, "/articles/rate", bytes.NewBuffer(bodyBadJSON))
	recBadJSON := httptest.NewRecorder()

	api.AddRatingHandler(recBadJSON, reqBadJSON)

	if recBadJSON.Code != http.StatusBadRequest {
		t.Errorf("Ожидался статус 400 Bad Request при битом JSON, получено: %d", recBadJSON.Code)
	}
	_ = saved
}
