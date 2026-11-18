-- RuangKu - Migrasi awal
-- Menjalankan: eksekusi file ini di PostgreSQL

CREATE EXTENSION IF NOT EXISTS btree_gist;

-- 1. USERS
CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(100) NOT NULL,
    email         VARCHAR(150) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20)  NOT NULL DEFAULT 'peminjam'
                  CHECK (role IN ('admin', 'pengelola', 'peminjam')),
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- 2. ROOMS
CREATE TABLE IF NOT EXISTS rooms (
    id         SERIAL PRIMARY KEY,
    code       VARCHAR(20)  NOT NULL UNIQUE,
    name       VARCHAR(100) NOT NULL,
    building   VARCHAR(100) NOT NULL,
    capacity   INT          NOT NULL CHECK (capacity > 0),
    is_active  BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

-- 3. BOOKINGS
CREATE TABLE IF NOT EXISTS bookings (
    id         SERIAL PRIMARY KEY,
    user_id    INT  NOT NULL REFERENCES users(id),
    room_id    INT  NOT NULL REFERENCES rooms(id),
    purpose    TEXT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time   TIMESTAMP NOT NULL,
    status     VARCHAR(20) NOT NULL DEFAULT 'pending'
               CHECK (status IN ('pending', 'approved', 'rejected', 'cancelled')),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CHECK (end_time > start_time),
    -- Cegah dua peminjaman aktif di ruang dan jam yang bentrok
    CONSTRAINT no_overlap_booking EXCLUDE USING gist (
        room_id WITH =,
        tsrange(start_time, end_time) WITH &&
    ) WHERE (status IN ('pending', 'approved'))
);

-- 4. BOOKING_LOGS
CREATE TABLE IF NOT EXISTS booking_logs (
    id          SERIAL PRIMARY KEY,
    booking_id  INT NOT NULL REFERENCES bookings(id) ON DELETE CASCADE,
    actor_id    INT NOT NULL REFERENCES users(id),
    from_status VARCHAR(20),
    to_status   VARCHAR(20) NOT NULL,
    note        TEXT,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Index untuk mempercepat query
CREATE INDEX IF NOT EXISTS idx_bookings_room_time ON bookings (room_id, start_time, end_time);
CREATE INDEX IF NOT EXISTS idx_bookings_user ON bookings (user_id);
CREATE INDEX IF NOT EXISTS idx_booking_logs_booking ON booking_logs (booking_id);