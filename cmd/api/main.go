package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	roomHttp "github.com/Mikee01-dev/ruangku-backend/internal/delivery/http"
	
	"github.com/Mikee01-dev/ruangku-backend/internal/repository"
	"github.com/Mikee01-dev/ruangku-backend/internal/usecase"
	"github.com/Mikee01-dev/ruangku-backend/pkg/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan")
	}

	db, err := database.NewPostgres()
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}
	defer db.Close()

	roomRepo := repository.NewRoomRepository(db)
	roomUsecase := usecase.NewRoomUsecase(roomRepo)
	
	roomDelivery := roomHttp.NewRoomDelivery(roomUsecase)

	mux := http.NewServeMux()
	roomDelivery.RegisterRoutes(mux)

	port := ":8080"
	fmt.Printf("Server RuangKu berjalan di http://localhost%s\n", port)
	
	err = http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
}