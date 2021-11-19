package main

import (
	"linker/pkg/api"
	"linker/pkg/storage"
	"linker/pkg/storage/memdb"
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

func dbfabric(s string) storage.Interface {
	if s == "memdb" {
		return memdb.New()
	}
	//  Создаём объект базы данных PostgreSQL.
	pwd := os.Getenv("pgpass")
	connstr := "postgres://postgres:" + pwd + "@0.0.0.0/linker"
	db, err := pgdb.New(connstr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	// Создаём объект сервера
	srv := server{}

	// Инициализируем хранилище сервера БД
	srv.db = dbfabric("memdb")

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
