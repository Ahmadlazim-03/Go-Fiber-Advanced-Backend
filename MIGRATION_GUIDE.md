# Database Migration System

Sistem migrasi database ini dirancang untuk membuat tabel secara otomatis hanya jika tabel tersebut belum ada di database.

## 📁 File Struktur

```
database/
├── connection.go      # Koneksi database PostgreSQL
└── migrations.go      # Sistem migrasi tabel dan index
```

## 🔧 Fitur Migrasi

### 1. **Smart Table Creation**
- Cek keberadaan tabel sebelum membuat
- Hanya membuat tabel yang belum ada
- Tidak menghapus data yang sudah ada

### 2. **Automatic Index Creation**
- Index untuk performa query
- Index pada kolom yang sering digunakan untuk search
- Index untuk foreign key relationships

### 3. **Health Check**
- Validasi koneksi database
- Ping database sebelum migrasi

## 📋 Tabel yang Dibuat

### 1. **users**
- Primary key: `id`
- Unique fields: `username`, `email`
- Indexes: `email`, `username`

### 2. **mahasiswas** 
- Primary key: `id`
- Unique field: `nim`
- Indexes: `nim`

### 3. **alumnis**
- Primary key: `id`
- Foreign key: `user_id` -> `users.id`
- Indexes: `nim`, `user_id`

### 4. **pekerjaan_alumnis**
- Primary key: `id`
- Foreign key: `alumni_id` -> `alumnis.id`
- Soft delete support: `deleted_at`
- Indexes: `alumni_id`, `deleted_at`

## 🚀 Cara Kerja

### Startup Sequence:
1. **Connect Database** - Buat koneksi ke PostgreSQL
2. **Health Check** - Pastikan database dapat diakses
3. **Run Migrations** - Buat tabel dan index jika belum ada
4. **Create Admin User** - Buat user admin default jika belum ada

### Migration Process:
```go
// 1. Cek apakah tabel ada
if !DB.Migrator().HasTable(&models.User{}) {
    // 2. Buat tabel jika belum ada
    DB.Migrator().CreateTable(&models.User{})
}

// 3. Buat index untuk performa
if !DB.Migrator().HasIndex(&models.User{}, "idx_users_email") {
    DB.Migrator().CreateIndex(&models.User{}, "email")
}
```

## ✅ Keuntungan Sistem Ini

### **Safe Migration**
- ✅ Tidak menghapus data existing
- ✅ Tidak mengubah struktur tabel yang sudah ada
- ✅ Idempotent - aman dijalankan berulang kali

### **Performance Optimized**
- ✅ Index otomatis untuk query yang sering digunakan
- ✅ Foreign key constraints untuk data integrity
- ✅ Soft delete support untuk pekerjaan alumni

### **Development Friendly**
- ✅ Logging yang jelas untuk setiap step
- ✅ Error handling yang proper
- ✅ No manual SQL script needed

## 📝 Log Output Example

```
Running database migrations...
✓ Users table already exists
Creating mahasiswas table...
✓ Mahasiswas table created successfully
✓ Alumnis table already exists
✓ Pekerjaan_alumnis table already exists
Creating database indexes...
✓ Created index on users.email
✓ Created index on mahasiswas.nim
Database migrations completed successfully!
✓ Database connection is healthy
Checking for default admin user...
✓ Admin user already exists
```

## 🔧 Configuration

Database connection dalam `database/connection.go`:
```go
dsn := "host=localhost user=postgres password= dbname=postgree port=5432 sslmode=disable"
```

## 🛠️ Manual Migration (Jika Diperlukan)

Jika perlu menjalankan migrasi manual:
```go
database.ConnectDB()
database.RunMigrations()
```

## ⚠️ Important Notes

1. **Backup Data**: Selalu backup data sebelum upgrade aplikasi
2. **Database Permissions**: User database harus memiliki permission CREATE TABLE dan CREATE INDEX
3. **Version Control**: Perubahan schema harus melalui kode, bukan manual SQL
4. **Testing**: Test migrasi di development environment dulu