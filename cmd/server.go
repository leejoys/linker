package main

import (
	"context"
	"linker/pkg/api"
	"linker/pkg/storage"
	"linker/pkg/storage/generator"
	"linker/pkg/storage/memdb"
	"linker/pkg/storage/pgdb"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	alphabet    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
	shortLength = 10
)

// Сервер линкера.
type server struct {
	db  storage.Interface
	api *api.API
}

func dbFabric(ctx context.Context, g *generator.Generator, inmemory bool) storage.Interface {
	if inmemory {
		return memdb.New(g)
	}
	//  Создаём объект базы данных PostgreSQL.
	pwd := os.Getenv("PGPASS")
	user := os.Getenv("PGUSER")
	addr := os.Getenv("PGADDR")
	//postgres://user:pwd@postgres:5432/db
	connstr := "postgres://" + user + ":" + pwd + "@" + addr
	db, err := pgdb.New(ctx, g, connstr)
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

	//todo gracefull shutdown
	// Создаём объект сервера
	srv := server{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	//Создаем генератор
	g := generator.New(alphabet, shortLength)

	// Инициализируем БД
	srv.db = dbFabric(ctx, g, isMemdb)

	// Освобождаем ресурс
	defer srv.db.Close()

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(ctx, srv.db)

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
