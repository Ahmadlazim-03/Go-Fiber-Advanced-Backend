# Complete CRUD Test with Authentication for PocketBase
Write-Host "`n============================================" -ForegroundColor Cyan
Write-Host " POCKETBASE CRUD TEST WITH AUTHENTICATION" -ForegroundColor Cyan
Write-Host "============================================`n" -ForegroundColor Cyan

$baseUrl = "http://localhost:8080"

# STEP 1: LOGIN AS ADMIN
Write-Host "[STEP 1] Login as admin..." -ForegroundColor Green
$loginData = @{
    email = "ahmad@gmail.com"
    password = "12345678"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method POST -ContentType "application/json" -Body $loginData
    $token = $loginResponse.data.token
    Write-Host "Success! Token received" -ForegroundColor Yellow
    Write-Host "User: $($loginResponse.data.user.email) (Role: $($loginResponse.data.user.role))`n" -ForegroundColor Cyan
} catch {
    Write-Host "Login failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Error details: $($_.ErrorDetails.Message)" -ForegroundColor Red
    exit 1
}

# Prepare headers with token
$headers = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

# TEST 1: CREATE Mahasiswa
Write-Host "[TEST 1] CREATE - Membuat mahasiswa baru..." -ForegroundColor Green
$mahasiswa1 = @{
    nim = "2025001"
    nama = "Ahmad Lazim"
    jurusan = "Teknik Informatika"
    angkatan = 2025
    email = "ahmad.lazim@example.com"
} | ConvertTo-Json

try {
    $created = Invoke-RestMethod -Uri "$baseUrl/api/mahasiswa" -Method POST -Headers $headers -Body $mahasiswa1
    Write-Host "Success!" -ForegroundColor Yellow
    Write-Host "  ID: $($created.id)" -ForegroundColor White
    Write-Host "  NIM: $($created.nim)" -ForegroundColor White
    Write-Host "  Nama: $($created.nama)" -ForegroundColor White
    Write-Host "  Jurusan: $($created.jurusan)`n" -ForegroundColor White
    $createdId = $created.id
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Details: $($_.ErrorDetails.Message)`n" -ForegroundColor Red
}

# Create second mahasiswa
Write-Host "Creating second mahasiswa..." -ForegroundColor Green
$mahasiswa2 = @{
    nim = "2025002"
    nama = "Budi Santoso"
    jurusan = "Sistem Informasi"
    angkatan = 2025
    email = "budi.santoso@example.com"
} | ConvertTo-Json

try {
    $created2 = Invoke-RestMethod -Uri "$baseUrl/api/mahasiswa" -Method POST -Headers $headers -Body $mahasiswa2
    Write-Host "Success! NIM: $($created2.nim), Nama: $($created2.nama)`n" -ForegroundColor Yellow
} catch {
    Write-Host "Error: $($_.Exception.Message)`n" -ForegroundColor Red
}

# TEST 2: READ ALL
Write-Host "[TEST 2] READ ALL - Get all mahasiswa..." -ForegroundColor Green
try {
    $getHeaders = @{
        "Authorization" = "Bearer $token"
    }
    $allData = Invoke-RestMethod -Uri "$baseUrl/api/mahasiswa" -Method GET -Headers $getHeaders
    Write-Host "Success! Total: $($allData.total)" -ForegroundColor Yellow
    Write-Host "Data:" -ForegroundColor Cyan
    foreach ($item in $allData.items) {
        Write-Host "  - NIM: $($item.nim), Nama: $($item.nama), Jurusan: $($item.jurusan)" -ForegroundColor White
    }
    Write-Host ""
} catch {
    Write-Host "Error: $($_.Exception.Message)`n" -ForegroundColor Red
}

# TEST 3: READ BY ID
Write-Host "[TEST 3] READ BY ID - Get mahasiswa by ID..." -ForegroundColor Green
try {
    $single = Invoke-RestMethod -Uri "$baseUrl/api/mahasiswa/$createdId" -Method GET -Headers $getHeaders
    Write-Host "Success!" -ForegroundColor Yellow
    Write-Host "  ID: $($single.id)" -ForegroundColor White
    Write-Host "  NIM: $($single.nim)" -ForegroundColor White
    Write-Host "  Nama: $($single.nama)" -ForegroundColor White
    Write-Host "  Email: $($single.email)`n" -ForegroundColor White
} catch {
    Write-Host "Error: $($_.Exception.Message)`n" -ForegroundColor Red
}

# TEST 4: UPDATE
Write-Host "[TEST 4] UPDATE - Update mahasiswa..." -ForegroundColor Green
$updateData = @{
    nim = "2025001"
    nama = "Ahmad Lazim (UPDATED)"
    jurusan = "Teknik Informatika & Komputer"
    angkatan = 2025
    email = "ahmad.updated@example.com"
} | ConvertTo-Json

try {
    $updated = Invoke-RestMethod -Uri "$baseUrl/api/mahasiswa/$createdId" -Method PUT -Headers $headers -Body $updateData
    Write-Host "Success!" -ForegroundColor Yellow
    Write-Host "  New Name: $($updated.nama)" -ForegroundColor White
    Write-Host "  New Jurusan: $($updated.jurusan)" -ForegroundColor White
    Write-Host "  New Email: $($updated.email)`n" -ForegroundColor White
} catch {
    Write-Host "Error: $($_.Exception.Message)`n" -ForegroundColor Red
}

# TEST 5: DELETE
Write-Host "[TEST 5] DELETE - Delete mahasiswa..." -ForegroundColor Green
try {
    Invoke-RestMethod -Uri "$baseUrl/api/mahasiswa/$createdId" -Method DELETE -Headers $getHeaders
    Write-Host "Success! Mahasiswa dengan ID $createdId berhasil dihapus`n" -ForegroundColor Yellow
} catch {
    Write-Host "Error: $($_.Exception.Message)`n" -ForegroundColor Red
}

# TEST 6: VERIFY
Write-Host "[TEST 6] VERIFY - Check remaining data..." -ForegroundColor Green
try {
    $finalData = Invoke-RestMethod -Uri "$baseUrl/api/mahasiswa" -Method GET -Headers $getHeaders
    Write-Host "Total mahasiswa tersisa: $($finalData.total)" -ForegroundColor Yellow
    if ($finalData.total -gt 0) {
        foreach ($item in $finalData.items) {
            Write-Host "  - $($item.nim): $($item.nama)" -ForegroundColor White
        }
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n============================================" -ForegroundColor Cyan
Write-Host " ALL CRUD TESTS COMPLETED!" -ForegroundColor Cyan
Write-Host "============================================`n" -ForegroundColor Cyan
