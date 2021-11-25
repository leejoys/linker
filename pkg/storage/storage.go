package storage

import "context"

//Link - хранимая ссылка
type Link struct {
	LongLink  string
	ShortLink string
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	GetLong(context.Context, Link) (Link, error)     // получение полной ссылки по сокращенной
	GetShort(context.Context, Link) (Link, error)    // получение сокращенной ссылки по полной
	CountShort(context.Context, string) (int, error) // проверка наличия сокращенной ссылки
	CountLong(context.Context, string) (int, error)  // проверка наличия полной ссылки
	StoreLink(context.Context, Link) error           // сохранение ссылки
	Close()                                          // освобождение ресурса
}
