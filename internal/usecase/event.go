package usecase

import (
	"context"
	"homework/internal/domain"
)

type Event struct {
	eventRepo  EventRepository
	sensorRepo SensorRepository
}

func NewEvent(er EventRepository, sr SensorRepository) *Event {
	return &Event{
		eventRepo:  er,
		sensorRepo: sr,
	}
}

func (e *Event) ReceiveEvent(ctx context.Context, event *domain.Event) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if e.eventRepo == nil && e.sensorRepo == nil {
		return ErrInvalidEventTimestamp
	}

	sensor, err := e.sensorRepo.GetSensorBySerialNumber(ctx, event.SensorSerialNumber)
	if err != nil {
		return err
	}

	event.SensorID = sensor.ID

	if err := e.eventRepo.SaveEvent(ctx, event); err != nil {
		return err
	}

	sensor.LastActivity = event.Timestamp
	sensor.CurrentState = event.Payload

	return e.sensorRepo.SaveSensor(ctx, sensor)
}

func (e *Event) GetLastEventBySensorID(ctx context.Context, sensorID int64) (*domain.Event, error) {
	return e.eventRepo.GetLastEventBySensorID(ctx, sensorID)
}
