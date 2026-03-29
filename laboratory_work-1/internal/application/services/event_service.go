package services

import (
	"github.com/n1jke/oop-bsuir-2025/lr-1/internal/domain"
	"github.com/n1jke/oop-bsuir-2025/lr-1/internal/infrastructure"
)

type EventService struct {
	store *infrastructure.MemoryEventStorage
}

func NewEventService(store *infrastructure.MemoryEventStorage) *EventService {
	return &EventService{store: store}
}

func (s *EventService) Publish(e *domain.Event) error {
	return s.store.Save(e)
}

func (s *EventService) QueryAll() []*domain.Event {
	return s.store.QueryAll()
}
