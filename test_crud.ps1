# CRUD Test Script for PocketBase - Simplified
Write-Host "`n========= POCKETBASE CRUD TEST - MAHASISWA =========`n" -ForegroundColor Cyan

$baseUrl = "http://localhost:8080/api"
$ErrorActionPreference = "Stop"

# TEST 1: CREATE
Write-Host "[TEST 1] CREATE - Membuat mahasiswa baru..." -ForegroundColor Green
$mahasiswa1 = @{
    nim = "2025001"
    nama = "Ahmad Lazim"
    jurusan = "Teknik Informatika"
    angkatan = 2025
    email = "ahmad.lazim@example.com"
} | ConvertTo-Json

try {
    $result1 = Invoke-RestMethod -Uri "$baseUrl/mahasiswas" -Method POST -ContentType "application/json" -Body $mahasiswa1
    Write-Host "Success! ID: $($result1.id), Nama: $($result1.nama)" -ForegroundColor Yellow
    $createdId = $result1.id
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

# TEST 2: READ ALL
Write-Host "`n[TEST 2] READ ALL - Mengambil semua mahasiswa..." -ForegroundColor Green
try {
    $allMahasiswa = Invoke-RestMethod -Uri "$baseUrl/mahasiswas" -Method GET
    Write-Host "Total data: $($allMahasiswa.total)" -ForegroundColor Yellow
    foreach ($mhs in $allMahasiswa.items) {
        Write-Host "  - NIM: $($mhs.nim), Nama: $($mhs.nama)" -ForegroundColor White
    }
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

# TEST 3: READ BY ID
Write-Host "`n[TEST 3] READ BY ID - Mengambil mahasiswa ID: $createdId..." -ForegroundColor Green
try {
    $single = Invoke-RestMethod -Uri "$baseUrl/mahasiswas/$createdId" -Method GET
    Write-Host "Found: $($single.nama) - $($single.email)" -ForegroundColor Yellow
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

# TEST 4: UPDATE
Write-Host "`n[TEST 4] UPDATE - Mengupdate data mahasiswa..." -ForegroundColor Green
$updateData = @{
    nim = "2025001"
    nama = "Ahmad Lazim (Updated)"
    jurusan = "Teknik Informatika & Komputer"
    angkatan = 2025
    email = "ahmad.updated@example.com"
} | ConvertTo-Json

try {
    $updated = Invoke-RestMethod -Uri "$baseUrl/mahasiswas/$createdId" -Method PUT -ContentType "application/json" -Body $updateData
    Write-Host "Updated! New name: $($updated.nama)" -ForegroundColor Yellow
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

# TEST 5: DELETE
Write-Host "`n[TEST 5] DELETE - Menghapus mahasiswa..." -ForegroundColor Green
try {
    Invoke-RestMethod -Uri "$baseUrl/mahasiswas/$createdId" -Method DELETE
    Write-Host "Deleted! ID: $createdId" -ForegroundColor Yellow
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

# TEST 6: VERIFY DELETE
Write-Host "`n[TEST 6] VERIFY - Cek data setelah delete..." -ForegroundColor Green
try {
    $finalData = Invoke-RestMethod -Uri "$baseUrl/mahasiswas" -Method GET
    Write-Host "Total mahasiswa tersisa: $($finalData.total)" -ForegroundColor Yellow
} catch {
    Write-Host "Error: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n========= CRUD TEST SELESAI =========`n" -ForegroundColor Cyan
