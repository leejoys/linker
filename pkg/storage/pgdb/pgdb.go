package pgdb

import (
	"context"
	"linker/pkg/storage"
	"linker/pkg/storage/generator"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Хранилище данных.
type Store struct {
	db *pgxpool.Pool
}

//New - Конструктор объекта хранилища.
func New(ctx context.Context, connstr string) (*Store, error) {

	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	// проверка связи с БД
	err = db.Ping(ctx)
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
func (s *Store) GetLong(ctx context.Context, l storage.Link) (storage.Link, error) {

	err := s.db.QueryRow(ctx,
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
func (s *Store) GetShort(ctx context.Context, l storage.Link) (storage.Link, error) {

	err := s.db.QueryRow(ctx,
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
func (s *Store) CountShort(ctx context.Context, short string) (int, error) {
	count := 0
	err := s.db.QueryRow(ctx,
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

//CountLong - проверка наличия полной ссылки
func (s *Store) CountLong(ctx context.Context, long string) (int, error) {
	count := 0
	err := s.db.QueryRow(ctx,
		`SELECT 
		count(*)
		FROM links
		WHERE longlink=$1;`, long).Scan(
		&count,
	)
	if err != nil {
		return 0, err
	}
	return count, err
}

//StoreLinkTX - сохранение новой ссылки через транзакцию
func (s *Store) StoreLink(ctx context.Context, l storage.Link) error {
	tx, err := s.db.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	count := 0

	for {
		l.ShortLink = generator.Do()
		err = tx.QueryRow(ctx,
			`SELECT 
		count(*)
		FROM links
		WHERE shortlink=$1;`, l.ShortLink).Scan(
			&count,
		)
		if err != nil {
			return err
		}
		if count == 0 {
			break
		}
	}

	_, err = tx.Exec(ctx, `
	INSERT INTO links (
		longlink,
		shortlink) 
	VALUES ($1,$2);`,
		l.LongLink,
		l.ShortLink)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return err
}
