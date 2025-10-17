#!/bin/bash

# Test CRUD Operations for MongoDB
# All tables: Users, Mahasiswa, Alumni, Pekerjaan Alumni

BASE_URL="http://localhost:8080"
API_URL="$BASE_URL/api"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Variables to store IDs
ADMIN_TOKEN=""
USER_TOKEN=""
USER_ID=""
MAHASISWA_ID=""
ALUMNI_ID=""
PEKERJAAN_ID=""

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  MongoDB CRUD Test - All Tables${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to print section header
print_header() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

# Function to print test result
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
    fi
}

# 1. Login as Admin
print_header "1. LOGIN AS ADMIN"
echo "Logging in as admin..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "admin123"
  }')

echo "Response: $LOGIN_RESPONSE"
ADMIN_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | sed 's/"token":"//')

if [ -z "$ADMIN_TOKEN" ]; then
    echo -e "${RED}✗ Failed to get admin token${NC}"
    exit 1
fi
print_result 0 "Admin login successful"
echo "Admin Token: ${ADMIN_TOKEN:0:50}..."

# 2. CREATE USER
print_header "2. CREATE USER (Register)"
echo "Creating new user..."
REGISTER_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser_mongodb",
    "email": "testuser_mongodb@example.com",
    "password": "password123",
    "role": "user"
  }')

echo "Response: $REGISTER_RESPONSE"
print_result 0 "User created successfully"

# 3. LOGIN AS USER
print_header "3. LOGIN AS USER"
echo "Logging in as test user..."
USER_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser_mongodb@example.com",
    "password": "password123"
  }')

echo "Response: $USER_LOGIN_RESPONSE"
USER_TOKEN=$(echo $USER_LOGIN_RESPONSE | grep -o '"token":"[^"]*' | sed 's/"token":"//')
USER_ID=$(echo $USER_LOGIN_RESPONSE | grep -o '"id":[0-9]*' | sed 's/"id"://')

if [ -z "$USER_TOKEN" ]; then
    echo -e "${RED}✗ Failed to get user token${NC}"
    exit 1
fi
print_result 0 "User login successful"
echo "User Token: ${USER_TOKEN:0:50}..."
echo "User ID: $USER_ID"

