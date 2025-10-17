#!/bin/bash

# Warna untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080/api"
SUCCESS_COUNT=0
FAIL_COUNT=0

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   COMPREHENSIVE ROUTES TEST - MongoDB ${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to test endpoint
test_endpoint() {
    local METHOD=$1
    local ENDPOINT=$2
    local DESCRIPTION=$3
    local DATA=$4
    local TOKEN=$5
    
    echo -e "${CYAN}Testing: ${DESCRIPTION}${NC}"
    
    if [ -n "$TOKEN" ]; then
        if [ "$METHOD" == "GET" ]; then
            RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X $METHOD "$BASE_URL$ENDPOINT" -H "Authorization: Bearer $TOKEN")
        else
            RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X $METHOD "$BASE_URL$ENDPOINT" -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$DATA")
        fi
    else
        if [ "$METHOD" == "GET" ]; then
            RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X $METHOD "$BASE_URL$ENDPOINT")
        else
            RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X $METHOD "$BASE_URL$ENDPOINT" -H "Content-Type: application/json" -d "$DATA")
        fi
    fi
    
    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d':' -f2)
    BODY=$(echo "$RESPONSE" | sed '/HTTP_CODE/d')
    
    if [ "$HTTP_CODE" -ge 200 ] && [ "$HTTP_CODE" -lt 300 ]; then
        echo -e "${GREEN}✅ SUCCESS (HTTP $HTTP_CODE)${NC}"
        echo -e "${GREEN}Response: $(echo $BODY | head -c 150)...${NC}"
        ((SUCCESS_COUNT++))
    else
        echo -e "${RED}❌ FAILED (HTTP $HTTP_CODE)${NC}"
        echo -e "${RED}Response: $(echo $BODY | head -c 150)${NC}"
        ((FAIL_COUNT++))
    fi
    echo ""
}

# ==========================================
# 1. AUTHENTICATION ROUTES
# ==========================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}1. AUTHENTICATION ROUTES${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# Register new user
RANDOM_STR=$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1)
test_endpoint "POST" "/register" "Register New User" "{\"username\":\"testroute_${RANDOM_STR}\",\"email\":\"testroute_${RANDOM_STR}@test.com\",\"password\":\"Test123!\",\"role\":\"user\"}" ""

# Login with admin
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}')

ADMIN_TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$ADMIN_TOKEN" ]; then
    echo -e "${RED}❌ Failed to get admin token. Exiting...${NC}"
    exit 1
fi

test_endpoint "POST" "/login" "Login Admin" '{"email":"admin@example.com","password":"admin123"}' ""

echo -e "${GREEN}Admin Token: ${ADMIN_TOKEN:0:50}...${NC}"
echo ""

# Login with regular user
USER_LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"testroute_${RANDOM_STR}@test.com\",\"password\":\"Test123!\"}")

USER_TOKEN=$(echo $USER_LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4)

echo -e "${GREEN}User Token: ${USER_TOKEN:0:50}...${NC}"
echo ""

# ==========================================
# 2. USER ROUTES
# ==========================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}2. USER ROUTES${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

test_endpoint "GET" "/users" "Get All Users (Admin)" "" "$ADMIN_TOKEN"
test_endpoint "GET" "/users?page=1&limit=5" "Get Users with Pagination (Admin)" "" "$ADMIN_TOKEN"
test_endpoint "GET" "/users/1" "Get User by ID (Admin)" "" "$ADMIN_TOKEN"
test_endpoint "GET" "/profile" "Get Own Profile (User)" "" "$USER_TOKEN"

# ==========================================
# 3. MAHASISWA ROUTES
# ==========================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}3. MAHASISWA ROUTES${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# Create Mahasiswa
MHS_NIM="M$(date +%s)"
CREATE_MHS_DATA="{\"nim\":\"$MHS_NIM\",\"nama\":\"Test Mahasiswa Route\",\"jurusan\":\"Teknik Informatika\",\"angkatan\":2024,\"email\":\"mhs_route_${RANDOM_STR}@test.com\"}"
CREATE_MHS_RESPONSE=$(curl -s -X POST "$BASE_URL/mahasiswa" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "$CREATE_MHS_DATA")

