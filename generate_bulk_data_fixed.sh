#!/bin/bash

# Warna untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080/api"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  BULK DATA GENERATION - MongoDB Fixed ${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Array untuk menyimpan ID yang dibuat
declare -a USER_IDS=()
declare -a MAHASISWA_IDS=()
declare -a ALUMNI_IDS=()
declare -a PEKERJAAN_IDS=()

# Function untuk generate random string
generate_random_string() {
    cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1
}

# Function untuk generate random number
generate_random_number() {
    echo $((RANDOM % 900000 + 100000))
}

# ==========================================
# 1. LOGIN WITH ADMIN
# ==========================================
echo -e "${YELLOW}Step 1: Login as Admin...${NC}"

ADMIN_EMAIL="admin@example.com"
ADMIN_PASSWORD="admin123"

LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$ADMIN_EMAIL\",
    \"password\": \"$ADMIN_PASSWORD\"
  }")

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ Failed to get token${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Admin logged in successfully${NC}"
echo ""

# ==========================================
# 2. CREATE MULTIPLE USERS (50 users)
# ==========================================
echo -e "${YELLOW}Step 2: Creating 50 Users...${NC}"

for i in {1..50}; do
    RANDOM_STRING=$(generate_random_string)
    USER_USERNAME="user_${i}_${RANDOM_STRING}"
    USER_EMAIL="user${i}_${RANDOM_STRING}@test.com"
    
    RESPONSE=$(curl -s -X POST "$BASE_URL/register" \
      -H "Content-Type: application/json" \
      -d "{
        \"username\": \"$USER_USERNAME\",
        \"email\": \"$USER_EMAIL\",
        \"password\": \"Password123!\",
        \"role\": \"user\"
      }")
    
    USER_ID=$(echo $RESPONSE | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    
    if [ -n "$USER_ID" ]; then
        USER_IDS+=("$USER_ID")
        echo -e "${GREEN}✅ User $i: $USER_USERNAME (ID: $USER_ID)${NC}"
    else
        echo -e "${RED}❌ Failed user $i: $(echo $RESPONSE | head -c 100)${NC}"
    fi
done

echo -e "${BLUE}Total Users Created: ${#USER_IDS[@]}${NC}"
echo ""

# ==========================================
# 3. CREATE MULTIPLE MAHASISWA (100 mahasiswa)
# ==========================================
echo -e "${YELLOW}Step 3: Creating 100 Mahasiswa...${NC}"

JURUSAN_LIST=("Teknik Informatika" "Sistem Informasi" "Teknik Komputer" "Teknik Elektro" "Manajemen Informatika")
ANGKATAN_LIST=(2019 2020 2021 2022 2023 2024)

for i in {1..100}; do
    RANDOM_STRING=$(generate_random_string)
    NIM="M$(generate_random_number)"
    NAMA="Mahasiswa Test $i"
    EMAIL="mhs${i}_${RANDOM_STRING}@student.ac.id"
    JURUSAN=${JURUSAN_LIST[$((RANDOM % ${#JURUSAN_LIST[@]}))]}
    ANGKATAN=${ANGKATAN_LIST[$((RANDOM % ${#ANGKATAN_LIST[@]}))]}
    
    RESPONSE=$(curl -s -X POST "$BASE_URL/mahasiswa" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"nim\": \"$NIM\",
        \"nama\": \"$NAMA\",
        \"jurusan\": \"$JURUSAN\",
        \"angkatan\": $ANGKATAN,
        \"email\": \"$EMAIL\"
      }")
    
    MAHASISWA_ID=$(echo $RESPONSE | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    
    if [ -n "$MAHASISWA_ID" ]; then
        MAHASISWA_IDS+=("$MAHASISWA_ID")
        echo -e "${GREEN}✅ Mahasiswa $i: $NAMA - $NIM (ID: $MAHASISWA_ID)${NC}"
    else
        echo -e "${RED}❌ Failed mahasiswa $i${NC}"
    fi
done

echo -e "${BLUE}Total Mahasiswa Created: ${#MAHASISWA_IDS[@]}${NC}"
echo ""

# ==========================================
# 4. CREATE MULTIPLE ALUMNI (80 alumni)
# ==========================================
echo -e "${YELLOW}Step 4: Creating 80 Alumni...${NC}"

for i in {1..80}; do
    # Pilih random user_id dari USER_IDS yang sudah dibuat
    if [ ${#USER_IDS[@]} -gt 0 ]; then
        USER_ID=${USER_IDS[$((RANDOM % ${#USER_IDS[@]}))]}
    else
        echo -e "${RED}❌ No users available${NC}"
        break
    fi
    
    RANDOM_STRING=$(generate_random_string)
    NIM="A$(generate_random_number)"
    NAMA="Alumni Test $i"
    JURUSAN=${JURUSAN_LIST[$((RANDOM % ${#JURUSAN_LIST[@]}))]}
    ANGKATAN=${ANGKATAN_LIST[$((RANDOM % 3))]}  # Angkatan lama (2019-2021)
    TAHUN_LULUS=$((ANGKATAN + 4))
    
    RESPONSE=$(curl -s -X POST "$BASE_URL/alumni" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"nim\": \"$NIM\",
        \"nama\": \"$NAMA\",
        \"jurusan\": \"$JURUSAN\",
        \"angkatan\": $ANGKATAN,
        \"tahun_lulus\": $TAHUN_LULUS,
        \"user_id\": $USER_ID,
        \"no_telepon\": \"0812345678$i\",
        \"alamat\": \"Jalan Test No. $i\"
      }")
    
    ALUMNI_ID=$(echo $RESPONSE | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    
    if [ -n "$ALUMNI_ID" ]; then
        ALUMNI_IDS+=("$ALUMNI_ID")
        echo -e "${GREEN}✅ Alumni $i: $NAMA - $NIM (ID: $ALUMNI_ID)${NC}"
    else
        echo -e "${RED}❌ Failed alumni $i${NC}"
    fi
done

echo -e "${BLUE}Total Alumni Created: ${#ALUMNI_IDS[@]}${NC}"
echo ""

# ==========================================
# 5. CREATE MULTIPLE PEKERJAAN ALUMNI (60 pekerjaan)
# ==========================================
echo -e "${YELLOW}Step 5: Creating 60 Pekerjaan Alumni...${NC}"

PERUSAHAAN_LIST=("PT Tech Indonesia" "CV Digital Nusantara" "PT Maju Jaya" "Startup Teknologi" "Bank Digital" "Konsultan IT" "E-Commerce Indo" "Fintech Solutions")
POSISI_LIST=("Software Engineer" "Frontend Developer" "Backend Developer" "Full Stack Developer" "DevOps Engineer" "Data Analyst" "Project Manager" "UI/UX Designer")
BIDANG_LIST=("Teknologi Informasi" "Perbankan" "E-Commerce" "Konsultansi" "Startup" "Pendidikan" "Kesehatan" "Fintech")
LOKASI_LIST=("Jakarta" "Bandung" "Surabaya" "Yogyakarta" "Semarang" "Medan" "Bali" "Makassar")
GAJI_RANGE_LIST=("5-7 juta" "7-10 juta" "10-15 juta" "15-20 juta" "20-30 juta" ">30 juta")

for i in {1..60}; do
    # Pilih random alumni_id dari ALUMNI_IDS yang sudah dibuat
    if [ ${#ALUMNI_IDS[@]} -gt 0 ]; then
        ALUMNI_ID=${ALUMNI_IDS[$((RANDOM % ${#ALUMNI_IDS[@]}))]}
    else
        echo -e "${RED}❌ No alumni available${NC}"
        break
    fi
    
    PERUSAHAAN=${PERUSAHAAN_LIST[$((RANDOM % ${#PERUSAHAAN_LIST[@]}))]}
    POSISI=${POSISI_LIST[$((RANDOM % ${#POSISI_LIST[@]}))]}
    BIDANG=${BIDANG_LIST[$((RANDOM % ${#BIDANG_LIST[@]}))]}
    LOKASI=${LOKASI_LIST[$((RANDOM % ${#LOKASI_LIST[@]}))]}
    GAJI=${GAJI_RANGE_LIST[$((RANDOM % ${#GAJI_RANGE_LIST[@]}))]}
    TAHUN=$((2020 + RANDOM % 5))
    BULAN=$((1 + RANDOM % 12))
    TANGGAL_MULAI="${TAHUN}-$(printf "%02d" $BULAN)-01T00:00:00Z"
    
    RESPONSE=$(curl -s -X POST "$BASE_URL/pekerjaan" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"alumni_id\": $ALUMNI_ID,
        \"nama_perusahaan\": \"$PERUSAHAAN\",
        \"posisi_jabatan\": \"$POSISI\",
        \"bidang_industri\": \"$BIDANG\",
        \"lokasi_kerja\": \"$LOKASI\",
        \"gaji_range\": \"$GAJI\",
        \"tanggal_mulai_kerja\": \"$TANGGAL_MULAI\",
        \"status_pekerjaan\": \"aktif\",
        \"deskripsi_pekerjaan\": \"Bertanggung jawab sebagai $POSISI di $PERUSAHAAN\"
      }")
    
    PEKERJAAN_ID=$(echo $RESPONSE | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    
    if [ -n "$PEKERJAAN_ID" ]; then
        PEKERJAAN_IDS+=("$PEKERJAAN_ID")
        echo -e "${GREEN}✅ Pekerjaan $i: $POSISI at $PERUSAHAAN (ID: $PEKERJAAN_ID)${NC}"
    else
        echo -e "${RED}❌ Failed pekerjaan $i${NC}"
    fi
done

echo -e "${BLUE}Total Pekerjaan Created: ${#PEKERJAAN_IDS[@]}${NC}"
echo ""

# ==========================================
# SUMMARY
# ==========================================
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}           GENERATION SUMMARY          ${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✅ Users Created: ${#USER_IDS[@]}${NC}"
echo -e "${GREEN}✅ Mahasiswa Created: ${#MAHASISWA_IDS[@]}${NC}"
echo -e "${GREEN}✅ Alumni Created: ${#ALUMNI_IDS[@]}${NC}"
echo -e "${GREEN}✅ Pekerjaan Alumni Created: ${#PEKERJAAN_IDS[@]}${NC}"
echo ""
echo -e "${BLUE}Total Records: $((${#USER_IDS[@]} + ${#MAHASISWA_IDS[@]} + ${#ALUMNI_IDS[@]} + ${#PEKERJAAN_IDS[@]}))${NC}"
echo ""

# ==========================================
# VERIFICATION - GET ALL DATA
# ==========================================
echo -e "${YELLOW}Verification: Fetching All Data Count...${NC}"
echo ""

USERS_COUNT=$(curl -s -X GET "$BASE_URL/users" -H "Authorization: Bearer $TOKEN" | grep -o '"total":[0-9]*' | head -1 | grep -o '[0-9]*')
echo -e "${BLUE}Users in database: ${GREEN}$USERS_COUNT${NC}"

MAHASISWA_COUNT=$(curl -s -X GET "$BASE_URL/mahasiswa?page=1&limit=1" -H "Authorization: Bearer $TOKEN" | grep -o '"total":[0-9]*' | head -1 | grep -o '[0-9]*')
echo -e "${BLUE}Mahasiswa in database: ${GREEN}$MAHASISWA_COUNT${NC}"

ALUMNI_COUNT=$(curl -s -X GET "$BASE_URL/alumni?page=1&limit=1" -H "Authorization: Bearer $TOKEN" | grep -o '"total":[0-9]*' | head -1 | grep -o '[0-9]*')
echo -e "${BLUE}Alumni in database: ${GREEN}$ALUMNI_COUNT${NC}"

PEKERJAAN_COUNT=$(curl -s -X GET "$BASE_URL/pekerjaan?page=1&limit=1" -H "Authorization: Bearer $TOKEN" | grep -o '"total":[0-9]*' | head -1 | grep -o '[0-9]*')
echo -e "${BLUE}Pekerjaan Alumni in database: ${GREEN}$PEKERJAAN_COUNT${NC}"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   BULK DATA GENERATION COMPLETED!    ${NC}"
echo -e "${GREEN}========================================${NC}"
