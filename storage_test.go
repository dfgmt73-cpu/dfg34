package main

import "testing"

func TestStorage_CreateAndGet(t *testing.T) {
	store := NewStorage()

	article := Article{
		Title:   "Тестирование в Go",
		Content: "Пишем надежные Unit-тесты для сервиса.",
		Tags:    []string{"testing", "golang"},
	}

	saved, err := store.CreateArticle(article)
	if err != nil {
		t.Fatalf("Не удалось создать статью: %v", err)
	}

	if saved.ID != 1 {
		t.Errorf("Ожидался ID=1, получено: %d", saved.ID)
	}

	fetched, err := store.GetArticleByID(saved.ID)
	if err != nil {
		t.Fatalf("Ошибка при получении статьи: %v", err)
	}

	if fetched.Title != article.Title {
		t.Errorf("Ожидался заголовок '%s', получено: '%s'", article.Title, fetched.Title)
	}

	_, err = store.GetArticleByID(999)
	if err == nil {
		t.Error("Ожидалась ошибка для несуществующей статьи, но получено nil")
	}
}

func TestStorage_AddRating_AverageCalculation(t *testing.T) {
	store := NewStorage()

	saved, err := store.CreateArticle(Article{
		Title:   "Статья для оценки",
		Content: "Проверяем математику среднего рейтинга.",
	})
	if err != nil {
		t.Fatalf("Ошибка подготовки данных: %v", err)
	}

	_, err = store.AddRating(Rating{ArticleID: saved.ID, Score: 5})
	if err != nil {
		t.Fatalf("Ошибка добавления первой оценки: %v", err)
	}

	avg, err := store.AddRating(Rating{ArticleID: saved.ID, Score: 4})
	if err != nil {
		t.Fatalf("Ошибка добавления второй оценки: %v", err)
	}

	expectedAvg := 4.5
	if avg != expectedAvg {
		t.Errorf("Ожидался средний рейтинг %.2f, получено: %.2f", expectedAvg, avg)
	}

	_, err = store.AddRating(Rating{ArticleID: 999, Score: 5})
	if err == nil {
		t.Error("Ожидалась ошибка при оценке несуществующей статьи, но получено nil")
	}
}
