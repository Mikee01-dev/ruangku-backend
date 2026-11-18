package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Mikee01-dev/ruangku-backend/internal/domain"
)

type roomRepository struct {
	db *sql.DB
}

// NewRoomRepository membuat implementasi domain.RoomRepository berbasis PostgreSQL.
func NewRoomRepository(db *sql.DB) domain.RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) FindAll(ctx context.Context) ([]domain.Room, error) {
	const query = `
		SELECT id, code, name, building, capacity, is_active, created_at
		FROM rooms
		ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := make([]domain.Room, 0)
	for rows.Next() {
		var room domain.Room
		if err := rows.Scan(
			&room.ID,
			&room.Code,
			&room.Name,
			&room.Building,
			&room.Capacity,
			&room.IsActive,
			&room.CreatedAt,
		); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (r *roomRepository) FindByID(ctx context.Context, id int64) (*domain.Room, error) {
	const query = `
		SELECT id, code, name, building, capacity, is_active, created_at
		FROM rooms
		WHERE id = $1`

	var room domain.Room
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&room.ID,
		&room.Code,
		&room.Name,
		&room.Building,
		&room.Capacity,
		&room.IsActive,
		&room.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrRoomNotFound
		}
		return nil, err
	}

	return &room, nil
}

func (r *roomRepository) Create(ctx context.Context, room *domain.Room) error {
	const query = `
		INSERT INTO rooms (code, name, building, capacity, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	return r.db.QueryRowContext(
		ctx, query,
		room.Code, room.Name, room.Building, room.Capacity, room.IsActive,
	).Scan(&room.ID, &room.CreatedAt)
}

func (r *roomRepository) Update(ctx context.Context, room *domain.Room) error {
	const query = `
		UPDATE rooms
		SET code = $1, name = $2, building = $3, capacity = $4, is_active = $5
		WHERE id = $6`

	res, err := r.db.ExecContext(
		ctx, query,
		room.Code, room.Name, room.Building, room.Capacity, room.IsActive, room.ID,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrRoomNotFound
	}

	return nil
}

func (r *roomRepository) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM rooms WHERE id = $1`, id)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return domain.ErrRoomNotFound
	}

	return nil
}