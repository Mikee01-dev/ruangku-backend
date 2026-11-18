package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Mikee01-dev/ruangku-backend/internal/domain"
)

type roomUsecaseStub struct {
	rooms       []domain.Room
	room        *domain.Room
	err         error
	created     *domain.Room
	updated     *domain.Room
	deletedID   int64
	requestedID int64
}

func (s *roomUsecaseStub) GetAll(context.Context) ([]domain.Room, error) {
	return s.rooms, s.err
}

func (s *roomUsecaseStub) GetByID(_ context.Context, id int64) (*domain.Room, error) {
	s.requestedID = id
	return s.room, s.err
}

func (s *roomUsecaseStub) Create(_ context.Context, room *domain.Room) error {
	s.created = room
	return s.err
}

func (s *roomUsecaseStub) Update(_ context.Context, room *domain.Room) error {
	s.updated = room
	return s.err
}

func (s *roomUsecaseStub) Delete(_ context.Context, id int64) error {
	s.deletedID = id
	return s.err
}

func newRoomMux(uc domain.RoomUsecase) *http.ServeMux {
	mux := http.NewServeMux()
	NewRoomDelivery(uc).RegisterRoutes(mux)
	return mux
}

func TestRoomDeliveryRoutes(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		usecaseErr error
		wantStatus int
		check      func(*testing.T, *roomUsecaseStub)
	}{
		{
			name:       "get all rooms",
			method:     http.MethodGet,
			path:       "/rooms",
			wantStatus: http.StatusOK,
		},
		{
			name:       "create room",
			method:     http.MethodPost,
			path:       "/rooms",
			body:       `{"code":"A101","name":"Ruang A","capacity":30}`,
			wantStatus: http.StatusCreated,
			check: func(t *testing.T, uc *roomUsecaseStub) {
				if uc.created == nil || uc.created.Code != "A101" {
					t.Fatalf("Create was not called with the decoded room: %#v", uc.created)
				}
			},
		},
		{
			name:       "reject malformed room body",
			method:     http.MethodPost,
			path:       "/rooms",
			body:       `{"code":"A101"} {"code":"A102"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "reject invalid room id",
			method:     http.MethodGet,
			path:       "/rooms/not-an-id",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "map missing room to not found",
			method:     http.MethodGet,
			path:       "/rooms/12",
			usecaseErr: domain.ErrRoomNotFound,
			wantStatus: http.StatusNotFound,
			check: func(t *testing.T, uc *roomUsecaseStub) {
				if uc.requestedID != 12 {
					t.Fatalf("GetByID received id %d, want 12", uc.requestedID)
				}
			},
		},
		{
			name:       "map invalid room input to bad request",
			method:     http.MethodPost,
			path:       "/rooms",
			body:       `{"code":"","name":"","capacity":0}`,
			usecaseErr: domain.ErrInvalidRoomInput,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "map unknown usecase error to internal server error",
			method:     http.MethodGet,
			path:       "/rooms",
			usecaseErr: errors.New("database unavailable"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "use path id for update",
			method:     http.MethodPut,
			path:       "/rooms/15",
			body:       `{"id":99,"code":"A101","name":"Ruang A","capacity":30}`,
			wantStatus: http.StatusOK,
			check: func(t *testing.T, uc *roomUsecaseStub) {
				if uc.updated == nil || uc.updated.ID != 15 {
					t.Fatalf("Update received room %#v, want path id 15", uc.updated)
				}
			},
		},
		{
			name:       "delete room",
			method:     http.MethodDelete,
			path:       "/rooms/20",
			wantStatus: http.StatusNoContent,
			check: func(t *testing.T, uc *roomUsecaseStub) {
				if uc.deletedID != 20 {
					t.Fatalf("Delete received id %d, want 20", uc.deletedID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &roomUsecaseStub{err: tt.usecaseErr}
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			newRoomMux(uc).ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body: %s", rec.Code, tt.wantStatus, rec.Body.String())
			}
			if tt.check != nil {
				tt.check(t, uc)
			}
		})
	}
}
