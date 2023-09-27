package storage

import "sync"

type Storage struct {
	sync.RWMutex
	data map[string]any
}

func NewStorage(size int) *Storage {
	return &Storage{
		data: make(map[string]any, size),
	}
}

func (storage *Storage) Set(key string, value any) {
	storage.Lock()
	defer storage.Unlock()
	storage.data[key] = value
}

func (storage *Storage) Get(key string) (any, bool) {
	storage.RLock()
	defer storage.RUnlock()

	val, ok := storage.data[key]

	return val, ok
}

func (storage *Storage) Data() map[string]any {
	return storage.data
}
