package api

import (
	"context"
	"io/ioutil"
	"linker/pkg/storage"
	"linker/pkg/storage/memdb"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPI_storeLink(t *testing.T) {
	// Создаём чистый объект API для теста.
	dbase := memdb.New()
	defer dbase.Close()

	api := New(context.Background(), dbase)
	// Создаём HTTP-запрос.
	req := httptest.NewRequest(http.MethodPost, "/links",
		strings.NewReader("https://habr.com/ru/news/t/568128/"))
	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	api.r.ServeHTTP(rr, req)

	// Проверяем код ответа.
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Читаем тело ответа.
	b, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось прочитать ответ сервера: %v", err)
	}
	// проверяем длину ответа
	if len(b) != 10 {
		t.Fatalf("короткая ссылка имеет неверный формат: %s", string(b))
	}

	// Проверяем, что второй раз вернется тот же результат.
	wantArr := b
	req = httptest.NewRequest(http.MethodPost, "/links",
		strings.NewReader("https://habr.com/ru/news/t/568128/"))
	rr = httptest.NewRecorder()
	api.r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	b, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось прочитать ответ сервера: %v", err)
	}
	if len(b) != 10 {
		t.Fatalf("короткая ссылка имеет неверный формат: %s", string(b))
	}
	if string(b) != string(wantArr) {
		t.Fatalf("получено %s, ожидалось %s", string(b), string(wantArr))
	}

	// Проверяем что пустая ссылка не сохраняется.
	req = httptest.NewRequest(http.MethodPost, "/links", strings.NewReader(""))
	rr = httptest.NewRecorder()
	api.r.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusBadRequest)
	}

}

func TestAPI_getLink(t *testing.T) {
	// Создаём чистый объект API для теста.
	dbase := memdb.New()
	defer dbase.Close()
	ctx := context.Background()

	links := []storage.Link{
		{ShortLink: "1234567890",
			LongLink: "https://github.com/microsoft/CBL-Mariner"},
		{ShortLink: "0987654321",
			LongLink: "https://habr.com/ru/news/t/568128/"},
	}
	for _, l := range links {
		dbase.StoreLink(ctx, l)
	}

	api := New(ctx, dbase)
	// Создаём HTTP-запрос.
	req := httptest.NewRequest(http.MethodGet, "/links/1234567890", nil)
	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()
	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	api.r.ServeHTTP(rr, req)

	// Проверяем код ответа.
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Читаем тело ответа.
	b, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось прочитать ответ сервера: %v", err)
	}
	// проверяем ответ
	if string(b) != links[0].LongLink {
		t.Fatalf("длинная ссылка не соответствует короткой, получили: %s, ожидали: %s",
			string(b), links[0].LongLink)
	}

	// Получаем второй ответ.
	req = httptest.NewRequest(http.MethodGet, "/links/0987654321", nil)
	rr = httptest.NewRecorder()
	api.r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	b, err = ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось прочитать ответ сервера: %v", err)
	}
	// проверяем ответ
	if string(b) != links[1].LongLink {
		t.Fatalf("длинная ссылка не соответствует короткой, получили: %s, ожидали: %s",
			string(b), links[1].LongLink)
	}

}
