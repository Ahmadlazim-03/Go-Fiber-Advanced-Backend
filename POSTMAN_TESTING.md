# POSTMAN Testing Guide - CRUD Go Fiber PostgreSQL

## Base URL
```
http://localhost:3000
```

## Authentication Flow

### 1. Register New User (Buat User Baru Dulu)
**POST** `/api/register`

**Headers:**
```json
{
  "Content-Type": "application/json"
}
```

**Body (JSON):**
```json
{
  "username": "test",
  "email": "test@gmail.com",
  "password": "Pembelajaranjarakjauh@123",
  "role": "user"
}
```

### 2. Login (Get Bearer Token)
**POST** `/api/login`

**Headers:**
```json
{
  "Content-Type": "application/json"
}
```

**Body (JSON):**
```json
{
  "email": "test@gmail.com",
  "password": "Pembelajaranjarakjauh@123"
}
```

**Expected Response:**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 2,
    "username": "test",
    "email": "test@gmail.com",
    "role": "user"
  }
}
```

### Alternative: Login dengan Admin
**Body (JSON):**
```json
{
  "email": "admin@example.com",
  "password": "admin123"
}
```

---

## Mahasiswa CRUD Operations

### 1. Get All Mahasiswa
**GET** `/api/mahasiswa`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

### 2. Get Mahasiswa Count
**GET** `/api/mahasiswa/count`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

### 3. Create New Mahasiswa
**POST** `/api/mahasiswa`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

**Body (JSON):**
```json
{
  "nim": "2021001001",
  "nama": "Ahmad Budi Santoso",
  "email": "ahmad.budi@student.ac.id",
  "jurusan": "Teknik Informatika",
  "angkatan": 2021
}
```

### 4. Get Mahasiswa by ID
**GET** `/api/mahasiswa/1`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

### 5. Update Mahasiswa
**PUT** `/api/mahasiswa/1`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

**Body (JSON):**
```json
{
  "nim": "2021001001",
  "nama": "Ahmad Budi Santoso Updated",
  "email": "ahmad.budi.updated@student.ac.id",
  "jurusan": "Sistem Informasi",
  "angkatan": 2021
}
```

### 6. Delete Mahasiswa
**DELETE** `/api/mahasiswa/1`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

---

## Alumni CRUD Operations

### 1. Get All Alumni
**GET** `/api/alumni`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

### 2. Get Alumni Count
**GET** `/api/alumni/count`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

### 3. Create New Alumni
**POST** `/api/alumni`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

**Body (JSON):**
```json
{
  "nim": "2018001001",
  "nama": "Siti Nur Aisyah",
  "email": "siti.aisyah@alumni.ac.id",
  "jurusan": "Teknik Informatika",
  "tahun_lulus": 2022,
  "alamat": "Jl. Merdeka No. 123, Jakarta Pusat"
}
```

### 4. Get Alumni by ID
**GET** `/api/alumni/1`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

### 5. Update Alumni
**PUT** `/api/alumni/1`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

**Body (JSON):**
```json
{
  "nim": "2018001001",
  "nama": "Siti Nur Aisyah, S.Kom",
  "email": "siti.aisyah.updated@alumni.ac.id",
  "jurusan": "Teknik Informatika",
  "tahun_lulus": 2022,
  "alamat": "Jl. Sudirman No. 456, Jakarta Selatan"
}
```

### 6. Delete Alumni
**DELETE** `/api/alumni/1`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

---

## Pekerjaan Alumni CRUD Operations

### 1. Get All Pekerjaan
**GET** `/api/pekerjaan`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

### 2. Create New Pekerjaan
**POST** `/api/pekerjaan`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

**Body (JSON):**
```json
{
  "alumni_id": 1,
  "nama_perusahaan": "PT. Tech Indonesia",
  "posisi_jabatan": "Software Developer",
  "bidang_industri": "Teknologi Informasi",
  "lokasi_kerja": "Jakarta",
  "gaji": 8000000,
  "tahun_mulai": 2022,
  "status": "aktif"
}
```

### 3. Get Pekerjaan by ID
**GET** `/api/pekerjaan/1`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

### 4. Update Pekerjaan
**PUT** `/api/pekerjaan/1`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

**Body (JSON):**
```json
{
  "alumni_id": 1,
  "nama_perusahaan": "PT. Tech Indonesia",
  "posisi_jabatan": "Senior Software Developer",
  "bidang_industri": "Teknologi Informasi",
  "lokasi_kerja": "Jakarta",
  "gaji": 12000000,
  "tahun_mulai": 2022,
  "status": "aktif"
}
```

### 5. Delete Pekerjaan
**DELETE** `/api/pekerjaan/1`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

---

## User Management (Admin Only)

### 1. Get All Users
**GET** `/api/users`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

### 2. Get User by ID
**GET** `/api/users/1`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

### 3. Update User
**PUT** `/api/users/1`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

**Body (JSON):**
```json
{
  "username": "updated_user",
  "email": "updated@example.com",
  "role": "moderator"
}
```

### 4. Delete User
**DELETE** `/api/users/1`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

---

## Dashboard Data

### Get Dashboard Statistics
**GET** `/api/dashboard`

**Headers:**
```json
{
  "Authorization": "Bearer YOUR_TOKEN_HERE",
  "Content-Type": "application/json"
}
```

**Expected Response:**
```json
{
  "total_mahasiswa": 150,
  "total_alumni": 1200,
  "total_pekerjaan": 800,
  "total_users": 25
}
```

---

## Testing Steps in Postman

1. **First**: Register/Login with user credentials to get Bearer token
2. **Copy the token** from login response
3. **Use the token** in Authorization header for all requests
4. **Test permissions based on role**:
   - **Admin**: Can perform ALL operations (GET, POST, PUT, DELETE)
   - **User**: Can only perform READ operations (GET only)

## Role-Based Access Control (RBAC)

### Admin Role (`admin`)
- ✅ **GET** (Read) - All endpoints
- ✅ **POST** (Create) - All data
- ✅ **PUT** (Update) - All data  
- ✅ **DELETE** (Delete) - All data
- ✅ **User Management** - Full access

### User Role (`user`)
- ✅ **GET** (Read) - All endpoints
- ❌ **POST** (Create) - Forbidden
- ❌ **PUT** (Update) - Forbidden
- ❌ **DELETE** (Delete) - Forbidden
- ❌ **User Management** - No access

## Default Credentials
- **Admin**: `admin@example.com` / `admin123`
- **Test User**: `test@gmail.com` / `Pembelajaranjarakjauh@123`

## Important Notes
- All `/api/*` endpoints require Bearer token authentication
- **Admin role**: Full CRUD access to all resources
- **User role**: Read-only access (GET requests only)
- **Moderator role**: Removed from system
- Token expires in 24 hours
