# Setup Database Fleksibel (PostgreSQL & MongoDB)

## Deskripsi
Proyek ini sudah dikonfigurasi untuk mendukung 2 jenis database secara fleksibel:
- **PostgreSQL** - Database relational
- **MongoDB** - Database NoSQL

Anda dapat dengan mudah beralih antara PostgreSQL dan MongoDB hanya dengan mengubah konfigurasi di file `.env`.

## Struktur Repository

```
repositories/
├── interfaces.go                    # Interface untuk semua repository
├── postgre/                         # Repository untuk PostgreSQL
│   ├── user_repository.go
│   ├── mahasiswa_repository.go
│   ├── alumni_repository.go
│   └── pekerjaan_alumni_repository.go
└── mongodb/                         # Repository untuk MongoDB
    ├── user_repository_mongo.go
    ├── mahasiswa_repository_mongo.go
    ├── alumni_repository_mongo.go
    └── pekerjaan_alumni_repository_mongo.go
```

## Konfigurasi Database

### 1. Setup File .env

Copy file `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```

Edit file `.env` sesuai kebutuhan:

```env
# Pilih tipe database: postgres atau mongodb
DB_TYPE=postgres

# Konfigurasi PostgreSQL
POSTGRES_DSN=postgresql://postgres:erjghiShjYBhXeOiQuQXrrGEabuCqxxP@switchyard.proxy.rlwy.net:54521/railway

# Konfigurasi MongoDB
MONGODB_URI=mongodb://mongo:pakgtnLdkcJlREVyWpuhiecIEQvnVOkh@caboose.proxy.rlwy.net:48828
MONGODB_DATABASE=railway

# Konfigurasi Server
SERVER_PORT=8080

# Konfigurasi JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=24h
```

### 2. Menggunakan PostgreSQL

Ubah variable `DB_TYPE` di `.env`:
```env
DB_TYPE=postgres
POSTGRES_DSN=postgresql://username:password@host:port/database
```

Format DSN PostgreSQL:
- Railway: `postgresql://user:password@host:port/database`
- Local: `host=localhost user=postgres password=yourpassword dbname=yourdb port=5432 sslmode=disable`

### 3. Menggunakan MongoDB

Ubah variable `DB_TYPE` di `.env`:
```env
DB_TYPE=mongodb
MONGODB_URI=mongodb://username:password@host:port
MONGODB_DATABASE=nama_database
```

Format URI MongoDB:
- Railway: `mongodb://user:password@host:port`
- Local: `mongodb://localhost:27017`
- Atlas: `mongodb+srv://user:password@cluster.mongodb.net`

## Cara Menjalankan

### 1. Install Dependencies

```bash
go mod tidy
```

### 2. Build Aplikasi

```bash
go build -o app
```

### 3. Jalankan Aplikasi

```bash
./app
```

Atau langsung run:
```bash
go run main.go
```

## Fitur Database Fleksibel

### Automatic Repository Selection
Aplikasi secara otomatis memilih repository yang sesuai berdasarkan `DB_TYPE` di `.env`:

```go
if database.IsPostgres() {
    userRepo = postgre.NewUserRepository(database.DB)
    mahasiswaRepo = postgre.NewMahasiswaRepository(database.DB)
    alumniRepo = postgre.NewAlumniRepository(database.DB)
    pekerjaanRepo = postgre.NewPekerjaanAlumniRepository(database.DB)
} else if database.IsMongoDB() {
    userRepo = mongodb.NewUserRepositoryMongo(database.MongoDB)
    mahasiswaRepo = mongodb.NewMahasiswaRepositoryMongo(database.MongoDB)
    alumniRepo = mongodb.NewAlumniRepositoryMongo(database.MongoDB)
    pekerjaanRepo = mongodb.NewPekerjaanAlumniRepositoryMongo(database.MongoDB)
}
```

### Interface-Based Design
Semua repository mengimplementasikan interface yang sama, memastikan kompatibilitas API:

```go
type UserRepository interface {
    GetAll() ([]models.User, error)
    GetWithPagination(pagination *models.PaginationRequest) ([]models.User, int64, error)
    GetByID(id int) (*models.User, error)
    GetByEmail(email string) (*models.User, error)
    GetByUsername(username string) (*models.User, error)
    Create(user *models.User) error
    Update(user *models.User) error
    Delete(id int) error
    Count() (int64, error)
}
```

## Migration

### PostgreSQL
Migration otomatis dijalankan saat aplikasi start jika menggunakan PostgreSQL:
- Membuat tabel jika belum ada
- Membuat index untuk performa
- Membuat user admin default

### MongoDB
Untuk MongoDB, collections akan dibuat otomatis saat data pertama kali diinsert.

## Default Admin User

Aplikasi akan otomatis membuat user admin jika belum ada:
- **Email**: admin@example.com
- **Password**: admin123
- **Role**: admin

## Troubleshooting

### Error: Database connection failed
- Pastikan koneksi database sudah benar
- Cek firewall dan network connectivity
- Untuk Railway, pastikan IP Anda sudah diizinkan

### Error: Package not found
```bash
go mod tidy
go get -u all
```

### Ganti Database
Cukup ubah `DB_TYPE` di `.env` dan restart aplikasi.

## Testing

### Test dengan PostgreSQL
```bash
# Edit .env
DB_TYPE=postgres

# Run aplikasi
go run main.go
```

### Test dengan MongoDB
```bash
# Edit .env
DB_TYPE=mongodb

# Run aplikasi
go run main.go
```

## Dependencies

- **Fiber v2** - Web framework
- **GORM** - ORM untuk PostgreSQL
- **MongoDB Driver** - Driver resmi MongoDB untuk Go
- **godotenv** - Load environment variables
- **JWT** - Authentication
- **bcrypt** - Password hashing

## Koneksi Database yang Dikonfigurasi

### PostgreSQL (Railway)
```
Host: switchyard.proxy.rlwy.net
Port: 54521
User: postgres
Password: erjghiShjYBhXeOiQuQXrrGEabuCqxxP
Database: railway
```

### MongoDB (Railway)
```
Host: caboose.proxy.rlwy.net
Port: 48828
User: mongo
Password: pakgtnLdkcJlREVyWpuhiecIEQvnVOkh
Database: railway
```

## API Endpoints

Semua endpoint REST API tetap sama, tidak terpengaruh oleh jenis database yang digunakan:

- `POST /api/auth/register` - Register user baru
- `POST /api/auth/login` - Login user
- `GET /api/mahasiswas` - Get all mahasiswa
- `GET /api/alumnis` - Get all alumni
- `GET /api/pekerjaan-alumnis` - Get all pekerjaan alumni
- Dan lainnya...

Lihat `API_DOCUMENTATION.md` untuk detail lengkap.

## Support

Jika ada pertanyaan atau masalah, silakan buat issue di repository ini.
