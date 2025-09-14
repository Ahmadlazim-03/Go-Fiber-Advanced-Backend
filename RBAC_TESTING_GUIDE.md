# Role-Based Access Control (RBAC) Testing Guide

## 🔐 **Permission Matrix**

| Endpoint | Method | Admin | User | Description |
|----------|--------|-------|------|-------------|
| **Authentication** |
| `/api/register` | POST | ✅ | ✅ | Anyone can register |
| `/api/login` | POST | ✅ | ✅ | Anyone can login |
| `/api/logout` | POST | ✅ | ✅ | Authenticated users |
| `/api/profile` | GET | ✅ | ✅ | Own profile |
| **Mahasiswa** |
| `/api/mahasiswa` | GET | ✅ | ✅ | View all students |
| `/api/mahasiswa/count` | GET | ✅ | ✅ | Count students |
| `/api/mahasiswa/:id` | GET | ✅ | ✅ | View student by ID |
| `/api/mahasiswa` | POST | ✅ | ❌ | Create student (Admin only) |
| `/api/mahasiswa/:id` | PUT | ✅ | ❌ | Update student (Admin only) |
| `/api/mahasiswa/:id` | DELETE | ✅ | ❌ | Delete student (Admin only) |
| **Alumni** |
| `/api/alumni` | GET | ✅ | ✅ | View all alumni |
| `/api/alumni/count` | GET | ✅ | ✅ | Count alumni |
| `/api/alumni/:id` | GET | ✅ | ✅ | View alumni by ID |
| `/api/alumni` | POST | ✅ | ❌ | Create alumni (Admin only) |
| `/api/alumni/:id` | PUT | ✅ | ❌ | Update alumni (Admin only) |
| `/api/alumni/:id` | DELETE | ✅ | ❌ | Delete alumni (Admin only) |
| **Pekerjaan** |
| `/api/pekerjaan` | GET | ✅ | ✅ | View all jobs |
| `/api/pekerjaan/count` | GET | ✅ | ✅ | Count jobs |
| `/api/pekerjaan/:id` | GET | ✅ | ✅ | View job by ID |
| `/api/pekerjaan/alumni/:id` | GET | ✅ | ✅ | View jobs by alumni |
| `/api/pekerjaan` | POST | ✅ | ❌ | Create job (Admin only) |
| `/api/pekerjaan/:id` | PUT | ✅ | ❌ | Update job (Admin only) |
| `/api/pekerjaan/:id` | DELETE | ✅ | ❌ | Delete job (Admin only) |
| **User Management** |
| `/api/users` | GET | ✅ | ❌ | View all users (Admin only) |
| `/api/users/count` | GET | ✅ | ❌ | Count users (Admin only) |
| `/api/users/:id` | GET | ✅ | ❌ | View user by ID (Admin only) |
| `/api/users/:id` | PUT | ✅ | ❌ | Update user (Admin only) |
| `/api/users/:id` | DELETE | ✅ | ❌ | Delete user (Admin only) |
| **Dashboard** |
| `/api/dashboard` | GET | ✅ | ✅ | Dashboard statistics |

## 🧪 **Testing Scenarios**

### **1. Test dengan Role Admin**
```bash
# Login sebagai admin
POST /api/login
{
  "email": "admin@example.com",
  "password": "admin123"
}

# Copy token, lalu test semua endpoint:
# ✅ Semua GET requests harus berhasil
# ✅ Semua POST requests harus berhasil  
# ✅ Semua PUT requests harus berhasil
# ✅ Semua DELETE requests harus berhasil
```

### **2. Test dengan Role User**
```bash
# Login sebagai user
POST /api/login
{
  "email": "test@gmail.com", 
  "password": "Pembelajaranjarakjauh@123"
}

# Copy token, lalu test:
# ✅ Semua GET requests harus berhasil
# ❌ POST /api/mahasiswa harus return 403 Forbidden
# ❌ PUT /api/mahasiswa/1 harus return 403 Forbidden  
# ❌ DELETE /api/mahasiswa/1 harus return 403 Forbidden
# ❌ GET /api/users harus return 403 Forbidden
```

## 🚨 **Expected Error Responses**

### **403 Forbidden (User trying admin operations):**
```json
{
  "error": "Akses ditolak: role tidak memiliki permission",
  "required_roles": ["admin"],
  "user_role": "user"
}
```

### **401 Unauthorized (No token):**
```json
{
  "error": "Token tidak ditemukan"
}
```

### **401 Unauthorized (Invalid token):**
```json
{
  "error": "Token tidak valid"
}
```

## 🎯 **Quick Test Commands**

### **Test User Permissions (Should fail):**
```bash
# Dengan token user, coba:
POST /api/mahasiswa    # Should return 403
PUT /api/alumni/1      # Should return 403  
DELETE /api/pekerjaan/1 # Should return 403
GET /api/users         # Should return 403
```

### **Test Admin Permissions (Should succeed):**
```bash
# Dengan token admin, coba:
POST /api/mahasiswa    # Should return 201
PUT /api/alumni/1      # Should return 200
DELETE /api/pekerjaan/1 # Should return 200  
GET /api/users         # Should return 200
```

## 🔑 **Role Summary**

### **Admin (`admin`)**
- **Permissions**: Full CRUD access
- **Use Case**: System administrator, data manager
- **Access Level**: Complete control over all resources

### **User (`user`)**  
- **Permissions**: Read-only access
- **Use Case**: End users, viewers, reporters
- **Access Level**: Can view data but cannot modify

### **Removed Roles**
- ~~**Moderator**~~ - Role dihapus dari sistem
