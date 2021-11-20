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

func dbFabric(inmemory bool) storage.Interface {
	if inmemory {
		return memdb.New()
	}
	//  Создаём объект базы данных PostgreSQL.
	// pwd := os.Getenv("PGPASS")
	// user := os.Getenv("PGUSER")
	// addr := os.Getenv("PGADDR")
	// //user://postgres:pwd@postgres:5432/db
	// connstr := user + "://postgres:" + pwd + "@" + addr
	connstr := "test://postgres:test@pg_db/test"
	db, err := pgdb.New(connstr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {

	if len(os.Args) > 1 && os.Args[1] != "-inmemory" {
		log.Fatal("usage: server [-inmemory]")
	}
	isMemdb := len(os.Args) > 1 && os.Args[1] == "-inmemory"

	// Создаём объект сервера
	srv := server{}

	// Инициализируем БД
	srv.db = dbFabric(isMemdb)

	// Освобождаем ресурс
	defer srv.db.Close()

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов.
	go func() {
		log.Fatal(http.ListenAndServe("0.0.0.0:8080", srv.api.Router()))
	}()
	log.Println("HTTP server is started")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	log.Println("HTTP server has been stopped")
}
