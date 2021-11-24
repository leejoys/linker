package memdb

import (
	"errors"
	"linker/pkg/storage"
	"sync"
)

//todo записывать в две разнонаправленных мапы
type inmemory struct {
	mutex       sync.Mutex
	longToShort map[string]string
	shortToLong map[string]string
	dbmap       map[string]string
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
	dbm := make(map[string]string)
	db := &inmemory{sync.Mutex{}, lts, stl, dbm}
	return &Store{db: db}
}

//Close - освобождение ресурса. Заглушка для реализации интерфейса.
func (s *Store) Close() {}

//GetLong - получение полной ссылки по сокращенной
func (s *Store) GetLong(l storage.Link) (storage.Link, error) {
	s.db.mutex.Lock()
	l.LongLink = s.db.shortToLong[l.ShortLink]
	s.db.mutex.Unlock()
	if l.LongLink == "" {
		return storage.Link{}, errors.New("memdb GetLong error: no data")
	}
	return l, nil
}

//GetShort - получение сокращенной ссылки по полной
func (s *Store) GetShort(l storage.Link) (storage.Link, error) {
	s.db.mutex.Lock()
	l.ShortLink = s.db.longToShort[l.LongLink]
	s.db.mutex.Unlock()
	if l.ShortLink == "" {
		return storage.Link{}, errors.New("memdb GetShort error: no data")
	}
	return l, nil
}

//CountShort - проверка наличия сокращенной ссылки
func (s *Store) CountShort(short string) (int, error) {
	s.db.mutex.Lock()
	defer s.db.mutex.Unlock()
	if _, ok := s.db.shortToLong[short]; !ok {
		return 0, nil
	}
	return 1, nil
}

//CountLong - проверка наличия полной ссылки
func (s *Store) CountLong(long string) (int, error) {
	s.db.mutex.Lock()
	defer s.db.mutex.Unlock()
	if _, ok := s.db.longToShort[long]; !ok {
		return 0, nil
	}
	return 1, nil
}

//StoreLink - сохранение новой ссылки
func (s *Store) StoreLink(l storage.Link) error {
	s.db.mutex.Lock()
	s.db.shortToLong[l.ShortLink] = l.LongLink
	s.db.longToShort[l.LongLink] = l.ShortLink
	s.db.mutex.Unlock()
	return nil
}
