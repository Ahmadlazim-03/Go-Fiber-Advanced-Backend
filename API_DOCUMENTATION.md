# 📚 API Documentation - CRUD Go Fiber PostgreSQL

## 🔗 Base URL
```
http://localhost:8080
```

---

## 🔐 Authentication Endpoints

### 1. Register User
**Endpoint:** `POST /api/register`  
**Deskripsi:** Mendaftarkan user baru ke sistem

**Request:**
```json
{
  "username": "john_doe",
  "email": "john@example.com", 
  "password": "password123",
  "role": "user"
}
```

**Response:**
```json
{
  "message": "User berhasil didaftarkan",
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2025-09-25T12:00:00Z"
  }
}
```

### 2. Login User
**Endpoint:** `POST /api/login`  
**Deskripsi:** Login untuk mendapatkan JWT token

**Request:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "role": "user"
    }
  }
}
```

---

## 👤 User CRUD Operations

**Headers Required:** `Authorization: Bearer <jwt_token>`

### 3. Get All Users (with Pagination & Search)
**Endpoint:** `GET /api/users?page=1&per_page=10&search=admin&sort_by=username&sort_order=asc`  
**Deskripsi:** Mendapatkan semua data user dengan pagination dan pencarian (Admin only)

**Query Parameters:**
- `page`: Halaman (default: 1)
- `per_page`: Data per halaman (default: 10)
- `search`: Kata kunci pencarian (username, email, role)
- `sort_by`: Kolom sorting (default: id)
- `sort_order`: Urutan sorting - ASC/DESC atau asc/desc (default: ASC)

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "username": "admin", 
      "email": "admin@example.com",
      "role": "admin",
      "is_active": true,
      "created_at": "2025-09-25T12:00:00Z"
    }
  ],
  "current_page": 1,
  "per_page": 10,
  "total_data": 1,
  "total_pages": 1,
  "has_next": false,
  "has_previous": false
}
```

### 4. Get User by ID
**Endpoint:** `GET /api/users/{id}`  
**Deskripsi:** Mendapatkan detail user berdasarkan ID

**Response:**
```json
{
  "data": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com", 
    "role": "user",
    "is_active": true,
    "created_at": "2025-09-25T12:00:00Z"
  }
}
```

### 5. Update User
**Endpoint:** `PUT /api/users/{id}`  
**Deskripsi:** Update data user

**Request:**
```json
{
  "username": "john_doe_updated",
  "email": "john.updated@example.com",
  "role": "admin",
  "is_active": true
}
```

### 6. Delete User
**Endpoint:** `DELETE /api/users/{id}`  
**Deskripsi:** Menghapus user (Hard delete)

**Response:**
```json
{
  "message": "User berhasil dihapus"
}
```

---

## 🎓 Mahasiswa CRUD Operations

**Headers Required:** `Authorization: Bearer <jwt_token>`

### 7. Get All Mahasiswa (with Pagination & Search)
**Endpoint:** `GET /api/mahasiswa?page=1&per_page=10&search=teknik&sort_by=nama&sort_order=asc`  
**Deskripsi:** Mendapatkan data mahasiswa dengan pagination dan pencarian

