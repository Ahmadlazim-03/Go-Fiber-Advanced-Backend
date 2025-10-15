# 🗄️ Multi-Database CRUD Application

## Overview
Aplikasi CRUD (Create, Read, Update, Delete) berbasis Go Fiber yang mendukung **3 jenis database**:
- **PostgreSQL** - Relational database
- **MongoDB** - NoSQL document database  
- **PocketBase** - All-in-one backend solution

## 🎯 Features

✅ **Multi-Database Support** - Ganti database hanya dengan mengubah environment variable
✅ **RESTful API** - Endpoints standar untuk semua operasi
✅ **Authentication & Authorization** - JWT-based auth dengan role-based access
✅ **Repository Pattern** - Clean architecture dengan interface abstraction
✅ **Auto Migrations** - Database schema dibuat otomatis
✅ **Pagination** - Support pagination untuk semua list endpoints
✅ **Soft Delete** - Data bisa di-restore (untuk pekerjaan alumni)
✅ **Web UI** - Simple web interface untuk testing

## 📦 Tech Stack

- **Framework**: Go Fiber v2
- **Databases**: PostgreSQL, MongoDB, PocketBase
- **Authentication**: JWT (JSON Web Tokens)
- **Password Hashing**: bcrypt
- **ORM**: GORM (untuk PostgreSQL)
- **MongoDB Driver**: Official Go MongoDB Driver
- **HTTP Client**: net/http (untuk PocketBase)

## 🚀 Quick Start

### 1. Clone Repository
```bash
git clone https://github.com/Ahmadlazim-03/CRUD-Go-Fiber-PostgreSQL.git
cd CRUD-Go-Fiber-PostgreSQL
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Configure Environment

Pilih salah satu database dan edit file `.env`:

#### Option A: PostgreSQL
```env
DB_TYPE=postgres
POSTGRES_DSN=postgresql://postgres:password@host:port/database
```

#### Option B: MongoDB
```env
DB_TYPE=mongodb
MONGODB_URI=mongodb://username:password@host:port/database?authSource=admin
MONGODB_DATABASE=database_name
```

#### Option C: PocketBase
```env
DB_TYPE=pocketbase
POCKETBASE_URL=https://your-pocketbase-url.com
POCKETBASE_ADMIN_EMAIL=admin@example.com
POCKETBASE_ADMIN_PASSWORD=admin_password
```

### 4. Run Application
```bash
go run main.go
```

Server akan berjalan di: **http://localhost:8080**

## 📁 Project Structure

```
.
├── main.go                      # Entry point
├── .env                         # Environment configuration
├── database/
│   ├── connection.go            # Database connections
│   └── migration/
│       ├── migrations.go        # Migration wrapper
│       ├── migrations_postgres.go
│       ├── migrations_mongodb.go
│       └── migrations_pocketbase.go
├── models/
│   ├── user.go                  # User model
│   ├── mahasiswa.go             # Student model
│   ├── alumni.go                # Alumni model
│   └── pagination.go            # Pagination helpers
├── repositories/
│   ├── interface/               # Repository interfaces
│   ├── postgres/                # PostgreSQL implementations
│   ├── mongodb/                 # MongoDB implementations
│   └── pocketbase/              # PocketBase implementations
├── services/
│   ├── auth_service.go          # Authentication logic
│   ├── mahasiswa_service.go     # Student business logic
│   ├── alumni_service.go        # Alumni business logic
│   └── pekerjaan_alumni_service.go
├── middleware/
│   └── auth.go                  # JWT authentication middleware
├── routes/
│   └── routes.go                # API route definitions
├── utils/
│   ├── jwt.go                   # JWT utilities
│   └── password.go              # Password hashing
├── static/
│   ├── css/
│   └── js/
└── templates/
    ├── login.html
    ├── register.html
    └── index.html
```

## 📚 API Documentation

### Authentication

#### Register
```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "role": "user"
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}

