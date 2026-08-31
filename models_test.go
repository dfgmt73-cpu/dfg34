package main

import "testing"

func TestArticle_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		article     Article
		expectError bool
	}{
		{
			name: "Корректная статья (Happy Path)",
			article: Article{
				Title:   "Основы Go для QA",
				Content: "Разбираем написание автотестов.",
			},
			expectError: false,
		},
		{
			name: "Граничное значение: заголовок ровно 3 символа (валидно)",
			article: Article{
				Title:   "Gol",
				Content: "Контент статьи...",
			},
			expectError: false,
		},
		{
			name: "Граничное значение: заголовок 2 символа (невалидно)",
			article: Article{
				Title:   "Go",
				Content: "Контент статьи...",
			},
			expectError: true,
		},
		{
			name: "Пустой заголовок (невалидно)",
			article: Article{
				Title:   "",
				Content: "Контент статьи...",
			},
			expectError: true,
		},
		{
			name: "Пустой контент (невалидно)",
			article: Article{
				Title:   "Длинный заголовок",
				Content: "",
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.article.Validate()
			if tc.expectError && err == nil {
				t.Errorf("Ожидалась ошибка, но получено nil для кейса: %s", tc.name)
			}
			if !tc.expectError && err != nil {
				t.Errorf("Не ожидалась ошибка, но получено: %v для кейса: %s", err, tc.name)
			}
		})
	}
}

func TestRating_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		rating      Rating
		expectError bool
	}{
		{
			name:        "Корректная оценка: 5",
			rating:      Rating{ArticleID: 1, Score: 5},
			expectError: false,
		},
		{
			name:        "Граничное значение снизу: 1 (валидно)",
			rating:      Rating{ArticleID: 1, Score: 1},
			expectError: false,
		},
		{
			name:        "Выход за нижнюю границу: 0 (невалидно)",
			rating:      Rating{ArticleID: 1, Score: 0},
			expectError: true,
		},
		{
			name:        "Выход за верхнюю границу: 6 (невалидно)",
			rating:      Rating{ArticleID: 1, Score: 6},
			expectError: true,
		},
		{
			name:        "Некорректный ArticleID: 0 (невалидно)",
			rating:      Rating{ArticleID: 0, Score: 4},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.rating.Validate()
			if tc.expectError && err == nil {
				t.Errorf("Ожидалась ошибка, но получено nil для кейса: %s", tc.name)
			}
			if !tc.expectError && err != nil {
				t.Errorf("Не ожидалась ошибка, но получено: %v для кейса: %s", err, tc.name)
			}
		})
	}
}
