package inmemory

import (
	"context"
	"homework/internal/domain"
	"sync"
)

type SensorOwnerRepository struct {
	mu             sync.RWMutex
	sensorsByOwner map[int64][]int64
}

func NewSensorOwnerRepository() *SensorOwnerRepository {
	return &SensorOwnerRepository{
		sensorsByOwner: make(map[int64][]int64),
	}
}

func (r *SensorOwnerRepository) SaveSensorOwner(ctx context.Context, sensorOwner domain.SensorOwner) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sensorsByOwner[sensorOwner.UserID] = append(r.sensorsByOwner[sensorOwner.UserID], sensorOwner.SensorID)
	return nil
}

func (r *SensorOwnerRepository) GetSensorsByUserID(ctx context.Context, userID int64) ([]domain.SensorOwner, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []domain.SensorOwner
	for _, v := range r.sensorsByOwner[userID] {
		result = append(result, domain.SensorOwner{UserID: userID, SensorID: v})
	}
	return result, nil
}
