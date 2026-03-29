package infrastructure

import (
	"errors"
	"sync"

	"github.com/google/uuid"

	"github.com/n1jke/oop-bsuir-2025/lr-2/internal/application/services"
	"github.com/n1jke/oop-bsuir-2025/lr-2/internal/domain"
)

var ErrNotFound = errors.New("not found")

type CacheStorage[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewCacheStorage[K comparable, V any]() *CacheStorage[K, V] {
	return &CacheStorage[K, V]{
		mu:   sync.RWMutex{},
		data: make(map[K]V),
	}
}

func (cs *CacheStorage[K, V]) Save(k K, v V) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.data[k] = v

	return nil
}

func (cs *CacheStorage[K, V]) ByID(k K) (V, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	var zero V

	val, exist := cs.data[k]
	if !exist {
		return zero, ErrNotFound
	}

	return val, nil
}

func (cs *CacheStorage[K, V]) Delete(k K) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if _, exists := cs.data[k]; !exists {
		return ErrNotFound
	}

	delete(cs.data, k)

	return nil
}

type AccountCache struct {
	*CacheStorage[uuid.UUID, services.PaymentAccount]
}

func NewAccountCacheStorage() *AccountCache {
	cs := NewCacheStorage[uuid.UUID, services.PaymentAccount]()
	return &AccountCache{cs}
}

func NewMemoryAccountStorage() *AccountCache {
	return NewAccountCacheStorage()
}

func (ac *AccountCache) UpdateStatus(accountID uuid.UUID, status domain.AccountStatus) error {
	ac.mu.Lock()
	defer ac.mu.Unlock()

	account, exists := ac.data[accountID]
	if !exists {
		return ErrNotFound
	}

	account.ChangeStatus(status)
	ac.data[accountID] = account

	return nil
}

type EventCache struct {
	*CacheStorage[uuid.UUID, domain.Event]
}

func NewEventCacheStorage() *EventCache {
	cs := NewCacheStorage[uuid.UUID, domain.Event]()
	return &EventCache{cs}
}

func NewMemoryEventStorage() *EventCache {
	return NewEventCacheStorage()
}

func (ec *EventCache) QueryAll() []*domain.Event {
	ec.mu.RLock()
	defer ec.mu.RUnlock()

	events := make([]*domain.Event, 0, len(ec.data))
	for _, v := range ec.data {
		event := v
		events = append(events, &event)
	}

	return events
}
