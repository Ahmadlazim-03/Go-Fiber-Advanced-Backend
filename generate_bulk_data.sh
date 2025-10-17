#!/bin/bash

# Warna untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080/api"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  BULK DATA GENERATION TEST - MongoDB  ${NC}"
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
    echo $((RANDOM % 9000 + 1000))
}

# ==========================================
# 1. LOGIN WITH EXISTING ADMIN OR CREATE NEW
# ==========================================
echo -e "${YELLOW}Step 1: Login as Admin...${NC}"

# Try default admin first
ADMIN_EMAIL="admin@example.com"
ADMIN_PASSWORD="admin123"

LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$ADMIN_EMAIL\",
    \"password\": \"$ADMIN_PASSWORD\"
  }")

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4)

# If default admin doesn't work, try to create new admin
if [ -z "$TOKEN" ]; then
    echo -e "${YELLOW}Default admin not found, creating new admin...${NC}"
    
    ADMIN_EMAIL="bulkadmin_$(generate_random_string)@test.com"
    ADMIN_PASSWORD="Admin123!"
    
    REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/register" \
      -H "Content-Type: application/json" \
      -d "{
        \"nama\": \"Bulk Admin\",
        \"email\": \"$ADMIN_EMAIL\",
        \"password\": \"$ADMIN_PASSWORD\",
        \"role\": \"admin\"
      }")
    
    # Login with new admin
    LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
      -H "Content-Type: application/json" \
      -d "{
        \"email\": \"$ADMIN_EMAIL\",
        \"password\": \"$ADMIN_PASSWORD\"
      }")
    
    TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4)
fi

if [ -z "$TOKEN" ]; then
    echo -e "${RED}❌ Failed to get token${NC}"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo -e "${GREEN}✅ Admin logged in successfully${NC}"
echo -e "${BLUE}Admin Email: $ADMIN_EMAIL${NC}"
echo ""

# ==========================================
# 2. CREATE MULTIPLE USERS (50 users)
# ==========================================
echo -e "${YELLOW}Step 2: Creating 50 Users...${NC}"

