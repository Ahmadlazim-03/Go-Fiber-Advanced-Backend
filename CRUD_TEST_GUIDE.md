# 🧪 Testing CRUD Operations dengan PocketBase

## Cara Testing

### 1. Start Server

**Option A: Manual (Recommended)**
```bash
# Buka terminal baru dan jalankan:
go run main.go

# Tunggu sampai muncul:
# Server running on http://localhost:8080
```

**Option B: Using Batch File**
```bash
# Double click file:
start_server.bat
```

### 2. Run CRUD Tests

**Setelah server running**, buka terminal baru dan jalankan:

```powershell
.\test_crud_complete.ps1
```

## Expected Output

```
============================================
 POCKETBASE CRUD TEST WITH AUTHENTICATION
============================================

[STEP 1] Login as admin...
Success! Token received
User: admin@example.com (Role: admin)

[TEST 1] CREATE - Membuat mahasiswa baru...
Success!
  ID: abc123xyz
  NIM: 2025001
  Nama: Ahmad Lazim
  Jurusan: Teknik Informatika

Creating second mahasiswa...
Success! NIM: 2025002, Nama: Budi Santoso

[TEST 2] READ ALL - Get all mahasiswa...
Success! Total: 2
Data:
  - NIM: 2025001, Nama: Ahmad Lazim, Jurusan: Teknik Informatika
  - NIM: 2025002, Nama: Budi Santoso, Jurusan: Sistem Informasi

[TEST 3] READ BY ID - Get mahasiswa by ID...
Success!
  ID: abc123xyz
  NIM: 2025001
  Nama: Ahmad Lazim
  Email: ahmad.lazim@example.com

[TEST 4] UPDATE - Update mahasiswa...
Success!
  New Name: Ahmad Lazim (UPDATED)
  New Jurusan: Teknik Informatika & Komputer
  New Email: ahmad.updated@example.com

[TEST 5] DELETE - Delete mahasiswa...
Success! Mahasiswa dengan ID abc123xyz berhasil dihapus

[TEST 6] VERIFY - Check remaining data...
Total mahasiswa tersisa: 1
  - 2025002: Budi Santoso

============================================
 ALL CRUD TESTS COMPLETED!
============================================
```

## Manual Testing dengan cURL/PowerShell

### 1. Login
```powershell
$body = @{email="admin@example.com"; password="admin123"} | ConvertTo-Json
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/login" -Method POST -ContentType "application/json" -Body $body
$token = $response.token
```

### 2. CREATE Mahasiswa
```powershell
$headers = @{Authorization="Bearer $token"; "Content-Type"="application/json"}
$mahasiswa = @{nim="2025003"; nama="Test User"; jurusan="TI"; angkatan=2025; email="test@example.com"} | ConvertTo-Json
$created = Invoke-RestMethod -Uri "http://localhost:8080/api/mahasiswa" -Method POST -Headers $headers -Body $mahasiswa
```

### 3. GET All Mahasiswa
```powershell
$getHeaders = @{Authorization="Bearer $token"}
$all = Invoke-RestMethod -Uri "http://localhost:8080/api/mahasiswa" -Method GET -Headers $getHeaders
$all.items | Format-Table
```

### 4. GET By ID
```powershell
$id = $created.id
$single = Invoke-RestMethod -Uri "http://localhost:8080/api/mahasiswa/$id" -Method GET -Headers $getHeaders
```

### 5. UPDATE Mahasiswa
```powershell
$updateData = @{nim="2025003"; nama="Updated Name"; jurusan="Updated Jurusan"; angkatan=2025; email="updated@example.com"} | ConvertTo-Json
$updated = Invoke-RestMethod -Uri "http://localhost:8080/api/mahasiswa/$id" -Method PUT -Headers $headers -Body $updateData
```

### 6. DELETE Mahasiswa
```powershell
Invoke-RestMethod -Uri "http://localhost:8080/api/mahasiswa/$id" -Method DELETE -Headers $getHeaders
```

## Verify Data in PocketBase Admin

After running tests, you can verify the data in PocketBase admin panel:

1. Open: `https://pocketbase-production-521e.up.railway.app/_/`
2. Login with your admin credentials
3. Click on `mahasiswas` collection
4. You should see the test data created

## Troubleshooting

### Error: "The underlying connection was closed"
**Solution:** Server is not running. Start the server first with `go run main.go`

### Error: "Token tidak ditemukan"
**Solution:** You need to login first to get JWT token

### Error: "Invalid token"
**Solution:** Token might be expired. Login again to get new token

### Error: "Unauthorized"
**Solution:** Make sure you're using admin account for CREATE/UPDATE/DELETE operations

## Testing Other Collections

### Alumni
Replace `/mahasiswa` with `/alumni` in the URLs

### Pekerjaan Alumni
Replace `/mahasiswa` with `/pekerjaan-alumni` in the URLs

## File Structure

```
PROJECT_GITHUB/CRUD-Go-Fiber-PostgreSQL/
├── test_crud_complete.ps1    # Complete CRUD test script with auth
├── start_server.bat           # Batch file to start server
├── CRUD_TEST_GUIDE.md        # This file
└── main.go                    # Main application
```

## Success Indicators

✅ Server starts without errors
✅ Login successful and receives JWT token
✅ CREATE operation returns new mahasiswa with ID
✅ READ ALL returns list of mahasiswas
✅ READ BY ID returns specific mahasiswa
✅ UPDATE changes mahasiswa data
✅ DELETE removes mahasiswa from database
✅ Data visible in PocketBase admin panel

---

**Status:** PocketBase CRUD operations fully functional! 🎉
