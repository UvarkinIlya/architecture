package storage

import (
	"sync"
)

type Storage interface {
	MessageStorage
	UserStorage
}

type StorageImpl struct {
	MessagesFilepath string
	UserFilepath     string
	lock             sync.Mutex
}

func NewStorageImpl(messageFilepath, userFilepath string) *StorageImpl {
	return &StorageImpl{
		MessagesFilepath: messageFilepath,
		UserFilepath:     userFilepath,
	}
}