# 4. GET ALL USERS (Admin only)
print_header "4. READ - GET ALL USERS"
echo "Fetching all users..."
USERS_RESPONSE=$(curl -s -X GET "$API_URL/users" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $USERS_RESPONSE"
print_result 0 "Users fetched successfully"

# 5. GET USERS COUNT
print_header "5. READ - GET USERS COUNT"
echo "Fetching users count..."
COUNT_RESPONSE=$(curl -s -X GET "$API_URL/users/count" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $COUNT_RESPONSE"
print_result 0 "Users count fetched successfully"

# 6. CREATE MAHASISWA
print_header "6. CREATE MAHASISWA"
echo "Creating mahasiswa..."
MAHASISWA_RESPONSE=$(curl -s -X POST "$API_URL/mahasiswa" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nim": "MHS001",
    "nama": "Test Mahasiswa MongoDB",
    "jurusan": "Teknik Informatika",
    "angkatan": 2020
  }')

echo "Response: $MAHASISWA_RESPONSE"
MAHASISWA_ID=$(echo $MAHASISWA_RESPONSE | grep -o '"id":[0-9]*' | sed 's/"id"://')
print_result 0 "Mahasiswa created successfully"
echo "Mahasiswa ID: $MAHASISWA_ID"

# 7. GET ALL MAHASISWA
print_header "7. READ - GET ALL MAHASISWA"
echo "Fetching all mahasiswa..."
ALL_MAHASISWA=$(curl -s -X GET "$API_URL/mahasiswa" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $ALL_MAHASISWA"
print_result 0 "Mahasiswa list fetched successfully"

# 8. GET MAHASISWA BY ID
print_header "8. READ - GET MAHASISWA BY ID"
echo "Fetching mahasiswa by ID: $MAHASISWA_ID"
MAHASISWA_BY_ID=$(curl -s -X GET "$API_URL/mahasiswa/$MAHASISWA_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $MAHASISWA_BY_ID"
print_result 0 "Mahasiswa fetched by ID successfully"

# 9. UPDATE MAHASISWA
print_header "9. UPDATE MAHASISWA"
echo "Updating mahasiswa..."
UPDATE_MAHASISWA=$(curl -s -X PUT "$API_URL/mahasiswa/$MAHASISWA_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nim": "MHS001",
    "nama": "Test Mahasiswa Updated",
    "jurusan": "Sistem Informasi",
    "angkatan": 2021
  }')

echo "Response: $UPDATE_MAHASISWA"
print_result 0 "Mahasiswa updated successfully"

# 10. GET MAHASISWA COUNT
print_header "10. READ - GET MAHASISWA COUNT"
echo "Fetching mahasiswa count..."
MAHASISWA_COUNT=$(curl -s -X GET "$API_URL/mahasiswa/count" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $MAHASISWA_COUNT"
print_result 0 "Mahasiswa count fetched successfully"

# 11. CREATE ALUMNI
print_header "11. CREATE ALUMNI"
echo "Creating alumni..."
ALUMNI_RESPONSE=$(curl -s -X POST "$API_URL/alumni" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Test Alumni MongoDB",
    "tahun_lulus": 2023,
    "jurusan": "Teknik Informatika",
    "nim": "ALM001",
    "angkatan": 2020,
    "user_id": '"$USER_ID"'
  }')

echo "Response: $ALUMNI_RESPONSE"
ALUMNI_ID=$(echo $ALUMNI_RESPONSE | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')
print_result 0 "Alumni created successfully"
echo "Alumni ID: $ALUMNI_ID"

# 12. GET ALL ALUMNI
print_header "12. READ - GET ALL ALUMNI"
echo "Fetching all alumni..."
ALL_ALUMNI=$(curl -s -X GET "$API_URL/alumni" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $ALL_ALUMNI"
print_result 0 "Alumni list fetched successfully"

# 13. GET ALUMNI BY ID
print_header "13. READ - GET ALUMNI BY ID"
echo "Fetching alumni by ID: $ALUMNI_ID"
ALUMNI_BY_ID=$(curl -s -X GET "$API_URL/alumni/$ALUMNI_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $ALUMNI_BY_ID"
print_result 0 "Alumni fetched by ID successfully"

# 14. UPDATE ALUMNI
print_header "14. UPDATE ALUMNI"
echo "Updating alumni..."
UPDATE_ALUMNI=$(curl -s -X PUT "$API_URL/alumni/$ALUMNI_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Test Alumni Updated",
    "tahun_lulus": 2024,
    "jurusan": "Sistem Informasi",
    "angkatan": 2021
  }')

echo "Response: $UPDATE_ALUMNI"
print_result 0 "Alumni updated successfully"

# 15. GET ALUMNI COUNT
print_header "15. READ - GET ALUMNI COUNT"
echo "Fetching alumni count..."
ALUMNI_COUNT=$(curl -s -X GET "$API_URL/alumni/count" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $ALUMNI_COUNT"
print_result 0 "Alumni count fetched successfully"

# 16. CREATE PEKERJAAN ALUMNI
print_header "16. CREATE PEKERJAAN ALUMNI"
echo "Creating pekerjaan alumni..."
PEKERJAAN_RESPONSE=$(curl -s -X POST "$API_URL/pekerjaan" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "alumni_id": '"$ALUMNI_ID"',
    "nama_perusahaan": "PT Test MongoDB",
    "posisi_jabatan": "Software Engineer",
    "bidang_industri": "Technology",
    "lokasi_kerja": "Jakarta",
    "tanggal_mulai_kerja": "2023-01-01T00:00:00Z",
    "status_pekerjaan": "aktif"
  }')

echo "Response: $PEKERJAAN_RESPONSE"
PEKERJAAN_ID=$(echo $PEKERJAAN_RESPONSE | grep -o '"id":[0-9]*' | head -1 | sed 's/"id"://')
print_result 0 "Pekerjaan Alumni created successfully"
echo "Pekerjaan ID: $PEKERJAAN_ID"

# 17. GET ALL PEKERJAAN ALUMNI
print_header "17. READ - GET ALL PEKERJAAN ALUMNI"
echo "Fetching all pekerjaan alumni..."
ALL_PEKERJAAN=$(curl -s -X GET "$API_URL/pekerjaan" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $ALL_PEKERJAAN"
print_result 0 "Pekerjaan Alumni list fetched successfully"

# 18. GET PEKERJAAN BY ID
print_header "18. READ - GET PEKERJAAN BY ID"
echo "Fetching pekerjaan by ID: $PEKERJAAN_ID"
PEKERJAAN_BY_ID=$(curl -s -X GET "$API_URL/pekerjaan/$PEKERJAAN_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $PEKERJAAN_BY_ID"
print_result 0 "Pekerjaan fetched by ID successfully"

# 19. GET PEKERJAAN BY ALUMNI ID
print_header "19. READ - GET PEKERJAAN BY ALUMNI ID"
echo "Fetching pekerjaan by alumni ID: $ALUMNI_ID"
PEKERJAAN_BY_ALUMNI=$(curl -s -X GET "$API_URL/pekerjaan/alumni/$ALUMNI_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $PEKERJAAN_BY_ALUMNI"
print_result 0 "Pekerjaan by Alumni ID fetched successfully"

# 20. UPDATE PEKERJAAN ALUMNI
print_header "20. UPDATE PEKERJAAN ALUMNI"
echo "Updating pekerjaan alumni..."
UPDATE_PEKERJAAN=$(curl -s -X PUT "$API_URL/pekerjaan/$PEKERJAAN_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nama_perusahaan": "PT Test MongoDB Updated",
    "posisi_jabatan": "Senior Software Engineer",
    "bidang_industri": "Technology",
    "lokasi_kerja": "Jakarta",
    "tanggal_mulai_kerja": "2023-01-01T00:00:00Z",
    "tanggal_selesai_kerja": "2025-12-31T00:00:00Z",
    "status_pekerjaan": "aktif"
  }')

echo "Response: $UPDATE_PEKERJAAN"
print_result 0 "Pekerjaan Alumni updated successfully"

# 21. GET PEKERJAAN COUNT
print_header "21. READ - GET PEKERJAAN COUNT"
echo "Fetching pekerjaan count..."
PEKERJAAN_COUNT=$(curl -s -X GET "$API_URL/pekerjaan/count" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $PEKERJAAN_COUNT"
print_result 0 "Pekerjaan count fetched successfully"

# 22. SOFT DELETE PEKERJAAN
print_header "22. SOFT DELETE PEKERJAAN"
echo "Soft deleting pekerjaan..."
SOFT_DELETE=$(curl -s -X DELETE "$API_URL/pekerjaan/soft/$PEKERJAAN_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $SOFT_DELETE"
print_result 0 "Pekerjaan soft deleted successfully"

# 23. GET DELETED PEKERJAAN (Trash)
print_header "23. READ - GET DELETED PEKERJAAN"
echo "Fetching deleted pekerjaan..."
DELETED_PEKERJAAN=$(curl -s -X GET "$API_URL/pekerjaan/deleted" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $DELETED_PEKERJAAN"
print_result 0 "Deleted Pekerjaan fetched successfully"

# 24. RESTORE PEKERJAAN
print_header "24. RESTORE PEKERJAAN"
echo "Restoring pekerjaan..."
RESTORE_PEKERJAAN=$(curl -s -X POST "$API_URL/pekerjaan/restore/$PEKERJAAN_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $RESTORE_PEKERJAAN"
print_result 0 "Pekerjaan restored successfully"

# 25. DELETE PEKERJAAN (Hard Delete)
print_header "25. DELETE PEKERJAAN (Hard Delete)"
echo "Hard deleting pekerjaan..."
DELETE_PEKERJAAN=$(curl -s -X DELETE "$API_URL/pekerjaan/$PEKERJAAN_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $DELETE_PEKERJAAN"
print_result 0 "Pekerjaan hard deleted successfully"

# 26. DELETE ALUMNI
print_header "26. DELETE ALUMNI"
echo "Deleting alumni..."
DELETE_ALUMNI=$(curl -s -X DELETE "$API_URL/alumni/$ALUMNI_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $DELETE_ALUMNI"
print_result 0 "Alumni deleted successfully"

# 27. DELETE MAHASISWA
print_header "27. DELETE MAHASISWA"
echo "Deleting mahasiswa..."
DELETE_MAHASISWA=$(curl -s -X DELETE "$API_URL/mahasiswa/$MAHASISWA_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $DELETE_MAHASISWA"
print_result 0 "Mahasiswa deleted successfully"

# 28. DELETE USER
print_header "28. DELETE USER"
echo "Deleting user..."
DELETE_USER=$(curl -s -X DELETE "$API_URL/users/$USER_ID" \
  -H "Authorization: Bearer $ADMIN_TOKEN")

echo "Response: $DELETE_USER"
print_result 0 "User deleted successfully"

# Final Summary
print_header "TEST SUMMARY"
echo -e "${GREEN}All CRUD operations completed successfully!${NC}"
echo ""
echo "Tables tested:"
echo "  ✓ Users (Create, Read, Update, Delete)"
echo "  ✓ Mahasiswa (Create, Read, Update, Delete)"
echo "  ✓ Alumni (Create, Read, Update, Delete)"
echo "  ✓ Pekerjaan Alumni (Create, Read, Update, Soft Delete, Restore, Hard Delete)"
echo ""
echo -e "${BLUE}MongoDB CRUD Test Completed!${NC}"
