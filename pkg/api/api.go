package api

import (
	"linker/pkg/storage"
	"net/http"

	"github.com/gorilla/mux"
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
	api.r.HandleFunc("/links/detailed", api.detailed).Methods(http.MethodGet)
	//метод добавления ссылки
	api.r.HandleFunc("/links/store", api.storeLink).Methods(http.MethodPost)
}
