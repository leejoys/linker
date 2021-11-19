package storage

//Link - хранимая ссылка
type Link struct {
	LongLink  string
	ShortLink string
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	GetLong(Link) (Link, error)     // получение полной ссылки по сокращенной
	GetShort(Link) (Link, error)    // получение сокращенной ссылки по полной
	CountShort(string) (int, error) // проверка наличия сокращенной ссылки
	CountLong(string) (int, error)  // проверка наличия полной ссылки
	StoreLink(Link) error           // сохранение ссылки
	Close()                         // освобождение ресурса
}