**Query Parameters:**
- `page`: Halaman (default: 1) 
- `per_page`: Data per halaman (default: 10)
- `search`: Kata kunci pencarian (nim, nama, email, jurusan, angkatan)
- `sort_by`: Kolom sorting (default: id)
- `sort_order`: Urutan sorting - ASC/DESC atau asc/desc (default: ASC)

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "nim": "2021001",
      "nama": "John Doe",
      "email": "john.student@example.com",
      "jurusan": "Teknik Informatika",
      "angkatan": 2021,
      "created_at": "2025-09-25T12:00:00Z"
    }
  ],
  "current_page": 1,
  "per_page": 10,
  "total_data": 1,
  "total_pages": 1,
  "has_next": false,
  "has_previous": false
}
```

### 8. Get Mahasiswa by ID
**Endpoint:** `GET /api/mahasiswa/{id}`  
**Deskripsi:** Mendapatkan detail mahasiswa berdasarkan ID

**Response:**
```json
{
  "data": {
    "id": 1,
    "nim": "2021001", 
    "nama": "John Doe",
    "email": "john.student@example.com",
    "jurusan": "Teknik Informatika",
    "angkatan": 2021,
    "created_at": "2025-09-25T12:00:00Z"
  }
}
```

### 9. Create Mahasiswa
**Endpoint:** `POST /api/mahasiswa`  
**Deskripsi:** Membuat data mahasiswa baru

**Request:**
```json
{
  "nim": "2021001",
  "nama": "John Doe",
  "email": "john.student@example.com", 
  "jurusan": "Teknik Informatika",
  "angkatan": 2021
}
```

### 10. Update Mahasiswa
**Endpoint:** `PUT /api/mahasiswa/{id}`  
**Deskripsi:** Update data mahasiswa

**Request:**
```json
{
  "nim": "2021001",
  "nama": "John Doe Updated",
  "email": "john.updated@example.com",
  "jurusan": "Sistem Informasi", 
  "angkatan": 2021
}
```

### 11. Delete Mahasiswa
**Endpoint:** `DELETE /api/mahasiswa/{id}`  
**Deskripsi:** Menghapus mahasiswa (Hard delete)

**Response:**
```json
{
  "message": "Mahasiswa berhasil dihapus"
}
```

---

## 🎯 Alumni CRUD Operations

**Headers Required:** `Authorization: Bearer <jwt_token>`

### 12. Get All Alumni (with Pagination & Search)
**Endpoint:** `GET /api/alumni?page=1&per_page=10&search=teknik&sort_by=nama&sort_order=desc`  
**Deskripsi:** Mendapatkan data alumni dengan relasi user dan pagination

**Query Parameters:**
- `page`: Halaman (default: 1)
- `per_page`: Data per halaman (default: 10) 
- `search`: Kata kunci pencarian (nim, nama, jurusan, tahun_lulus, user.email)
- `sort_by`: Kolom sorting (default: id)
- `sort_order`: Urutan sorting - ASC/DESC atau asc/desc (default: ASC)

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "user_id": 2,
      "user": {
        "id": 2,
        "username": "jane_alumni",
        "email": "jane@example.com",
        "role": "user"
      },
      "nim": "2018001",
      "nama": "Jane Smith",
      "jurusan": "Teknik Informatika", 
      "angkatan": 2018,
      "tahun_lulus": 2022,
      "no_telepon": "08123456789",
      "alamat": "Jakarta Selatan",
      "created_at": "2025-09-25T12:00:00Z",
      "pekerjaan_alumni": []
    }
  ],
  "current_page": 1,
  "per_page": 10,
  "total_data": 1,
  "total_pages": 1
}
```

### 13. Get Alumni by ID
**Endpoint:** `GET /api/alumni/{id}`  
**Deskripsi:** Mendapatkan detail alumni dengan relasi user dan pekerjaan

**Response:**
```json
{
  "data": {
    "id": 1,
    "user_id": 2,
    "user": {
      "id": 2,
      "username": "jane_alumni",
      "email": "jane@example.com",
      "role": "user"
    },
    "nim": "2018001",
    "nama": "Jane Smith",
    "jurusan": "Teknik Informatika",
    "angkatan": 2018, 
    "tahun_lulus": 2022,
    "no_telepon": "08123456789",
    "alamat": "Jakarta Selatan",
    "pekerjaan_alumni": [
      {
        "id": 1,
        "nama_perusahaan": "PT Tech Indonesia",
        "posisi_jabatan": "Software Engineer"
      }
    ]
  }
}
```

### 14. Create Alumni
**Endpoint:** `POST /api/alumni`  
**Deskripsi:** Membuat data alumni baru (terhubung dengan user)

**Request:**
```json
{
  "user_id": 2,
  "nim": "2018001",
  "nama": "Jane Smith",
  "jurusan": "Teknik Informatika",
  "angkatan": 2018,
  "tahun_lulus": 2022,
  "no_telepon": "08123456789",
  "alamat": "Jakarta Selatan"
}
```

### 15. Update Alumni  
**Endpoint:** `PUT /api/alumni/{id}`  
**Deskripsi:** Update data alumni

**Request:**
```json
{
  "user_id": 2,
  "nim": "2018001", 
  "nama": "Jane Smith Updated",
  "jurusan": "Sistem Informasi",
  "angkatan": 2018,
  "tahun_lulus": 2022,
  "no_telepon": "08123456789",
  "alamat": "Jakarta Pusat"
}
```

### 16. Delete Alumni
**Endpoint:** `DELETE /api/alumni/{id}`  
**Deskripsi:** Menghapus alumni (Hard delete)

**Response:**
```json
{
  "message": "Alumni berhasil dihapus"
}
```

---

## 💼 Pekerjaan Alumni CRUD Operations (with Soft Delete)

**Headers Required:** `Authorization: Bearer <jwt_token>`

