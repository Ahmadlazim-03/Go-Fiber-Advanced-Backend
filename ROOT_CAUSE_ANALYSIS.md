# 🔍 ROOT CAUSE ANALYSIS - PocketBase Login Failure

## Masalah yang Terjadi

```
Login failed: The remote server returned an error: (401) Unauthorized.
Error details: {"error":"Email atau password salah"}
```

## Investigasi

### ✅ Yang Sudah Benar:
1. Server berjalan di `http://localhost:8080`
2. Database type: `pocketbase`
3. PocketBase URL: `https://pocketbase-production-521e.up.railway.app`
4. User berhasil dibuat (register success)

### ❌ Yang Bermasalah:
1. Login selalu gagal dengan 401 Unauthorized
2. Error message: "Email atau password salah"

## Root Causes yang Ditemukan

### 1. **User `is_active = False` (CONFIRMED)**
Di PocketBase admin panel terlihat:
- ✗ test@example.com → is_active: **False**
- ✗ admin@example.com → is_active: **False**

**Penyebab:** Code tidak set `is_active = true` saat register

**Lokasi Bug:** `services/auth_service.go` line ~80
```go
newUser := &models.User{
    Username: req.Username,
    Email:    req.Email,
    Password: hashedPassword,
    Role:     role,
    // IsActive: true,  ← MISSING!
}
```

### 2. **Password Hashing Issue**
Password di-hash dengan bcrypt di aplikasi, tapi PocketBase juga hash password sendiri.
Jadi terjadi **double hashing**!

**Flow sekarang (WRONG):**
```
User input: "12345678"
  ↓
App hash dengan bcrypt: "$2a$10$abc..."
  ↓  
PocketBase terima: "$2a$10$abc..."
  ↓
PocketBase hash lagi: "$2a$10$xyz..."  ← DOUBLE HASH!
```

**Saat login:**
```
User input: "12345678"
  ↓
App hash: "$2a$10$abc..."
  ↓
Compare dengan PocketBase: "$2a$10$xyz..."  ← NOT MATCH!
```

## Solusi

### Solution 1: Perbaiki Code (RECOMMENDED)

#### A. Fix `is_active` Default Value

Edit `services/auth_service.go`:

```go
newUser := &models.User{
    Username: req.Username,
    Email:    req.Email,
    Password: hashedPassword,
    Role:     role,
    IsActive: true,  // ← ADD THIS
}
```

#### B. Fix Password Hashing untuk PocketBase

Edit `repositories/pocketbase/user_repository_pocketbase.go`:

**BEFORE (WRONG):**
```go
func (r *UserRepositoryPocketBase) Create(user *models.User) error {
    payload := map[string]interface{}{
        "username":       user.Username,
        "email":          user.Email,
        "password":       user.Password,  // ← Already hashed by bcrypt!
        "passwordConfirm": user.Password,
        "role":           user.Role,
    }
    // ...
}
```

**AFTER (CORRECT):**
```go
func (r *UserRepositoryPocketBase) Create(user *models.User) error {
    // PocketBase akan hash password sendiri, jadi kirim plain password!
    // Tapi kita sudah terima hashed password dari service...
    // MASALAHNYA ADA DI SINI!
    
    payload := map[string]interface{}{
        "username":        user.Username,
        "email":           user.Email,
        "password":        user.Password,  // Sudah di-hash
        "passwordConfirm": user.Password,
        "role":            user.Role,
        "is_active":       user.IsActive,  // ← ADD THIS
    }
    // ...
}
```

**BETTER SOLUTION:**

Ubah flow - jangan hash password di service untuk PocketBase!

Edit `services/auth_service.go`:

```go
func (s *AuthService) Register(c *fiber.Ctx) error {
    // ... validation code ...

    // Check database type
    if database.IsPocketBase() {
        // PocketBase hash password sendiri - kirim plain password
        newUser := &models.User{
            Username: req.Username,
            Email:    req.Email,
            Password: req.Password,  // ← PLAIN PASSWORD for PocketBase
            Role:     role,
            IsActive: true,
        }
    } else {
        // PostgreSQL/MongoDB - hash password dulu
        hashedPassword, _ := utils.HashPassword(req.Password)
        newUser := &models.User{
            Username: req.Username,
            Email:    req.Email,
            Password: hashedPassword,
            Role:     role,
            IsActive: true,
        }
    }
    
    // Create user
    if err := s.userRepo.Create(newUser); err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Gagal mendaftarkan user"})
    }
    // ...
}
```

### Solution 2: Manual Fix (QUICK)

**Untuk test SEKARANG:**

1. **Buka PocketBase Admin:**
   ```
   https://pocketbase-production-521e.up.railway.app/_/
   ```

2. **Edit User:**
   - Collections → users
   - Click edit pada `ahmad@gmail.com`
   - Set `is_active` = **True**
   - **IMPORTANT:** Reset password manually:
     - Delete existing password hash
     - Set new password: `12345678`
     - PocketBase akan hash dengan benar
   - Save

3. **Test Login:**
   ```powershell
   .\test_crud_complete.ps1
   ```

### Solution 3: Create User via PocketBase Admin (EASIEST)

1. **Buka:** `https://pocketbase-production-521e.up.railway.app/_/`
2. **Collections** → **users** → **New record**
3. **Fill form:**
   - email: `ahmad@gmail.com`
   - password: `12345678`
   - username: `ahmad`
   - role: `admin`
   - is_active: ✅ **True**
4. **Save**
5. **Test:** `.\test_crud_complete.ps1`

## Verification Checklist

- [ ] User exists in PocketBase
- [ ] `is_active = True`
- [ ] Password set correctly (via PocketBase, not double-hashed)
- [ ] Login works
- [ ] CRUD operations work

## Next Steps

**IMMEDIATE (untuk test sekarang):**
→ Gunakan **Solution 3** - buat user baru via PocketBase admin

**PERMANENT (untuk fix code):**
→ Implement **Solution 1B** - conditional password hashing based on database type

---

**Status:** Root cause identified - Double password hashing + is_active=false issue