Response:
{
  "user": {...},
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Mahasiswa (Students)

```http
GET    /api/mahasiswa              # List all students
GET    /api/mahasiswa/:id          # Get student by ID
POST   /api/mahasiswa              # Create new student
PUT    /api/mahasiswa/:id          # Update student
DELETE /api/mahasiswa/:id          # Delete student
```

### Alumni

```http
GET    /api/alumni                 # List all alumni
GET    /api/alumni/:id             # Get alumni by ID
POST   /api/alumni                 # Create new alumni
PUT    /api/alumni/:id             # Update alumni
DELETE /api/alumni/:id             # Delete alumni
```

### Pekerjaan Alumni (Alumni Jobs)

```http
GET    /api/pekerjaan              # List all jobs
GET    /api/pekerjaan/:id          # Get job by ID
POST   /api/pekerjaan              # Create new job
PUT    /api/pekerjaan/:id          # Update job
DELETE /api/pekerjaan/:id          # Soft delete job
GET    /api/pekerjaan/trash        # List deleted jobs
POST   /api/pekerjaan/restore/:id  # Restore deleted job
```

## 🗄️ Database Comparison

| Feature | PostgreSQL | MongoDB | PocketBase |
|---------|-----------|---------|------------|
| Type | Relational | NoSQL | Embedded SQLite + API |
| Schema | Strict | Flexible | Defined via API |
| Relationships | Foreign Keys | References | Manual |
| Transactions | ✅ | ✅ | Limited |
| Migrations | GORM Auto | Collections | HTTP API |
| Indexes | ✅ | ✅ | ✅ |
| Full-text Search | ✅ | ✅ | ✅ |
| Admin UI | External | External | Built-in |
| Best For | Production | Scalable apps | Prototyping |

## 🔧 Configuration Details

### PostgreSQL
```env
# Railway Example
POSTGRES_DSN=postgresql://postgres:password@host.railway.app:5432/railway

# Local Example
POSTGRES_DSN=host=localhost user=postgres password=password dbname=mydb port=5432 sslmode=disable
```

**Tables Created:**
- `users` - User accounts
- `mahasiswas` - Student data
- `alumnis` - Alumni data
- `pekerjaan_alumnis` - Alumni job data

### MongoDB
```env
# Railway Example
MONGODB_URI=mongodb://mongo:password@host.railway.app:27017/railway?authSource=admin
MONGODB_DATABASE=railway

# Local Example  
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=mydb
```

**Collections Created:**
- `users`
- `mahasiswas`
- `alumnis`
- `pekerjaan_alumnis`

### PocketBase
```env
# Railway Example
POCKETBASE_URL=https://pocketbase-production.up.railway.app
POCKETBASE_ADMIN_EMAIL=admin@example.com
POCKETBASE_ADMIN_PASSWORD=admin123456
```

**Collections Created:**
- `users` (built-in auth collection)
- `mahasiswas`
- `alumnis`
- `pekerjaan_alumnis`

**Admin Dashboard:** `{POCKETBASE_URL}/_/`

## 🎨 Web Interface

Akses web interface di **http://localhost:8080**

Pages:
- `/` - Home page
- `/login` - Login page
- `/register` - Register page
- `/welcome` - Dashboard (after login)

## 🔐 Default Admin Account

Setelah aplikasi start, akan otomatis membuat admin account:

```
Email: admin@example.com
Password: admin123
Role: admin
```

## 📝 Environment Variables

```env
# Database Selection
DB_TYPE=postgres|mongodb|pocketbase

# PostgreSQL
POSTGRES_DSN=connection_string

# MongoDB
MONGODB_URI=connection_string
MONGODB_DATABASE=database_name

# PocketBase
POCKETBASE_URL=https://your-url.com
POCKETBASE_ADMIN_EMAIL=admin@example.com
POCKETBASE_ADMIN_PASSWORD=password

# Server
SERVER_PORT=8080

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=24h
```

## 🧪 Testing

### Using curl

```bash
# Register user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@example.com","password":"pass123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"pass123"}'

# Get students (with token)
curl http://localhost:8080/api/mahasiswa \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### Using Postman

1. Import collection dari `API_DOCUMENTATION.md`
2. Set environment variable `token` setelah login
3. Test semua endpoints

## 📖 Documentation Files

- `README.md` - This file
- `API_DOCUMENTATION.md` - Complete API reference
- `DATABASE_SETUP.md` - Database setup guide
- `DATABASE_MIGRATIONS.md` - Migration system documentation
- `POCKETBASE_INTEGRATION.md` - PocketBase specific guide
- `IMPLEMENTATION_SUMMARY.md` - Technical implementation details
- `MIGRATION_GUIDE.md` - Migration guide between databases
- `DEPLOYMENT_STATUS.md` - Deployment status & notes

## 🐛 Troubleshooting

### PostgreSQL Connection Failed
```
Error: failed to connect to PostgreSQL
```
**Solution:** Check POSTGRES_DSN format and credentials

### MongoDB Connection Timeout
```
Error: server selection error: context deadline exceeded
```
**Solution:** 
- Check MONGODB_URI format
- Ensure `authSource=admin` is included
- Verify network connectivity

### PocketBase Authentication Failed
```
Error: authentication failed (status 400)
```
**Solution:** Check POCKETBASE_ADMIN_EMAIL and POCKETBASE_ADMIN_PASSWORD

### Port Already in Use
```
Error: bind: address already in use
```
**Solution:** Change SERVER_PORT in .env or kill process using port 8080

## 🚀 Deployment

### Railway Deployment

1. **Push to GitHub**
```bash
git add .
git commit -m "Initial commit"
git push origin main
```

2. **Deploy on Railway**
- Go to https://railway.app
- Create new project from GitHub repo
- Add database service (PostgreSQL/MongoDB)
- Set environment variables
- Deploy!

### Docker Deployment

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]
```

```bash
docker build -t crud-app .
docker run -p 8080:8080 --env-file .env crud-app
```

## 📊 Performance Considerations

### PostgreSQL
- ✅ Best for complex queries and transactions
- ✅ ACID compliance
- ⚠️ Requires proper indexing for large datasets

### MongoDB
- ✅ Great for flexible schemas
- ✅ Horizontal scaling
- ⚠️ May need disk space for indexes

### PocketBase
- ✅ Lightweight and fast for small/medium apps
- ✅ No external database needed
- ⚠️ SQLite limitations for high traffic

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

This project is licensed under the MIT License.

## 👨‍💻 Author

**Ahmadlazim-03**
- GitHub: [@Ahmadlazim-03](https://github.com/Ahmadlazim-03)

## 🙏 Acknowledgments

- Go Fiber Framework
- GORM ORM
- MongoDB Go Driver
- PocketBase
- Railway Hosting

---

**Happy Coding! 🚀**
