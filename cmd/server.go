package main

import (
	"linker/pkg/api"
	"linker/pkg/storage"
	"linker/pkg/storage/pgdb"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// Сервер линкера.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// Создаём объект сервера
	srv := server{}

	//  Создаём объект базы данных PostgreSQL.
	pwd := os.Getenv("pgpass")
	connstr := "postgres://postgres:" + pwd + "@0.0.0.0/linksrv"
	db, err := pgdb.New(connstr)
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем хранилище сервера БД
	srv.db = db

	// Освобождаем ресурс
	defer srv.db.Close()

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов.
	go func() {
		log.Fatal(http.ListenAndServe("localhost:8080", srv.api.Router()))
	}()
	log.Println("HTTP server is started on localhost:8080")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	log.Println("HTTP server has been stopped")
}
