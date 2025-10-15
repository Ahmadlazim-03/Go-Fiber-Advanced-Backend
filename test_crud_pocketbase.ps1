# ========================================
# CRUD Test Script for PocketBase
# ========================================

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  POCKETBASE CRUD TEST - MAHASISWA" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

$baseUrl = "http://localhost:8080/api"

# ========================================
# TEST 1: CREATE - Buat Mahasiswa Baru
# ========================================
Write-Host "`n[TEST 1] CREATE - Membuat mahasiswa baru..." -ForegroundColor Green

$mahasiswa1 = @{
    nim = "2025001"
    nama = "Ahmad Lazim"
    jurusan = "Teknik Informatika"
    angkatan = 2025
    email = "ahmad.lazim@example.com"
} | ConvertTo-Json

try {
    $result1 = Invoke-RestMethod -Uri "$baseUrl/mahasiswas" -Method POST -ContentType "application/json" -Body $mahasiswa1
    Write-Host "✓ Mahasiswa berhasil dibuat!" -ForegroundColor Green
    Write-Host "  ID: $($result1.id)" -ForegroundColor Yellow
    Write-Host "  NIM: $($result1.nim)" -ForegroundColor Yellow
    Write-Host "  Nama: $($result1.nama)" -ForegroundColor Yellow
    $createdId = $result1.id
} catch {
    Write-Host "✗ Error: $_" -ForegroundColor Red
    exit 1
}

# Buat mahasiswa kedua
$mahasiswa2 = @{
    nim = "2025002"
    nama = "Budi Santoso"
    jurusan = "Sistem Informasi"
    angkatan = 2025
    email = "budi.santoso@example.com"
} | ConvertTo-Json

try {
    $result2 = Invoke-RestMethod -Uri "$baseUrl/mahasiswas" -Method POST -ContentType "application/json" -Body $mahasiswa2
    Write-Host "✓ Mahasiswa kedua berhasil dibuat!" -ForegroundColor Green
    Write-Host "  ID: $($result2.id)" -ForegroundColor Yellow
    Write-Host "  NIM: $($result2.nim)" -ForegroundColor Yellow
    Write-Host "  Nama: $($result2.nama)" -ForegroundColor Yellow
} catch {
    Write-Host "✗ Error: $_" -ForegroundColor Red
}

# ========================================
# TEST 2: READ ALL - Ambil Semua Mahasiswa
# ========================================
Write-Host "`n[TEST 2] READ ALL - Mengambil semua mahasiswa..." -ForegroundColor Green

try {
    $allMahasiswa = Invoke-RestMethod -Uri "$baseUrl/mahasiswas" -Method GET
    Write-Host "✓ Berhasil mengambil data mahasiswa!" -ForegroundColor Green
    Write-Host "  Total data: $($allMahasiswa.total)" -ForegroundColor Yellow
    Write-Host "`n  Daftar Mahasiswa:" -ForegroundColor Cyan
    foreach ($mhs in $allMahasiswa.items) {
        $nimValue = $mhs.nim
        $namaValue = $mhs.nama
        $jurusanValue = $mhs.jurusan
        Write-Host "    - $nimValue - $namaValue - $jurusanValue" -ForegroundColor White
    }
} catch {
    Write-Host "✗ Error: $_" -ForegroundColor Red
}

# ========================================
# TEST 3: READ BY ID - Ambil Mahasiswa Tertentu
# ========================================
Write-Host "`n[TEST 3] READ BY ID - Mengambil mahasiswa dengan ID: $createdId..." -ForegroundColor Green

try {
    $singleMahasiswa = Invoke-RestMethod -Uri "$baseUrl/mahasiswas/$createdId" -Method GET
    Write-Host "✓ Berhasil mengambil data mahasiswa!" -ForegroundColor Green
    Write-Host "  ID: $($singleMahasiswa.id)" -ForegroundColor Yellow
    Write-Host "  NIM: $($singleMahasiswa.nim)" -ForegroundColor Yellow
    Write-Host "  Nama: $($singleMahasiswa.nama)" -ForegroundColor Yellow
    Write-Host "  Jurusan: $($singleMahasiswa.jurusan)" -ForegroundColor Yellow
    Write-Host "  Angkatan: $($singleMahasiswa.angkatan)" -ForegroundColor Yellow
    Write-Host "  Email: $($singleMahasiswa.email)" -ForegroundColor Yellow
} catch {
    Write-Host "✗ Error: $_" -ForegroundColor Red
}

# ========================================
# TEST 4: UPDATE - Update Data Mahasiswa
# ========================================
Write-Host "`n[TEST 4] UPDATE - Mengupdate data mahasiswa..." -ForegroundColor Green

$updateData = @{
    nim = "2025001"
    nama = "Ahmad Lazim (Updated)"
    jurusan = "Teknik Informatika & Komputer"
    angkatan = 2025
    email = "ahmad.lazim.updated@example.com"
} | ConvertTo-Json

try {
    $updated = Invoke-RestMethod -Uri "$baseUrl/mahasiswas/$createdId" -Method PUT -ContentType "application/json" -Body $updateData
    Write-Host "✓ Data mahasiswa berhasil diupdate!" -ForegroundColor Green
    Write-Host "  Nama baru: $($updated.nama)" -ForegroundColor Yellow
    Write-Host "  Jurusan baru: $($updated.jurusan)" -ForegroundColor Yellow
    Write-Host "  Email baru: $($updated.email)" -ForegroundColor Yellow
} catch {
    Write-Host "✗ Error: $_" -ForegroundColor Red
}

# ========================================
# TEST 5: DELETE - Hapus Mahasiswa
# ========================================
Write-Host "`n[TEST 5] DELETE - Menghapus mahasiswa..." -ForegroundColor Green

try {
    Invoke-RestMethod -Uri "$baseUrl/mahasiswas/$createdId" -Method DELETE
    Write-Host "✓ Mahasiswa berhasil dihapus!" -ForegroundColor Green
    Write-Host "  ID yang dihapus: $createdId" -ForegroundColor Yellow
} catch {
    Write-Host "✗ Error: $_" -ForegroundColor Red
}

# Verifikasi penghapusan
Write-Host "`n  Verifikasi: Cek apakah data sudah terhapus..." -ForegroundColor Cyan
try {
    Invoke-RestMethod -Uri "$baseUrl/mahasiswas/$createdId" -Method GET
    Write-Host "  ⚠️  Data masih ada (mungkin soft delete)" -ForegroundColor Yellow
} catch {
    Write-Host "  ✓ Data berhasil terhapus!" -ForegroundColor Green
}

# ========================================
# TEST 6: READ ALL AGAIN - Cek Data Akhir
# ========================================
Write-Host "`n[TEST 6] READ ALL AGAIN - Data mahasiswa setelah operasi CRUD..." -ForegroundColor Green

try {
    $finalData = Invoke-RestMethod -Uri "$baseUrl/mahasiswas" -Method GET
    Write-Host "✓ Total mahasiswa tersisa: $($finalData.total)" -ForegroundColor Green
    if ($finalData.total -gt 0) {
        Write-Host "`n  Daftar Mahasiswa:" -ForegroundColor Cyan
        foreach ($mhs in $finalData.items) {
            $nimValue = $mhs.nim
            $namaValue = $mhs.nama
            $jurusanValue = $mhs.jurusan
            Write-Host "    - $nimValue - $namaValue - $jurusanValue" -ForegroundColor White
        }
    }
} catch {
    Write-Host "✗ Error: $_" -ForegroundColor Red
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "  CRUD TEST SELESAI!" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan
