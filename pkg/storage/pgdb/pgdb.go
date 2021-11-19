package pgdb

import (
	"context"
	"linker/pkg/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

//New - Конструктор объекта хранилища.
func New(connstr string) (*Store, error) {

	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	// проверка связи с БД
	err = db.Ping(context.Background())
	if err != nil {
		db.Close()
		return nil, err
	}

	return &Store{db: db}, nil
}

//Close - освобождение ресурса
func (s *Store) Close() {
	s.db.Close()
}

//GetLong - получение полной ссылки по сокращенной
func (s *Store) GetLong(l storage.Link) (storage.Link, error) {

	err := s.db.QueryRow(context.Background(),
		`SELECT 
	links.longlink
	FROM links
	WHERE shortlink=$1;`, l.ShortLink).Scan(
		&l.LongLink,
	)

	if err != nil {
		return storage.Link{}, err
	}

	return l, err
}

//GetShort - получение сокращенной ссылки по полной
func (s *Store) GetShort(l storage.Link) (storage.Link, error) {

	err := s.db.QueryRow(context.Background(),
		`SELECT 
	links.shortlink
	FROM links
	WHERE longlink=$1;`, l.LongLink).Scan(
		&l.ShortLink,
	)

	if err != nil {
		return storage.Link{}, err
	}

	return l, err
}

//CountShort - проверка наличия сокращенной ссылки
func (s *Store) CountShort(short string) (int, error) {
	count := 0
	err := s.db.QueryRow(context.Background(),
		`SELECT 
		count(*)
		FROM links
		WHERE shortlink=$1;`, short).Scan(
		&count,
	)
	if err != nil {
		return 0, err
	}
	return count, err
}

//CountLong - получение сокращенной ссылки по полной
func (s *Store) CountLong(short string) (int, error) {
	count := 0
	err := s.db.QueryRow(context.Background(),
		`SELECT 
		count(*)
		FROM links
		WHERE longlink=$1;`, short).Scan(
		&count,
	)
	if err != nil {
		return 0, err
	}
	return count, err
}

//StoreLink - сохранение новой ссылки
func (s *Store) StoreLink(l storage.Link) error {
	_, err := s.db.Exec(context.Background(), `
	INSERT INTO links (
		longlink,
		shortlink) 
	VALUES ($1,$2);`,
		l.LongLink,
		l.ShortLink)

	return err
}
