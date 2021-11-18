package storage

//Link - хранимая ссылка
type Link struct {
	FullLink  string
	ShortLink string
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	StoreLink(int) (Link, error) // получение ссылки
	AddLink(Link) error          // сохранение ссылки
	Close()                      // освобождение ресурса
}
