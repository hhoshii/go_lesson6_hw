package inmemory

import (
	"context"
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
	"time"
)

type SensorRepository struct {
	mu                    sync.RWMutex
	nextID                int64
	sensorsById           map[int64]*domain.Sensor
	sensorsBySerialNumber map[string]*domain.Sensor
}

func NewSensorRepository() *SensorRepository {
	return &SensorRepository{
		nextID:                1,
		sensorsById:           make(map[int64]*domain.Sensor),
		sensorsBySerialNumber: make(map[string]*domain.Sensor),
	}
}

func (r *SensorRepository) SaveSensor(ctx context.Context, sensor *domain.Sensor) error {
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

	if sensor == nil {
		return fmt.Errorf("sensor is nil")
	}

	if _, ok := r.sensorsBySerialNumber[sensor.SerialNumber]; ok {
		return fmt.Errorf("sensor with serial number %s already exists", sensor.SerialNumber)
	}

	sensor.RegisteredAt = time.Now()
	sensor.ID = r.nextID
	r.nextID++
	r.sensorsById[sensor.ID] = sensor
	r.sensorsBySerialNumber[sensor.SerialNumber] = sensor
	return nil
}

func (r *SensorRepository) GetSensors(ctx context.Context) ([]domain.Sensor, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	var sensors []domain.Sensor
	for _, sensor := range r.sensorsById {
		sensors = append(sensors, *sensor)
	}
	return sensors, nil
}

func (r *SensorRepository) GetSensorByID(ctx context.Context, id int64) (*domain.Sensor, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	sensor, ok := r.sensorsById[id]
	if !ok {
		return nil, usecase.ErrSensorNotFound
	}
	return sensor, nil
}

func (r *SensorRepository) GetSensorBySerialNumber(ctx context.Context, sn string) (*domain.Sensor, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	sensor, ok := r.sensorsBySerialNumber[sn]
	if !ok {
		return nil, usecase.ErrSensorNotFound
	}

	return sensor, nil
}
