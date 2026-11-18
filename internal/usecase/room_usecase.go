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

func (u *roomUsecase) GetAll(ctx context.Context) ([]domain.Room, error) {
	return u.repo.FindAll(ctx)
}

func (u *roomUsecase) GetByID(ctx context.Context, id int64) (*domain.Room, error) {
	return u.repo.FindByID(ctx, id)
}

func (u *roomUsecase) Create(ctx context.Context, room *domain.Room) error {
	if err := validateRoom(room); err != nil {
		return err
	}
	return u.repo.Create(ctx, room)
}

func (u *roomUsecase) Update(ctx context.Context, room *domain.Room) error {
	if err := validateRoom(room); err != nil {
		return err
	}
	return u.repo.Update(ctx, room)
}

func (u *roomUsecase) Delete(ctx context.Context, id int64) error {
	return u.repo.Delete(ctx, id)
}

func validateRoom(room *domain.Room) error {
	if room.Code == "" || room.Name == "" || room.Capacity <= 0 {
		return domain.ErrInvalidRoomInput
	}
	return nil
}