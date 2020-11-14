package storage

import (
	"sort"
	"sync"
	"time"
	"users-api/models"
)

type inMemoryStorage struct {
	storage map[int]*models.User
	mu      *sync.RWMutex
}

func NewInMemoryStorage() *inMemoryStorage {
	return &inMemoryStorage{
		storage: map[int]*models.User{},
		mu:      &sync.RWMutex{},
	}
}

func (i inMemoryStorage) Get(id int) (*models.User, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	user, exists := i.storage[id]
	if !exists {
		return user, ErrNotExists
	}

	return user, nil
}

func (i inMemoryStorage) GetAll() ([]*models.User, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	//ordering result
	keys := make([]int, 0)
	for k := range i.storage {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	users := make([]*models.User, 0)
	for _, k := range keys {
		users = append(users, i.storage[k])
	}

	return users, nil
}

func (i inMemoryStorage) Store(newUser *models.User) (newId int, err error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	newId = len(i.storage) + 1
	newUser.ID = newId
	newUser.CreatedAt = time.Now()

	i.storage[newId] = newUser
	return
}

func (i inMemoryStorage) Update(id int, userData *models.User) (*models.User, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	userData.ID = id
	i.storage[id] = userData

	return userData, nil
}

func (i inMemoryStorage) Delete(id int) error {
	i.mu.RLock()
	defer i.mu.RUnlock()

	delete(i.storage, id)

	return nil
}

func (i inMemoryStorage) UserExist(user *models.User) (bool, error) {
	for _, storedUser := range i.storage {
		if storedUser.Last == user.Last && storedUser.First == user.First {
			return true, nil
		}
	}
	return false, nil
}