### 17. Get All Pekerjaan (with Pagination & Search)
**Endpoint:** `GET /api/pekerjaan?page=1&per_page=10&search=software&sort_by=created_at&sort_order=desc`  
**Deskripsi:** Mendapatkan data pekerjaan alumni dengan pencarian (exclude soft deleted)

**Query Parameters:**
- `page`: Halaman (default: 1)
- `per_page`: Data per halaman (default: 10)
- `search`: Kata kunci pencarian (posisi_jabatan, nama_perusahaan, bidang_industri, lokasi_kerja, alumni.nama)
- `sort_by`: Kolom sorting (default: id) 
- `sort_order`: Urutan sorting - ASC/DESC atau asc/desc (default: ASC)

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "alumni_id": 1,
      "alumni": {
        "id": 1,
        "nama": "Jane Smith",
        "nim": "2018001"
      },
      "nama_perusahaan": "PT Tech Indonesia",
      "posisi_jabatan": "Software Engineer", 
      "bidang_industri": "Technology",
      "lokasi_kerja": "Jakarta",
      "gaji_range": "15-20 juta",
      "tanggal_mulai_kerja": "2022-07-01T00:00:00Z",
      "tanggal_selesai_kerja": null,
      "status_pekerjaan": "aktif",
      "deskripsi_pekerjaan": "Mengembangkan aplikasi web",
      "created_at": "2025-09-25T12:00:00Z"
    }
  ],
  "current_page": 1,
  "per_page": 10,
  "total_data": 1,
  "total_pages": 1
}
```

### 18. Get Pekerjaan by ID
**Endpoint:** `GET /api/pekerjaan/{id}`  
**Deskripsi:** Mendapatkan detail pekerjaan alumni berdasarkan ID

**Response:**
```json
{
  "data": {
    "id": 1,
    "alumni_id": 1,
    "alumni": {
      "id": 1,
      "nama": "Jane Smith",
      "nim": "2018001", 
      "jurusan": "Teknik Informatika"
    },
    "nama_perusahaan": "PT Tech Indonesia",
    "posisi_jabatan": "Software Engineer",
    "bidang_industri": "Technology",
    "lokasi_kerja": "Jakarta", 
    "gaji_range": "15-20 juta",
    "tanggal_mulai_kerja": "2022-07-01T00:00:00Z",
    "status_pekerjaan": "aktif",
    "deskripsi_pekerjaan": "Mengembangkan aplikasi web dan mobile"
  }
}
```

### 19. Create Pekerjaan Alumni
**Endpoint:** `POST /api/pekerjaan`  
**Deskripsi:** Membuat data pekerjaan alumni baru

**Request:**
```json
{
  "alumni_id": 1,
  "nama_perusahaan": "PT Tech Indonesia",
  "posisi_jabatan": "Software Engineer",
  "bidang_industri": "Technology",
  "lokasi_kerja": "Jakarta",
  "gaji_range": "15-20 juta", 
  "tanggal_mulai_kerja": "2022-07-01T00:00:00Z",
  "tanggal_selesai_kerja": null,
  "status_pekerjaan": "aktif",
  "deskripsi_pekerjaan": "Mengembangkan aplikasi web dan mobile"
}
```

### 20. Update Pekerjaan Alumni
**Endpoint:** `PUT /api/pekerjaan/{id}`  
**Deskripsi:** Update data pekerjaan alumni

**Request:**
```json
{
  "alumni_id": 1,
  "nama_perusahaan": "PT Tech Indonesia", 
  "posisi_jabatan": "Senior Software Engineer",
  "bidang_industri": "Technology",
  "lokasi_kerja": "Jakarta",
  "gaji_range": "20-25 juta",
  "status_pekerjaan": "aktif",
  "deskripsi_pekerjaan": "Lead development team"
}
```

### 21. Hard Delete Pekerjaan
**Endpoint:** `DELETE /api/pekerjaan/{id}`  
**Deskripsi:** Menghapus pekerjaan secara permanen (Hard delete)

**Response:**
```json
{
  "message": "Pekerjaan berhasil dihapus"
}
```

---

## 🗑️ Soft Delete Operations (NEW FEATURES)

### 22. Soft Delete Pekerjaan
**Endpoint:** `DELETE /api/pekerjaan/soft/{id}`  
**Deskripsi:** Menghapus pekerjaan sementara (dapat dikembalikan)  
**Authorization:** 
- Admin: dapat soft delete semua pekerjaan
- User: hanya pekerjaan alumni miliknya sendiri

**Response:**
```json
{
  "message": "Pekerjaan berhasil dihapus sementara"
}
```

### 23. Restore Pekerjaan  
**Endpoint:** `POST /api/pekerjaan/restore/{id}`  
**Deskripsi:** Mengembalikan pekerjaan yang di-soft delete

**Response:**
```json
{
  "message": "Pekerjaan berhasil dikembalikan"
}
```

### 24. Get Deleted Pekerjaan
**Endpoint:** `GET /api/pekerjaan/deleted`  
**Deskripsi:** Mendapatkan semua pekerjaan yang di-soft delete (Admin only)

**Response:**
```json
{
  "data": [
    {
      "id": 2,
      "alumni_id": 1, 
      "nama_perusahaan": "PT Old Company",
      "posisi_jabatan": "Junior Developer",
      "deleted_at": "2025-09-25T15:30:00Z"
    }
  ]
}
```

### 25. Soft Delete All Pekerjaan by Alumni
**Endpoint:** `DELETE /api/alumni/{alumni_id}/pekerjaan/soft`  
**Deskripsi:** Soft delete semua pekerjaan milik alumni tertentu (Admin only)

**Response:**
```json
{
  "message": "Semua pekerjaan alumni berhasil dihapus sementara"
}
```

---

## � Sorting & Search Examples

### Sorting Examples
```bash
# Sort users by username ascending
GET /api/users?sort_by=username&sort_order=ASC

