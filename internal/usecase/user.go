package usecase

import (
	"context"
	"homework/internal/domain"
)

type User struct {
	userRepo        UserRepository
	sensorRepo      SensorRepository
	sensorOwnerRepo SensorOwnerRepository
}

func NewUser(ur UserRepository, sor SensorOwnerRepository, sr SensorRepository) *User {
	return &User{
		userRepo:        ur,
		sensorRepo:      sr,
		sensorOwnerRepo: sor,
	}
}

func (u *User) RegisterUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	if user.Name == "" {
		return nil, ErrInvalidUserName
	}
	if err := u.userRepo.SaveUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) AttachSensorToUser(ctx context.Context, userID, sensorID int64) error {
	_, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	_, err = u.sensorRepo.GetSensorByID(ctx, sensorID)
	if err != nil {
		return err
	}

	return u.sensorOwnerRepo.SaveSensorOwner(ctx, domain.SensorOwner{UserID: userID, SensorID: sensorID})
}

func (u *User) GetUserSensors(ctx context.Context, userID int64) ([]domain.Sensor, error) {
	_, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	sensors, err := u.sensorOwnerRepo.GetSensorsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var result []domain.Sensor
	for _, sensor := range sensors {
		buf, err := u.sensorRepo.GetSensorByID(ctx, sensor.SensorID)
		if err != nil {
			return nil, err
		}
		result = append(result, *buf)
	}

	return result, nil
}