MHS_ID=$(echo $CREATE_MHS_RESPONSE | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')

test_endpoint "POST" "/mahasiswa" "Create Mahasiswa (Admin)" "$CREATE_MHS_DATA" "$ADMIN_TOKEN"
test_endpoint "GET" "/mahasiswa" "Get All Mahasiswa (User)" "" "$USER_TOKEN"
test_endpoint "GET" "/mahasiswa?page=1&limit=5" "Get Mahasiswa with Pagination" "" "$USER_TOKEN"
test_endpoint "GET" "/mahasiswa/count" "Get Mahasiswa Count" "" "$USER_TOKEN"

if [ -n "$MHS_ID" ]; then
    test_endpoint "GET" "/mahasiswa/$MHS_ID" "Get Mahasiswa by ID" "" "$USER_TOKEN"
    test_endpoint "PUT" "/mahasiswa/$MHS_ID" "Update Mahasiswa (Admin)" "{\"nama\":\"Updated Mahasiswa\",\"jurusan\":\"Sistem Informasi\",\"angkatan\":2024,\"email\":\"mhs_updated_${RANDOM_STR}@test.com\"}" "$ADMIN_TOKEN"
    test_endpoint "GET" "/mahasiswa/$MHS_ID" "Get Updated Mahasiswa" "" "$USER_TOKEN"
fi

test_endpoint "GET" "/mahasiswa/search?q=Test" "Search Mahasiswa by Query" "" "$USER_TOKEN"
test_endpoint "GET" "/mahasiswa/filter?jurusan=Teknik%20Informatika" "Filter Mahasiswa by Jurusan" "" "$USER_TOKEN"

# ==========================================
# 4. ALUMNI ROUTES
# ==========================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}4. ALUMNI ROUTES${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# Get a user ID for alumni
USER_ID=$(echo $USER_LOGIN_RESPONSE | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')

# Create Alumni
ALUMNI_NIM="A$(date +%s)"
CREATE_ALUMNI_DATA="{\"nim\":\"$ALUMNI_NIM\",\"nama\":\"Test Alumni Route\",\"jurusan\":\"Teknik Informatika\",\"angkatan\":2020,\"tahun_lulus\":2024,\"user_id\":$USER_ID,\"no_telepon\":\"081234567890\",\"alamat\":\"Jalan Test\"}"
CREATE_ALUMNI_RESPONSE=$(curl -s -X POST "$BASE_URL/alumni" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d "$CREATE_ALUMNI_DATA")

ALUMNI_ID=$(echo $CREATE_ALUMNI_RESPONSE | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')

test_endpoint "POST" "/alumni" "Create Alumni (Admin)" "$CREATE_ALUMNI_DATA" "$ADMIN_TOKEN"
test_endpoint "GET" "/alumni" "Get All Alumni (User)" "" "$USER_TOKEN"
test_endpoint "GET" "/alumni?page=1&limit=5" "Get Alumni with Pagination" "" "$USER_TOKEN"
test_endpoint "GET" "/alumni/count" "Get Alumni Count" "" "$USER_TOKEN"

if [ -n "$ALUMNI_ID" ]; then
    test_endpoint "GET" "/alumni/$ALUMNI_ID" "Get Alumni by ID" "" "$USER_TOKEN"
    test_endpoint "PUT" "/alumni/$ALUMNI_ID" "Update Alumni (Admin)" "{\"nama\":\"Updated Alumni\",\"jurusan\":\"Sistem Informasi\",\"angkatan\":2020,\"tahun_lulus\":2024,\"no_telepon\":\"081234567899\",\"alamat\":\"Jalan Updated\"}" "$ADMIN_TOKEN"
    test_endpoint "GET" "/alumni/$ALUMNI_ID" "Get Updated Alumni" "" "$USER_TOKEN"
fi

test_endpoint "GET" "/alumni/search?q=Test" "Search Alumni by Query" "" "$USER_TOKEN"
test_endpoint "GET" "/alumni/filter?jurusan=Teknik%20Informatika" "Filter Alumni by Jurusan" "" "$USER_TOKEN"
test_endpoint "GET" "/alumni/my-profile" "Get Alumni Profile (User)" "" "$USER_TOKEN"

# ==========================================
# 5. PEKERJAAN ALUMNI ROUTES
# ==========================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}5. PEKERJAAN ALUMNI ROUTES${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

if [ -n "$ALUMNI_ID" ]; then
    # Create Pekerjaan
    CREATE_PEKERJAAN_DATA="{\"alumni_id\":$ALUMNI_ID,\"nama_perusahaan\":\"PT Test Company\",\"posisi_jabatan\":\"Software Engineer\",\"bidang_industri\":\"Teknologi Informasi\",\"lokasi_kerja\":\"Jakarta\",\"gaji_range\":\"10-15 juta\",\"tanggal_mulai_kerja\":\"2024-01-01T00:00:00Z\",\"status_pekerjaan\":\"aktif\",\"deskripsi_pekerjaan\":\"Test job description\"}"
    CREATE_PEKERJAAN_RESPONSE=$(curl -s -X POST "$BASE_URL/pekerjaan" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $ADMIN_TOKEN" \
      -d "$CREATE_PEKERJAAN_DATA")

    PEKERJAAN_ID=$(echo $CREATE_PEKERJAAN_RESPONSE | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')

    test_endpoint "POST" "/pekerjaan" "Create Pekerjaan Alumni (Admin)" "$CREATE_PEKERJAAN_DATA" "$ADMIN_TOKEN"
    test_endpoint "GET" "/pekerjaan" "Get All Pekerjaan (User)" "" "$USER_TOKEN"
    test_endpoint "GET" "/pekerjaan?page=1&limit=5" "Get Pekerjaan with Pagination" "" "$USER_TOKEN"
    test_endpoint "GET" "/pekerjaan/count" "Get Pekerjaan Count" "" "$USER_TOKEN"

    if [ -n "$PEKERJAAN_ID" ]; then
        test_endpoint "GET" "/pekerjaan/$PEKERJAAN_ID" "Get Pekerjaan by ID" "" "$USER_TOKEN"
        test_endpoint "PUT" "/pekerjaan/$PEKERJAAN_ID" "Update Pekerjaan (Admin)" "{\"nama_perusahaan\":\"PT Updated Company\",\"posisi_jabatan\":\"Senior Software Engineer\",\"bidang_industri\":\"Teknologi Informasi\",\"lokasi_kerja\":\"Bandung\",\"gaji_range\":\"15-20 juta\",\"tanggal_mulai_kerja\":\"2024-01-01T00:00:00Z\",\"status_pekerjaan\":\"aktif\",\"deskripsi_pekerjaan\":\"Updated job description\"}" "$ADMIN_TOKEN"
        test_endpoint "GET" "/pekerjaan/$PEKERJAAN_ID" "Get Updated Pekerjaan" "" "$USER_TOKEN"
    fi

    test_endpoint "GET" "/pekerjaan/alumni/$ALUMNI_ID" "Get Pekerjaan by Alumni ID" "" "$USER_TOKEN"
    test_endpoint "GET" "/pekerjaan/search?q=Test" "Search Pekerjaan by Query" "" "$USER_TOKEN"
    test_endpoint "GET" "/pekerjaan/filter?status_pekerjaan=aktif" "Filter Pekerjaan by Status" "" "$USER_TOKEN"
fi

# ==========================================
# 6. TRASH/SOFT DELETE ROUTES
# ==========================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}6. TRASH/SOFT DELETE ROUTES${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

# Soft delete pekerjaan
if [ -n "$PEKERJAAN_ID" ]; then
    test_endpoint "DELETE" "/pekerjaan/$PEKERJAAN_ID" "Soft Delete Pekerjaan (Admin)" "" "$ADMIN_TOKEN"
    test_endpoint "GET" "/trash/pekerjaan" "Get Trashed Pekerjaan (Admin)" "" "$ADMIN_TOKEN"
    test_endpoint "POST" "/trash/pekerjaan/$PEKERJAAN_ID/restore" "Restore Pekerjaan (Admin)" "" "$ADMIN_TOKEN"
    test_endpoint "DELETE" "/pekerjaan/$PEKERJAAN_ID" "Soft Delete Pekerjaan Again (Admin)" "" "$ADMIN_TOKEN"
    test_endpoint "DELETE" "/trash/pekerjaan/$PEKERJAAN_ID" "Permanent Delete Pekerjaan (Admin)" "" "$ADMIN_TOKEN"
fi

# Delete Alumni
if [ -n "$ALUMNI_ID" ]; then
    test_endpoint "DELETE" "/alumni/$ALUMNI_ID" "Delete Alumni (Admin)" "" "$ADMIN_TOKEN"
fi

# Delete Mahasiswa
if [ -n "$MHS_ID" ]; then
    test_endpoint "DELETE" "/mahasiswa/$MHS_ID" "Delete Mahasiswa (Admin)" "" "$ADMIN_TOKEN"
fi

# ==========================================
# 7. STATISTICS & REPORTS
# ==========================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}7. STATISTICS & REPORTS${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

test_endpoint "GET" "/alumni/stats/by-year" "Get Alumni Statistics by Year" "" "$USER_TOKEN"
test_endpoint "GET" "/alumni/stats/by-jurusan" "Get Alumni Statistics by Jurusan" "" "$USER_TOKEN"
test_endpoint "GET" "/pekerjaan/stats/by-industry" "Get Pekerjaan Statistics by Industry" "" "$USER_TOKEN"
test_endpoint "GET" "/pekerjaan/stats/by-location" "Get Pekerjaan Statistics by Location" "" "$USER_TOKEN"

# ==========================================
# 8. PERMISSION TESTS (User trying Admin routes)
# ==========================================
echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}8. PERMISSION TESTS (Should Fail)${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""

echo -e "${CYAN}Testing: User trying to create Mahasiswa (Should Fail)${NC}"
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "$BASE_URL/mahasiswa" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $USER_TOKEN" \
  -d '{"nim":"FAIL123","nama":"Should Fail","jurusan":"Test","angkatan":2024,"email":"fail@test.com"}')
HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d':' -f2)
if [ "$HTTP_CODE" -eq 403 ] || [ "$HTTP_CODE" -eq 401 ]; then
    echo -e "${GREEN}✅ Correctly denied (HTTP $HTTP_CODE)${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${RED}❌ Permission check failed (HTTP $HTTP_CODE)${NC}"
    ((FAIL_COUNT++))
fi
echo ""

echo -e "${CYAN}Testing: User trying to delete Mahasiswa (Should Fail)${NC}"
RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X DELETE "$BASE_URL/mahasiswa/1" \
  -H "Authorization: Bearer $USER_TOKEN")
HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d':' -f2)
if [ "$HTTP_CODE" -eq 403 ] || [ "$HTTP_CODE" -eq 401 ]; then
    echo -e "${GREEN}✅ Correctly denied (HTTP $HTTP_CODE)${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${RED}❌ Permission check failed (HTTP $HTTP_CODE)${NC}"
    ((FAIL_COUNT++))
fi
echo ""

# ==========================================
# FINAL SUMMARY
# ==========================================
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}           TEST SUMMARY                ${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "${GREEN}✅ Successful Tests: $SUCCESS_COUNT${NC}"
echo -e "${RED}❌ Failed Tests: $FAIL_COUNT${NC}"
echo -e "${BLUE}Total Tests: $((SUCCESS_COUNT + FAIL_COUNT))${NC}"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}   ALL TESTS PASSED! 🎉               ${NC}"
    echo -e "${GREEN}========================================${NC}"
else
    echo -e "${YELLOW}========================================${NC}"
    echo -e "${YELLOW}   Some tests failed. Review above.   ${NC}"
    echo -e "${YELLOW}========================================${NC}"
fi

# Show database counts
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}      CURRENT DATABASE STATE           ${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

USERS_COUNT=$(curl -s "http://localhost:8080/api/users" -H "Authorization: Bearer $ADMIN_TOKEN" | grep -o '"total_data":[0-9]*' | grep -o '[0-9]*')
echo -e "${CYAN}Users in database: ${GREEN}$USERS_COUNT${NC}"

MHS_COUNT=$(curl -s "http://localhost:8080/api/mahasiswa?page=1&limit=1" -H "Authorization: Bearer $ADMIN_TOKEN" | grep -o '"total_data":[0-9]*' | grep -o '[0-9]*')
echo -e "${CYAN}Mahasiswa in database: ${GREEN}$MHS_COUNT${NC}"

ALUMNI_COUNT=$(curl -s "http://localhost:8080/api/alumni?page=1&limit=1" -H "Authorization: Bearer $ADMIN_TOKEN" | grep -o '"total_data":[0-9]*' | grep -o '[0-9]*')
echo -e "${CYAN}Alumni in database: ${GREEN}$ALUMNI_COUNT${NC}"

PEKERJAAN_COUNT=$(curl -s "http://localhost:8080/api/pekerjaan?page=1&limit=1" -H "Authorization: Bearer $ADMIN_TOKEN" | grep -o '"total_data":[0-9]*' | grep -o '[0-9]*')
echo -e "${CYAN}Pekerjaan in database: ${GREEN}$PEKERJAAN_COUNT${NC}"

echo ""
echo -e "${BLUE}Total Records: ${GREEN}$((USERS_COUNT + MHS_COUNT + ALUMNI_COUNT + PEKERJAAN_COUNT))${NC}"
echo -e "${BLUE}========================================${NC}"
