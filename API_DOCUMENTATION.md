# API Documentation - CRUD Go Fiber PostgreSQL dengan JWT Authentication

## Refactored Architecture

Aplikasi telah di-refactor dengan struktur yang lebih terorganisir:

### Folder Structure
```
├── models/
│   └── user.go              # User model + auth structs
├── utils/
│   ├── password.go          # Password hashing utilities
│   └── jwt.go               # JWT token utilities
├── middleware/
│   └── auth.go              # JWT & RBAC middleware
├── repositories/
│   └── user_repository.go   # User CRUD operations
├── usecases/
│   └── auth_usecase.go      # Auth business logic
├── services/
│   ├── auth_service.go      # Auth HTTP handlers (was controllers)
│   ├── alumni_service.go    # Alumni HTTP handlers
│   ├── mahasiswa_service.go # Mahasiswa HTTP handlers
│   └── pekerjaan_alumni_service.go # Pekerjaan HTTP handlers
└── routes/
    └── routes.go            # Protected routes setup
```

## Autentikasi dan Autorisasi

Aplikasi ini mengimplementasikan:
- **JWT (JSON Web Token)** untuk autentikasi
- **Role-Based Access Control (RBAC)** dengan 3 role: `admin`, `moderator`, `user`
- **Password hashing** menggunakan bcrypt
- **Utilities separation** untuk clean architecture

## User Model

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