# Sort users by username descending (both uppercase and lowercase work)
GET /api/users?sort_by=username&sort_order=DESC
GET /api/users?sort_by=username&sort_order=desc

# Sort alumni by angkatan (year) descending
GET /api/alumni?sort_by=angkatan&sort_order=DESC

# Sort pekerjaan by company name ascending
GET /api/pekerjaan?sort_by=nama_perusahaan&sort_order=ASC

# Sort by creation date (newest first)
GET /api/pekerjaan?sort_by=created_at&sort_order=DESC
```

### Search Examples  
```bash
# Search users by username, email, or role
GET /api/users?search=admin

# Search mahasiswa by nim, name, email, jurusan, or angkatan
GET /api/mahasiswa?search=teknik

# Search alumni by nim, name, jurusan, graduation year, or user email
GET /api/alumni?search=2020

# Search pekerjaan by position, company, industry, location, or alumni name
GET /api/pekerjaan?search=software
```

### Combined Search + Sort Examples
```bash
# Search for "teknik" and sort by name descending
GET /api/alumni?search=teknik&sort_by=nama&sort_order=DESC

# Search for "engineer" and sort by company name ascending  
GET /api/pekerjaan?search=engineer&sort_by=nama_perusahaan&sort_order=ASC

# Search with pagination and sorting
GET /api/users?search=admin&sort_by=username&sort_order=ASC&page=1&per_page=5
```

---

## �🔒 Authorization Rules

| Role | Permissions |
|------|-------------|
| **Admin** | • Full access ke semua endpoint<br>• Dapat CRUD semua data<br>• Dapat soft delete semua pekerjaan<br>• Dapat melihat data deleted |
| **User** | • Dapat melihat semua data<br>• Dapat CRUD data sendiri<br>• Dapat soft delete hanya pekerjaan alumni miliknya<br>• Tidak dapat melihat data deleted |

---

## 🚫 Error Responses

### Authentication Error
```json
{
  "error": "Token tidak ditemukan"
}
```

### Authorization Error  
```json
{
  "error": "Access denied. You can only delete your own job records."
}
```

### Validation Error
```json
{
  "error": "Email sudah digunakan"
}
```

### Not Found Error
```json
{
  "error": "Data tidak ditemukan"
}
```

---

## ⚙️ Setup & Run

1. **Clone repository**
2. **Setup PostgreSQL database**
3. **Run:** `go run main.go` 
4. **Server runs on:** `http://localhost:8080`

---

## 📈 Database Schema

### Users Table
- id (Primary Key)
- username (Unique)
- email (Unique) 
- password (Hashed)
- role (admin/user)
- is_active (Boolean)
- created_at, updated_at

### Alumni Table  
- id (Primary Key)
- user_id (Foreign Key → users.id)
- nim (Unique)
- nama
- jurusan
- angkatan
- tahun_lulus
- no_telepon
- alamat
- created_at, updated_at

### Pekerjaan Alumni Table
- id (Primary Key)
- alumni_id (Foreign Key → alumnis.id) 
- nama_perusahaan
- posisi_jabatan
- bidang_industri
- lokasi_kerja
- gaji_range
- tanggal_mulai_kerja
- tanggal_selesai_kerja
- status_pekerjaan
- deskripsi_pekerjaan
- created_at, updated_at
- **deleted_at** (Soft Delete)

