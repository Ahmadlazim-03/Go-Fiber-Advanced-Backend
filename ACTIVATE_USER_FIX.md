# 🔧 Quick Fix: Activate User in PocketBase

## Problem
User berhasil dibuat tapi `is_active = False`, sehingga login gagal dengan error:
```
✗ Error: The remote server returned an error: (401) Unauthorized.
✗ Login FAILED: {"error":"Email atau password salah"}
```

## Solution: Activate User via PocketBase Admin

### Step 1: Open PocketBase Admin
```
https://pocketbase-production-521e.up.railway.app/_/
```

### Step 2: Go to Users Collection
1. Click **Collections** (left sidebar)
2. Click **users**

### Step 3: Edit User
1. Find user: `test@example.com` (username: testuser)
2. Click the **pencil icon** (Edit) on the right
3. Find field: **is_active**
4. Change from `False` → **`True`** ✅
5. Click **Save**

### Step 4: Verify
User should now show:
- ✅ Email: test@example.com
- ✅ Username: testuser
- ✅ Role: admin
- ✅ **is_active: True** ← IMPORTANT!

### Step 5: Test Login Again

Run test script:
```powershell
.\test_crud_complete.ps1
```

Or test manually:
```powershell
$body = @{email="test@example.com"; password="test12345"} | ConvertTo-Json
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/login" `
    -Method POST `
    -ContentType "application/json" `
    -Body $body

Write-Host "Success! Token: $($response.token.Substring(0,30))..."
```

## Alternative: Fix via Code

Update `auth_service.go` to set `is_active = true` by default:

```go
// In Register function
newUser := &models.User{
    Username: req.Username,
    Email:    req.Email,
    Password: hashedPassword,
    Role:     role,
    IsActive: true,  // ← Add this line
}
```

Then restart server and register new user.

## Quick Test Command

After activating user:
```powershell
# Login
$login = @{email="test@example.com"; password="test12345"} | ConvertTo-Json
$auth = Invoke-RestMethod -Uri "http://localhost:8080/api/login" -Method POST -ContentType "application/json" -Body $login
Write-Host "✓ Login success! Token: $($auth.token.Substring(0,20))..."

# Create Mahasiswa
$headers = @{Authorization="Bearer $($auth.token)"; "Content-Type"="application/json"}
$mhs = @{nim="2025001"; nama="Test User"; jurusan="TI"; angkatan=2025; email="test@test.com"} | ConvertTo-Json
$created = Invoke-RestMethod -Uri "http://localhost:8080/api/mahasiswa" -Method POST -Headers $headers -Body $mhs
Write-Host "✓ Mahasiswa created! ID: $($created.id)"
```

---

**Action Required:** 
1. Go to PocketBase admin panel
2. Edit user `test@example.com`
3. Set `is_active = True`
4. Save
5. Run test script!
