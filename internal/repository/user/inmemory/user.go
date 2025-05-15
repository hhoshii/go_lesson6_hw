package inmemory

import (
	"context"
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecase"
	"sync"
)

type UserRepository struct {
	mu     sync.RWMutex
	nextID int64
	users  map[int64]*domain.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		nextID: 1,
		users:  make(map[int64]*domain.User),
	}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if user == nil {
		return fmt.Errorf("user is nil")
	}

	r.users[user.ID] = user
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, usecase.ErrUserNotFound
	}

	return user, nil
}
