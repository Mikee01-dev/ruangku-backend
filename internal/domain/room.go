package domain

import (
	"context"
	"errors"
	"time"
)

type Room struct {
	ID        int64     `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Building  string    `json:"building"`
	Capacity  int       `json:"capacity"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrRoomNotFound     = errors.New("room not found")
	ErrInvalidRoomInput = errors.New("invalid room input")
)

type RoomRepository interface {
	FindAll(ctx context.Context) ([]Room, error)
	FindByID(ctx context.Context, id int64) (*Room, error)
	Create(ctx context.Context, room *Room) error
	Update(ctx context.Context, room *Room) error
	Delete(ctx context.Context, id int64) error
}

type RoomUsecase interface {
	GetAll(ctx context.Context) ([]Room, error)
	GetByID(ctx context.Context, id int64) (*Room, error)
	Create(ctx context.Context, room *Room) error
	Update(ctx context.Context, room *Room) error
	Delete(ctx context.Context, id int64) error
}