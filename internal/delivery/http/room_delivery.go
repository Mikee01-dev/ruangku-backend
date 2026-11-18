package http

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Mikee01-dev/ruangku-backend/internal/domain"
)

type RoomDelivery struct {
	usecase domain.RoomUsecase
}

func NewRoomDelivery(uc domain.RoomUsecase) *RoomDelivery {
	return &RoomDelivery{usecase: uc}
}

func (d *RoomDelivery) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /rooms", d.getAll)
	mux.HandleFunc("GET /rooms/{id}", d.getByID)
	mux.HandleFunc("POST /rooms", d.create)
	mux.HandleFunc("PUT /rooms/{id}", d.update)
	mux.HandleFunc("DELETE /rooms/{id}", d.delete)
}

func (d *RoomDelivery) getAll(w http.ResponseWriter, r *http.Request) {
	rooms, err := d.usecase.GetAll(r.Context())
	if err != nil {
		d.writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, rooms)
}

func (d *RoomDelivery) getByID(w http.ResponseWriter, r *http.Request) {
	id, ok := roomID(w, r)
	if !ok {
		return
	}

	room, err := d.usecase.GetByID(r.Context(), id)
	if err != nil {
		d.writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, room)
}

func (d *RoomDelivery) create(w http.ResponseWriter, r *http.Request) {
	var room domain.Room
	if err := decodeRoom(r, &room); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := d.usecase.Create(r.Context(), &room); err != nil {
		d.writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, room)
}

func (d *RoomDelivery) update(w http.ResponseWriter, r *http.Request) {
	id, ok := roomID(w, r)
	if !ok {
		return
	}

	var room domain.Room
	if err := decodeRoom(r, &room); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	room.ID = id

	if err := d.usecase.Update(r.Context(), &room); err != nil {
		d.writeUsecaseError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, room)
}

func (d *RoomDelivery) delete(w http.ResponseWriter, r *http.Request) {
	id, ok := roomID(w, r)
	if !ok {
		return
	}

	if err := d.usecase.Delete(r.Context(), id); err != nil {
		d.writeUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func roomID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid room id")
		return 0, false
	}

	return id, true
}

func decodeRoom(r *http.Request, room *domain.Room) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(room); err != nil {
		return err
	}

	var extra json.RawMessage
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON object")
	}

	return nil
}

func (d *RoomDelivery) writeUsecaseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidRoomInput):
		writeError(w, http.StatusBadRequest, domain.ErrInvalidRoomInput.Error())
	case errors.Is(err, domain.ErrRoomNotFound):
		writeError(w, http.StatusNotFound, domain.ErrRoomNotFound.Error())
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	body, err := json.Marshal(value)
	if err != nil {
		log.Printf("room delivery: failed to encode JSON response: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(append(body, '\n')); err != nil {
		log.Printf("room delivery: failed to write JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, struct {
		Error string `json:"error"`
	}{Error: message})
}
