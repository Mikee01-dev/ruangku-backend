# RuangKu - Sistem Peminjaman Ruang Kampus

## Deskripsi Sistem
RuangKu adalah sistem berbasis web yang membantu mahasiswa dan organisasi kemahasiswaan meminjam ruang kelas, lab, atau aula secara online. Peminjam mengajukan peminjaman, lalu pengelola ruang menyetujui atau menolaknya. Sistem mencegah dua peminjaman pada ruang dan jam yang sama.

## Role pada Sistem
* **Admin**: Mengelola user dan data ruang.
* **Pengelola Ruang**: Menyetujui atau menolak peminjaman.
* **Peminjam**: Melihat ruang, mengajukan & membatalkan peminjaman, melihat riwayat.

## Fitur Utama
1. Registrasi dan login (Semua)
2. Kelola data ruang / CRUD (Admin)
3. Cek ketersediaan ruang per tanggal (Peminjam)
4. Ajukan peminjaman (Peminjam)
5. Setujui / tolak peminjaman (Pengelola Ruang)
6. Batalkan peminjaman (Peminjam)
7. Riwayat peminjaman (Peminjam, Pengelola Ruang)

## Rancangan Database
* **users**: id, name, email, password_hash, role, created_at
* **rooms**: id, code, name, building, capacity, is_active, created_at
* **bookings**: id, user_id, room_id, purpose, start_time, end_time, status, created_at
* **booking_logs**: id, booking_id, actor_id, from_status, to_status, note, created_at
*(Status peminjaman: pending, approved, rejected, cancelled)*

## Ringkasan Pembagian Tugas

| Orang | Nama | Branch | File Utama |
| :--- | :--- | :--- | :--- |
| **A (Git Master)** | **Michael** | `main`, `feat/usecase-room`, `feat/wiring` | `room_usecase.go`, `cmd/api/main.go`, `README.md` |
| **B** | **Sovia** | `feat/migration`, `feat/repository-room`, `feat/database-pkg`, `feat/delivery-room` | `001_init.sql`, `room_repository.go`, `postgres.go`, `room_delivery.go` |
| **C** | **Rafael** | `feat/domain` | `room.go`, `user.go`, `booking.go` |

## Cara Menjalankan Proyek
1. Clone repository ini:
   ```bash
   git clone https://github.com/Mikee01-dev/ruangku-backend.git
   ```
2. Masuk ke folder proyek:
   ```bash
   cd ruangku-backend
   ```
3. Salin file environment:
   ```bash
   cp .env.example .env
   ```
4. Sesuaikan konfigurasi database (user, password, nama db) di dalam file `.env` dengan yang ada di laptop Anda.
5. Jalankan migrasi database (eksekusi file `migrations/001_init.sql` di PostgreSQL).
6. Jalankan aplikasi:
   ```bash
   go run cmd/api/main.go
   ```
7. Aplikasi akan berjalan di `http://localhost:8080`.