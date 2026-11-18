package usecase

import (
	"context"

	"github.com/Mikee01-dev/ruangku-backend/internal/domain"
)

type roomUsecase struct {
	repo domain.RoomRepository
}

func NewRoomUsecase(repo domain.RoomRepository) domain.RoomUsecase {
	return &roomUsecase{repo: repo}
}