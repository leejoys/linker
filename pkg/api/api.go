package api

import (
	"fmt"
	"io"
	"linker/pkg/storage"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const (
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
	shortLength = 10
)

// Программный интерфейс сервиса
type API struct {
	db storage.Interface
	r  *mux.Router
}

// Конструктор объекта API
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	//метод получения полной ссылки
	api.r.HandleFunc("/links/{short}", api.getLink).Methods(http.MethodGet)
	//метод добавления ссылки
	api.r.HandleFunc("/links", api.storeLink).Methods(http.MethodPost)
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.r
}

// получение полной ссылки по сокращенной
func (api *API) getLink(w http.ResponseWriter, r *http.Request) {
	l := storage.Link{}
	l.ShortLink = mux.Vars(r)["short"]
	if len(l.ShortLink) != 10 {
		http.Error(w, "invalid link", http.StatusBadRequest)
		return
	}
	l, err := api.db.GetLong(l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(l.LongLink))
}

// сохранение ссылки
func (api *API) storeLink(w http.ResponseWriter, r *http.Request) {
	bLongLink, err := io.ReadAll(r.Body)
	if err != nil || len(bLongLink) < 1 {
		http.Error(w, fmt.Sprintf("storeLink ReadAll error: %s", err.Error()), http.StatusBadRequest)
		return
	}

	l := storage.Link{LongLink: string(bLongLink)}

	count, err := api.db.CountLong(l.LongLink)
	if err != nil {
		http.Error(w, fmt.Sprintf("storeLink CountLong error: %s", err.Error()), http.StatusBadRequest)
		return
	}
	if count > 0 {
		exist, err := api.db.GetShort(l)
		if err != nil {
			http.Error(w, fmt.Sprintf("storeLink GetShort error: %s", err.Error()), http.StatusBadRequest)
			return
		}
		w.Write([]byte(exist.ShortLink))
		return
	}

	rand.Seed(time.Now().UnixNano())
	abc := []byte(alphabet)
	var short []byte
	for {
		short = []byte{}
		for i := 1; i <= shortLength; i++ {
			short = append(short, abc[rand.Intn(len(abc)-1)])
		}
		count, err := api.db.CountShort(string(short))
		if err != nil {
			http.Error(w, fmt.Sprintf("storeLink CountShort error: %s", err.Error()), http.StatusBadRequest)
			return
		}
		if count > 0 {
			continue
		}
		break
	}
	l.ShortLink = string(short)
	err = api.db.StoreLink(l)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(l.ShortLink))
}
