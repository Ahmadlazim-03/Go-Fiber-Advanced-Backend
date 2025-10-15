# Database Migrations

## Overview
File migrations telah dipisahkan berdasarkan tipe database untuk memudahkan maintenance dan pemahaman kode. Semua file migrations berada di folder `database/migration/`.

## Folder Structure

```
database/
├── connection.go          # Koneksi database PostgreSQL, MongoDB & PocketBase
└── migration/
    ├── migrations.go              # Main wrapper
    ├── migrations_postgres.go     # PostgreSQL migrations
    ├── migrations_mongodb.go      # MongoDB migrations
    └── migrations_pocketbase.go   # PocketBase migrations
```

## File Structure

### 1. `migration/migrations.go` (Main Wrapper)
File wrapper utama yang menentukan migrations mana yang akan dijalankan berdasarkan tipe database yang digunakan.

```go
func RunMigrations()
```
- Membaca tipe database dari environment variable `DB_TYPE`
- Memanggil `RunPostgresMigrations()` jika menggunakan PostgreSQL
- Memanggil `RunMongoDBMigrations()` jika menggunakan MongoDB
- Memanggil `RunPocketBaseMigrations()` jika menggunakan PocketBase

### 2. `migration/migrations_postgres.go` (PostgreSQL Migrations)
Berisi semua migrations untuk PostgreSQL menggunakan GORM Auto Migrate.

**Functions:**
- `RunPostgresMigrations()` - Membuat tabel-tabel PostgreSQL
- `createPostgresIndexes()` - Membuat indexes untuk optimasi query

**Tables Created:**
- `users` - Tabel user untuk autentikasi
- `mahasiswas` - Tabel data mahasiswa
- `alumnis` - Tabel data alumni
- `pekerjaan_alumnis` - Tabel data pekerjaan alumni

**Indexes Created:**
- `idx_users_email` - Index pada users.email (unique)
- `idx_users_username` - Index pada users.username (unique)
- `idx_mahasiswas_nim` - Index pada mahasiswas.nim (unique)
- `idx_alumnis_nim` - Index pada alumnis.nim (unique)
- `idx_alumnis_user_id` - Index pada alumnis.user_id (foreign key)
- `idx_pekerjaan_alumni_id` - Index pada pekerjaan_alumnis.alumni_id (foreign key)
- `idx_pekerjaan_deleted_at` - Index pada pekerjaan_alumnis.deleted_at (soft delete)

### 3. `migrations_mongodb.go` (MongoDB Migrations)
Berisi semua migrations untuk MongoDB menggunakan MongoDB Driver.

**Functions:**
- `RunMongoDBMigrations()` - Membuat collections dan indexes MongoDB
- `createMongoDBIndexes()` - Membuat indexes untuk optimasi query
- `createMongoIndex()` - Helper function untuk membuat single index
- `DropMongoDBCollections()` - Menghapus semua collections (untuk development)

**Collections Created:**
- `users` - Collection user untuk autentikasi
- `mahasiswas` - Collection data mahasiswa
- `alumnis` - Collection data alumni
- `pekerjaan_alumnis` - Collection data pekerjaan alumni

**Indexes Created:**
- `idx_users_email` - Index pada users.email (unique)
- `idx_users_username` - Index pada users.username (unique)
- `idx_mahasiswas_nim` - Index pada mahasiswas.nim (unique)
- `idx_mahasiswas_email` - Index pada mahasiswas.email (unique)
- `idx_alumnis_nim` - Index pada alumnis.nim (unique)
- `idx_alumnis_user_id` - Index pada alumnis.user_id
- `idx_pekerjaan_alumni_id` - Index pada pekerjaan_alumnis.alumni_id
- `idx_pekerjaan_deleted_at` - Index pada pekerjaan_alumnis.deleted_at

## Usage

Migrations akan otomatis dijalankan saat aplikasi start di `main.go`:

```go
import (
    "modul4crud/database"
    "modul4crud/database/migration"
)

func main() {
    // Connect to database
    database.ConnectDB()

    // Run migrations
    migration.RunMigrations()
}
```

## Environment Variables

Pastikan `.env` file sudah dikonfigurasi dengan benar:

### Untuk PostgreSQL:
```env
DB_TYPE=postgres
POSTGRES_DSN=postgresql://user:password@host:port/database
```

### Untuk MongoDB:
```env
DB_TYPE=mongodb
MONGODB_URI=mongodb://username:password@host:port
MONGODB_DATABASE=database_name
```

## Features

### PostgreSQL Migrations
- ✅ Auto-create tables with GORM
- ✅ Check if table exists before creating
- ✅ Create indexes automatically
- ✅ Support foreign key relationships
- ✅ Support soft deletes with DeletedAt

### MongoDB Migrations
- ✅ Auto-create collections
- ✅ Check if collection exists before creating
- ✅ Create indexes with unique constraints
- ✅ Support compound indexes
- ✅ Helper function untuk drop collections

## Notes

1. **PostgreSQL** menggunakan GORM Migrator yang akan:
   - Membuat tabel otomatis berdasarkan struct models
   - Menambahkan foreign keys
   - Membuat indexes yang diperlukan

2. **MongoDB** menggunakan Native MongoDB Driver yang akan:
   - Membuat collections secara eksplisit
   - Membuat indexes secara manual
   - Tidak perlu schema definition (schemaless)

3. **Automatic Execution**: Migrations akan berjalan otomatis setiap kali aplikasi start, tetapi akan skip jika tabel/collection sudah ada.

4. **Idempotent**: Migrations bersifat idempotent, artinya aman untuk dijalankan berkali-kali.

## Troubleshooting

### PostgreSQL
- Pastikan POSTGRES_DSN format sudah benar
- Check koneksi database dengan `psql` atau pgAdmin
- Verifikasi user memiliki permission untuk CREATE TABLE

### MongoDB
- Pastikan MONGODB_URI format sudah benar
- Check koneksi dengan MongoDB Compass atau mongosh
- Verifikasi user memiliki permission untuk createCollection dan createIndex

## Development Tips

### Menambah Tabel/Collection Baru

**PostgreSQL:**
1. Tambahkan model baru di folder `models/`
2. Tambahkan migration di `database/migration/migrations_postgres.go` → `RunPostgresMigrations()`
3. Tambahkan index jika diperlukan di `createPostgresIndexes()`

**MongoDB:**
1. Tambahkan collection name di array `collections` di `database/migration/migrations_mongodb.go`
2. Tambahkan index di `createMongoDBIndexes()`
3. Implementasikan repository di `repositories/mongodb/`

### Reset Database (Development Only)

**PostgreSQL:**
```sql
DROP DATABASE yourdb;
CREATE DATABASE yourdb;
```

**MongoDB:**
```go
// Call this function (use with caution!)
// Di file database/migration/migrations_mongodb.go
migration.DropMongoDBCollections()
```
