package memdb

import (
	"context"
	"errors"
	"linker/pkg/storage"
	"linker/pkg/storage/generator"
	"sync"
)

type inmemory struct {
	mutex       sync.RWMutex
	longToShort map[string]string
	shortToLong map[string]string
}

// Хранилище данных.
type Store struct {
	db *inmemory
}

//New - Конструктор объекта хранилища.
func New() *Store {
	lts := make(map[string]string)
	stl := make(map[string]string)
	db := &inmemory{sync.RWMutex{}, lts, stl}
	return &Store{db: db}
}

//Close - освобождение ресурса. Заглушка для реализации интерфейса.
func (s *Store) Close() {}

//GetLong - получение полной ссылки по сокращенной
func (s *Store) GetLong(ctx context.Context, l storage.Link) (storage.Link, error) {
	s.db.mutex.RLock()
	l.LongLink = s.db.shortToLong[l.ShortLink]
	s.db.mutex.RUnlock()
	if l.LongLink == "" {
		return storage.Link{}, errors.New("memdb GetLong error: no data")
	}
	return l, nil
}

//GetShort - получение сокращенной ссылки по полной
func (s *Store) GetShort(ctx context.Context, l storage.Link) (storage.Link, error) {
	s.db.mutex.RLock()
	l.ShortLink = s.db.longToShort[l.LongLink]
	s.db.mutex.RUnlock()
	if l.ShortLink == "" {
		return storage.Link{}, errors.New("memdb GetShort error: no data")
	}
	return l, nil
}

//CountShort - проверка наличия сокращенной ссылки
func (s *Store) CountShort(ctx context.Context, short string) (int, error) {
	s.db.mutex.RLock()
	defer s.db.mutex.RUnlock()
	if _, ok := s.db.shortToLong[short]; !ok {
		return 0, nil
	}
	return 1, nil
}

//CountLong - проверка наличия полной ссылки
func (s *Store) CountLong(ctx context.Context, long string) (int, error) {
	s.db.mutex.RLock()
	defer s.db.mutex.RUnlock()
	if _, ok := s.db.longToShort[long]; !ok {
		return 0, nil
	}
	return 1, nil
}

//StoreLink - сохранение новой ссылки
func (s *Store) StoreLink(ctx context.Context, l storage.Link) error {
	s.db.mutex.Lock()
	for {
		l.ShortLink = generator.Do()
		if _, ok := s.db.shortToLong[l.ShortLink]; !ok {
			break
		}
	}
	s.db.shortToLong[l.ShortLink] = l.LongLink
	s.db.longToShort[l.LongLink] = l.ShortLink
	s.db.mutex.Unlock()
	return nil
}