for i in {1..50}; do
    RANDOM_STRING=$(generate_random_string)
    USER_EMAIL="user${i}_${RANDOM_STRING}@test.com"
    USER_NAME="user${i}_${RANDOM_STRING}"
    
    RESPONSE=$(curl -s -X POST "$BASE_URL/register" \
      -H "Content-Type: application/json" \
      -d "{
        \"username\": \"$USER_NAME\",
        \"email\": \"$USER_EMAIL\",
        \"password\": \"Password123!\",
        \"role\": \"user\"
      }")
    
    USER_ID=$(echo $RESPONSE | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    
    if [ -n "$USER_ID" ]; then
        USER_IDS+=("$USER_ID")
        echo -e "${GREEN}✅ Created user $i: $USER_NAME (ID: $USER_ID)${NC}"
    else
        echo -e "${RED}❌ Failed to create user $i${NC}"
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
    NIM="$(generate_random_number)$i"
    NAMA="Mahasiswa Test $i"
    JURUSAN=${JURUSAN_LIST[$((RANDOM % ${#JURUSAN_LIST[@]}))]}
    ANGKATAN=${ANGKATAN_LIST[$((RANDOM % ${#ANGKATAN_LIST[@]}))]}
    
    RESPONSE=$(curl -s -X POST "$BASE_URL/mahasiswa" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"nim\": \"$NIM\",
        \"nama\": \"$NAMA\",
        \"jurusan\": \"$JURUSAN\",
        \"angkatan\": $ANGKATAN
      }")
    
    MAHASISWA_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    
    if [ -n "$MAHASISWA_ID" ]; then
        MAHASISWA_IDS+=("$MAHASISWA_ID")
        echo -e "${GREEN}✅ Created mahasiswa $i: $NAMA (NIM: $NIM, ID: $MAHASISWA_ID)${NC}"
    else
        echo -e "${RED}❌ Failed to create mahasiswa $i${NC}"
    fi
done

echo -e "${BLUE}Total Mahasiswa Created: ${#MAHASISWA_IDS[@]}${NC}"
echo ""

# ==========================================
# 4. CREATE MULTIPLE ALUMNI (80 alumni)
# ==========================================
echo -e "${YELLOW}Step 4: Creating 80 Alumni...${NC}"

STATUS_LIST=("Bekerja" "Wirausaha" "Melanjutkan Studi" "Mencari Kerja")

for i in {1..80}; do
    # Pilih random user_id dari USER_IDS yang sudah dibuat
    if [ ${#USER_IDS[@]} -gt 0 ]; then
        USER_ID=${USER_IDS[$((RANDOM % ${#USER_IDS[@]}))]}
    else
        echo -e "${RED}❌ No users available for alumni $i${NC}"
        continue
    fi
    
    NIM="$(generate_random_number)$i"
    NAMA="Alumni Test $i"
    JURUSAN=${JURUSAN_LIST[$((RANDOM % ${#JURUSAN_LIST[@]}))]}
    ANGKATAN=${ANGKATAN_LIST[$((RANDOM % ${#ANGKATAN_LIST[@]}))]}
    TAHUN_LULUS=$((ANGKATAN + 4))
    STATUS=${STATUS_LIST[$((RANDOM % ${#STATUS_LIST[@]}))]}
    
    RESPONSE=$(curl -s -X POST "$BASE_URL/alumni" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"nim\": \"$NIM\",
        \"nama\": \"$NAMA\",
        \"jurusan\": \"$JURUSAN\",
        \"angkatan\": $ANGKATAN,
        \"tahun_lulus\": $TAHUN_LULUS,
        \"status\": \"$STATUS\",
        \"user_id\": \"$USER_ID\"
      }")
    
    ALUMNI_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    
    if [ -n "$ALUMNI_ID" ]; then
        ALUMNI_IDS+=("$ALUMNI_ID")
        echo -e "${GREEN}✅ Created alumni $i: $NAMA (NIM: $NIM, ID: $ALUMNI_ID)${NC}"
    else
        echo -e "${RED}❌ Failed to create alumni $i${NC}"
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

for i in {1..60}; do
    # Pilih random alumni_id dari ALUMNI_IDS yang sudah dibuat
    if [ ${#ALUMNI_IDS[@]} -gt 0 ]; then
        ALUMNI_ID=${ALUMNI_IDS[$((RANDOM % ${#ALUMNI_IDS[@]}))]}
    else
        echo -e "${RED}❌ No alumni available for pekerjaan $i${NC}"
        continue
    fi
    
    PERUSAHAAN=${PERUSAHAAN_LIST[$((RANDOM % ${#PERUSAHAAN_LIST[@]}))]}
    POSISI=${POSISI_LIST[$((RANDOM % ${#POSISI_LIST[@]}))]}
    GAJI=$((5000000 + RANDOM % 15000000))
    TAHUN_MULAI=$((2020 + RANDOM % 5))
    
    RESPONSE=$(curl -s -X POST "$BASE_URL/pekerjaan-alumni" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{
        \"alumni_id\": \"$ALUMNI_ID\",
        \"perusahaan\": \"$PERUSAHAAN\",
        \"posisi\": \"$POSISI\",
        \"gaji\": $GAJI,
        \"tahun_mulai\": $TAHUN_MULAI
      }")
    
    PEKERJAAN_ID=$(echo $RESPONSE | grep -o '"id":"[^"]*"' | head -1 | cut -d'"' -f4)
    
    if [ -n "$PEKERJAAN_ID" ]; then
        PEKERJAAN_IDS+=("$PEKERJAAN_ID")
        echo -e "${GREEN}✅ Created pekerjaan $i: $POSISI at $PERUSAHAAN (ID: $PEKERJAAN_ID)${NC}"
    else
        echo -e "${RED}❌ Failed to create pekerjaan $i${NC}"
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
echo -e "${YELLOW}Verification: Fetching All Data...${NC}"
echo ""

echo -e "${BLUE}Users in database:${NC}"
curl -s -X GET "$BASE_URL/users" -H "Authorization: Bearer $TOKEN" | grep -o '"total":[0-9]*' | head -1
echo ""

echo -e "${BLUE}Mahasiswa in database:${NC}"
curl -s -X GET "$BASE_URL/mahasiswa?page=1&limit=1" -H "Authorization: Bearer $TOKEN" | grep -o '"total":[0-9]*' | head -1
echo ""

echo -e "${BLUE}Alumni in database:${NC}"
curl -s -X GET "$BASE_URL/alumni?page=1&limit=1" -H "Authorization: Bearer $TOKEN" | grep -o '"total":[0-9]*' | head -1
echo ""

echo -e "${BLUE}Pekerjaan Alumni in database:${NC}"
curl -s -X GET "$BASE_URL/pekerjaan-alumni?page=1&limit=1" -H "Authorization: Bearer $TOKEN" | grep -o '"total":[0-9]*' | head -1
echo ""

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}   BULK DATA GENERATION COMPLETED!    ${NC}"
echo -e "${GREEN}========================================${NC}"
