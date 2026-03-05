package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Event - event for future log system.
type Event struct {
	id          uuid.UUID
	aggregateID uuid.UUID
	eventType   string
	occurredAt  time.Time
}

func NewEvent(id, aggregateID uuid.UUID, eventType string) *Event {
	return &Event{
		id:          id,
		aggregateID: aggregateID,
		eventType:   eventType,
		occurredAt:  time.Now(),
	}
}

func (e Event) ID() uuid.UUID {
	return e.id
}

func (e Event) AggregateID() uuid.UUID {
	return e.aggregateID
}

func (e Event) Type() string {
	return e.eventType
}

func (e Event) OccurredAt() time.Time {
	return e.occurredAt
}

func (e Event) String() string {
	return fmt.Sprintf(
		"Event{id: %s, aggregateID: %s, type: %q, occurredAt: %s}",
		e.id.String(), e.aggregateID.String(), e.eventType, e.occurredAt.Format(time.RFC3339))
}
