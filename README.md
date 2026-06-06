# pim-api

Backend REST API untuk sistem kas Karang Taruna RT 22/06. Dibangun dengan Go + Gin + GORM + PostgreSQL.

## Tech Stack

- **Go** 1.26+
- **Gin** — HTTP framework
- **GORM** — ORM untuk PostgreSQL
- **Cloudinary** — penyimpanan foto bukti transfer & nota
- **JWT** — autentikasi admin via HTTP-only cookie
- **bcrypt** — hashing password

## Prasyarat

- Go 1.21+
- PostgreSQL 14+
- Akun Cloudinary

## Setup Database

Buat database PostgreSQL terlebih dahulu:

```bash
psql -U postgres -c "CREATE DATABASE pim;"
psql -U postgres -c "CREATE USER mongsky WITH PASSWORD '135504';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE pim TO mongsky;"
```

Jalankan schema + seed:

```bash
psql -U mongsky -d pim -f database/schema.sql
```

Ini akan membuat semua tabel, enum, index, dan akun superadmin default.

## Konfigurasi Environment

Buat file `.env` di root project:

```env
# Database
DATABASE_URL=postgres://<user>:<password>@localhost:5432/<dbname>?sslmode=disable

# JWT
JWT_SECRET=<random_string_panjang>

# Enkripsi nomor telepon
ENCRYPTION_KEY=<64_karakter_hex>

# Cloudinary
CLOUDINARY_CLOUD_NAME=<cloud_name>
CLOUDINARY_API_KEY=<api_key>
CLOUDINARY_API_SECRET=<api_secret>

# Server
PORT=3002
FRONTEND_URL=http://localhost:5173
NODE_ENV=development
```

Generate `ENCRYPTION_KEY` (64 karakter hex = 32 byte):

```bash
openssl rand -hex 32
```

## Menjalankan Server

```bash
# Install dependencies
go mod download

# Jalankan development
go run main.go

# Atau build dulu
go build -o pim-api main.go
./pim-api
```

Server berjalan di `http://localhost:3002`.

## Struktur Project

```
pim-api/
├── main.go
├── database/
│   └── schema.sql        # Schema + seed data
├── internal/
│   ├── config/           # Koneksi DB & env loader
│   ├── handler/          # HTTP handlers (thin layer)
│   ├── middleware/        # Auth JWT, rate limiter
│   ├── model/            # GORM models
│   ├── repository/       # Query DB
│   └── service/          # Business logic
└── pkg/
    ├── cloudinary/       # Upload helper
    ├── crypto/           # Enkripsi/dekripsi nomor telepon
    └── walink/           # Generator link WhatsApp notifikasi
```

## API Endpoints

### Public

| Method | Path | Deskripsi |
|--------|------|-----------|
| GET | `/api/public/dashboard` | Ringkasan kas (total masuk, keluar, saldo, chart) |
| GET | `/api/public/leaderboard` | Papan donatur terverifikasi |
| GET | `/api/public/contributions` | List donasi terverifikasi (publik) |
| POST | `/api/public/contributions` | Submit donasi baru |
| GET | `/api/public/expenses` | List pengeluaran |

### Admin (butuh login)

| Method | Path | Deskripsi |
|--------|------|-----------|
| POST | `/api/admin/auth/login` | Login |
| POST | `/api/admin/auth/logout` | Logout |
| GET | `/api/admin/auth/me` | Info user aktif |
| GET | `/api/admin/contributions` | List donasi (semua status) |
| PATCH | `/api/admin/contributions/:id/verify` | Verifikasi donasi |
| PATCH | `/api/admin/contributions/:id/reject` | Tolak donasi |
| GET | `/api/admin/contributions/export` | Export CSV |
| POST | `/api/admin/expenses` | Tambah pengeluaran |
| GET | `/api/admin/expenses` | List pengeluaran |
| GET | `/api/admin/users` | List admin (SUPER_ADMIN) |
| POST | `/api/admin/users` | Tambah admin (SUPER_ADMIN) |
| PATCH | `/api/admin/users/:id/role` | Ubah role (SUPER_ADMIN) |
| PATCH | `/api/admin/users/:id/password` | Reset password (SUPER_ADMIN) |
| DELETE | `/api/admin/users/:id` | Hapus admin (SUPER_ADMIN) |

## Akun Default

Setelah menjalankan `schema.sql`:

| Field | Value |
|-------|-------|
| Username | `superadmin` |
| Password | `@IRPK2206` |
| Role | `SUPER_ADMIN` |

**Ganti password setelah pertama kali login.**