---

## 🎯 Key Features

✅ **JWT Authentication & RBAC**  
✅ **Clean Architecture Pattern**  
✅ **Pagination & Search**  
✅ **Soft Delete with Authorization**  
✅ **Database Relations**  
✅ **Input Validation**  
✅ **Error Handling**  
✅ **RESTful API Design**

```go
type User struct {
    ID        int       `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    Role      string    `json:"role"`
    CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type LoginResponse struct {
    User  User   `json:"user"`
    Token string `json:"token"`
}

type JWTClaims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}
```

## Endpoint Authentication

### 1. Register User
```
POST /auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "role": "user" // optional, default: "user"
}
```

**Response:**
```json
{
  "message": "User berhasil didaftarkan",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "role": "user",
    "created_at": "2025-09-14T10:00:00Z"
  }
}
```

### 2. Login
```
POST /auth/login
Content-Type: application/json

{
  "username": "johndoe",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "Login berhasil",
  "data": {
    "user": {
      "id": 1,
      "username": "johndoe",
      "email": "john@example.com",
      "role": "user",
      "created_at": "2025-09-14T10:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

## Protected Endpoints

Semua endpoint dibawah memerlukan header:
```
Authorization: Bearer <your_jwt_token>
```

### User Profile
```
GET /api/profile
```

### User Management (Admin Only)
```
GET /api/users              # Get all users
GET /api/users/count         # Get user count
GET /api/users/{id}          # Get user by ID
PUT /api/users/{id}          # Update user
DELETE /api/users/{id}       # Delete user
```

### Mahasiswa Endpoints
```
GET /api/mahasiswa           # All users
GET /api/mahasiswa/count     # All users
GET /api/mahasiswa/{id}      # All users
POST /api/mahasiswa          # Admin/Moderator only
PUT /api/mahasiswa/{id}      # Admin/Moderator only
DELETE /api/mahasiswa/{id}   # Admin only
```

### Alumni Endpoints
```
GET /api/alumni              # All users
GET /api/alumni/count        # All users
GET /api/alumni/{id}         # All users
POST /api/alumni             # Admin/Moderator only
PUT /api/alumni/{id}         # Admin/Moderator only
DELETE /api/alumni/{id}      # Admin only
```

### Pekerjaan Alumni Endpoints
```
GET /api/pekerjaan                    # All users
GET /api/pekerjaan/count              # All users
GET /api/pekerjaan/{id}               # All users
GET /api/pekerjaan/alumni/{alumni_id} # All users
POST /api/pekerjaan                   # Admin/Moderator only
PUT /api/pekerjaan/{id}               # Admin/Moderator only
DELETE /api/pekerjaan/{id}            # Admin only
```

### Perusahaan Statistics
```
GET /api/perusahaan/{nama_perusahaan} # All users
```

**Response:**
```json
{
  "total_alumni": 5,
  "nama_perusahaan": "PT Telkom Indonesia",
  "message": "Data jumlah alumni di perusahaan berhasil diambil"
}
```

## Architecture Changes

### 1. **Utils Layer**
- `utils/password.go`: Password hashing functions
- `utils/jwt.go`: JWT token generation and validation

### 2. **Services Layer** (formerly Controllers)
- Clean separation of HTTP handling logic
- Better naming convention
- Consistent error handling

### 3. **Middleware Updates**
- Uses utils for JWT validation
- Simplified context handling
- Better error responses

## Role Permissions

| Role | Permissions |
|------|-------------|
| **admin** | Full access: Create, Read, Update, Delete semua data |
| **moderator** | Read all data, Create/Update mahasiswa/alumni/pekerjaan |
| **user** | Read access only |

## Testing dengan Postman/curl

### 1. Register Admin User
```bash
curl -X POST http://localhost:3000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "admin123",
    "role": "admin"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:3000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### 3. Access Protected Endpoint
```bash
curl -X GET http://localhost:3000/api/mahasiswa \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Environment Variables

Untuk production, set environment variable:
```
JWT_SECRET=your-super-secret-jwt-key-for-production
```

Jika tidak di-set, akan menggunakan default secret (hanya untuk development).

## Security Features

1. **Password Hashing**: Menggunakan bcrypt dengan default cost
2. **JWT Expiration**: Token berlaku 24 jam
3. **Role Validation**: Setiap endpoint memiliki role requirement
4. **Input Validation**: Validasi username dan password minimum
5. **Error Handling**: Response error yang informatif tanpa expose sensitive data
6. **Clean Architecture**: Separation of concerns dengan utils layer
