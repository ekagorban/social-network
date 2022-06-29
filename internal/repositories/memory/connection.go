package memory

import "sync"

type Store struct {
	storageData   sync.Map
	storageAccess sync.Map
}

func New() *Store {
	return &Store{storageData: sync.Map{}, storageAccess: sync.Map{}}
}
