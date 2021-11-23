package memdb

import (
	"errors"
	"linker/pkg/storage"
	"sync"
)

//todo записывать в две разнонаправленных мапы
type inmemory struct {
	sync.RWMutex
	LongToShort map[string]string
	ShortToLong map[string]string
}

//TODO RWMutex
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
func (s *Store) GetLong(l storage.Link) (storage.Link, error) {
	s.db.RLock()
	l.LongLink = s.db.ShortToLong[l.ShortLink]
	s.db.RUnlock()
	if l.LongLink == "" {
		return storage.Link{}, errors.New("memdb GetLong error: no data")
	}
	return l, nil
}

//GetShort - получение сокращенной ссылки по полной
func (s *Store) GetShort(l storage.Link) (storage.Link, error) {
	s.db.RLock()
	l.ShortLink = s.db.LongToShort[l.LongLink]
	s.db.RUnlock()
	if l.ShortLink == "" {
		return storage.Link{}, errors.New("memdb GetShort error: no data")
	}
	return l, nil
}

//CountShort - проверка наличия сокращенной ссылки
func (s *Store) CountShort(short string) (int, error) {
	s.db.RLock()
	if _, ok := s.db.ShortToLong[short]; !ok {
		return 0, nil
	}
	s.db.RUnlock()
	return 1, nil
}

//CountLong - проверка наличия полной ссылки
func (s *Store) CountLong(long string) (int, error) {
	s.db.RLock()
	if _, ok := s.db.LongToShort[long]; !ok {
		return 0, nil
	}
	s.db.RUnlock()
	return 1, nil
}

//StoreLink - сохранение новой ссылки
func (s *Store) StoreLink(l storage.Link) error {
	s.db.Lock()
	s.db.ShortToLong[l.ShortLink] = l.LongLink
	s.db.LongToShort[l.LongLink] = l.ShortLink
	s.db.Unlock()
	return nil
}
