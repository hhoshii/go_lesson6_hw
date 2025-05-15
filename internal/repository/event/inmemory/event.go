package inmemory

import (
	"context"
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
)

type EventRepository struct {
	mu     sync.Mutex
	events map[int64][]*domain.Event
}

func NewEventRepository() *EventRepository {
	return &EventRepository{
		events: make(map[int64][]*domain.Event),
	}
}

func (r *EventRepository) SaveEvent(ctx context.Context, event *domain.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return err
	}

	if event == nil {
		return fmt.Errorf("event is nil")
	}

	r.events[event.SensorID] = append(r.events[event.SensorID], event)
	return nil
}

func (r *EventRepository) GetLastEventBySensorID(ctx context.Context, id int64) (*domain.Event, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if err := ctx.Err(); err != nil {
		return nil, err
	}

	events, ok := r.events[id]
	if !ok || len(events) == 0 {
		return nil, usecase.ErrEventNotFound
	}

	result := events[0]
	for i := 0; i < len(events); i++ {
		if events[i].Timestamp.After(result.Timestamp) {
			result = events[i]
		}
	}

	return result, nil
}
