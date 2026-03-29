package application

import (
	"sync"

	"github.com/n1jke/oop-bsuir-2025/lr-3/internal/domain"
)

type LocalCacheService struct {
	data map[string]*domain.Order
	mu   *sync.RWMutex
}

func NewLocalCacheService() *LocalCacheService {
	return &LocalCacheService{
		data: make(map[string]*domain.Order),
		mu:   new(sync.RWMutex),
	}
}

func (lcs *LocalCacheService) TryAddOrder(order *domain.Order) bool {
	lcs.mu.Lock()
	defer lcs.mu.Unlock()

	if _, exists := lcs.data[order.ID]; exists {
		return false
	}

	lcs.data[order.ID] = cloneOrder(order)

	return true
}

func (lcs *LocalCacheService) RemoveOrder(id string) {
	lcs.mu.Lock()
	defer lcs.mu.Unlock()

	delete(lcs.data, id)
}

func (lcs *LocalCacheService) FindOrder(id string) (*domain.Order, bool) {
	lcs.mu.RLock()
	defer lcs.mu.RUnlock()

	order, exist := lcs.data[id]
	if !exist {
		return nil, false
	}

	return cloneOrder(order), true
}

func cloneOrder(order *domain.Order) *domain.Order {
	if order == nil {
		return nil
	}

	orderCopy := *order
	if len(order.Items) > 0 {
		orderCopy.Items = make([]domain.Item, len(order.Items))
		copy(orderCopy.Items, order.Items)
	}

	return &orderCopy
}
