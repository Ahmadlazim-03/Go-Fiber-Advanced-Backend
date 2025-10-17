#!/bin/bash

# Warna untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
MAGENTA='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080/api"

# Counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Function to log test result
log_test() {
    local test_name=$1
    local status=$2
    local message=$3
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if [ "$status" == "PASS" ]; then
        PASSED_TESTS=$((PASSED_TESTS + 1))
        echo -e "${GREEN}✅ PASS${NC} | $test_name"
        [ -n "$message" ] && echo -e "   ${CYAN}→ $message${NC}"
    else
        FAILED_TESTS=$((FAILED_TESTS + 1))
        echo -e "${RED}❌ FAIL${NC} | $test_name"
        [ -n "$message" ] && echo -e "   ${RED}→ $message${NC}"
    fi
}

# Function to check HTTP status
check_status() {
    local response=$1
    local expected=$2
    
    if echo "$response" | grep -q "HTTP.*$expected"; then
        echo "PASS"
    else
        echo "FAIL"
    fi
}

# Function to check if response contains data
check_data() {
    local response=$1
    local field=$2
    
    if echo "$response" | grep -q "\"$field\""; then
        echo "PASS"
    else
        echo "FAIL"
    fi
}

echo -e "${CYAN}============================================${NC}"
echo -e "${CYAN}   COMPREHENSIVE DATABASE ROUTES TEST     ${NC}"
echo -e "${CYAN}============================================${NC}"
echo ""

# ==========================================
# TEST 1: LOGIN & AUTHENTICATION
# ==========================================
echo -e "${MAGENTA}═══════════════════════════════════════════${NC}"
echo -e "${MAGENTA} TEST 1: AUTHENTICATION                    ${NC}"
echo -e "${MAGENTA}═══════════════════════════════════════════${NC}"
echo ""

# Login as admin
LOGIN_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"admin123"}')

HTTP_STATUS=$(echo "$LOGIN_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$LOGIN_RESPONSE" | sed '/HTTP_STATUS/d')

if [ "$HTTP_STATUS" == "200" ]; then
    TOKEN=$(echo "$BODY" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    if [ -n "$TOKEN" ]; then
        log_test "Admin Login" "PASS" "Token received"
    else
        log_test "Admin Login" "FAIL" "No token in response"
        exit 1
    fi
else
    log_test "Admin Login" "FAIL" "HTTP $HTTP_STATUS"
    exit 1
fi
echo ""

# ==========================================
# TEST 2: USERS ROUTES
# ==========================================
echo -e "${MAGENTA}═══════════════════════════════════════════${NC}"
echo -e "${MAGENTA} TEST 2: USERS ROUTES                      ${NC}"
echo -e "${MAGENTA}═══════════════════════════════════════════${NC}"
echo ""

# GET all users
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/users" \
  -H "Authorization: Bearer $TOKEN")
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_STATUS/d')

if [ "$HTTP_STATUS" == "200" ] && echo "$BODY" | grep -q '"total_data"'; then
    TOTAL=$(echo "$BODY" | grep -o '"total_data":[0-9]*' | grep -o '[0-9]*')
    log_test "GET /users" "PASS" "Total users: $TOTAL"
else
    log_test "GET /users" "FAIL" "HTTP $HTTP_STATUS"
fi

# GET users count
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/users/count" \
  -H "Authorization: Bearer $TOKEN")
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)

if [ "$HTTP_STATUS" == "200" ]; then
    log_test "GET /users/count" "PASS"
else
    log_test "GET /users/count" "FAIL" "HTTP $HTTP_STATUS"
fi

# GET user by ID
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/users/1" \
  -H "Authorization: Bearer $TOKEN")
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)

if [ "$HTTP_STATUS" == "200" ]; then
    log_test "GET /users/:id" "PASS"
else
    log_test "GET /users/:id" "FAIL" "HTTP $HTTP_STATUS"
fi

# Create new user
RANDOM_STRING=$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1)
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"testuser_$RANDOM_STRING\",\"email\":\"test_${RANDOM_STRING}@test.com\",\"password\":\"Test123!\",\"role\":\"user\"}")
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_STATUS/d')

if [ "$HTTP_STATUS" == "201" ] || [ "$HTTP_STATUS" == "200" ]; then
    NEW_USER_ID=$(echo "$BODY" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    log_test "POST /register (Create User)" "PASS" "User ID: $NEW_USER_ID"
else
    log_test "POST /register (Create User)" "FAIL" "HTTP $HTTP_STATUS"
fi

# Update user
if [ -n "$NEW_USER_ID" ]; then
    RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X PUT "$BASE_URL/users/$NEW_USER_ID" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{\"username\":\"updated_user_$RANDOM_STRING\",\"email\":\"updated_${RANDOM_STRING}@test.com\",\"role\":\"user\"}")
    HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
    
    if [ "$HTTP_STATUS" == "200" ]; then
        log_test "PUT /users/:id" "PASS"
    else
        log_test "PUT /users/:id" "FAIL" "HTTP $HTTP_STATUS"
    fi
    
    # Delete user
    RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X DELETE "$BASE_URL/users/$NEW_USER_ID" \
      -H "Authorization: Bearer $TOKEN")
    HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
    
    if [ "$HTTP_STATUS" == "200" ]; then
        log_test "DELETE /users/:id" "PASS"
    else
        log_test "DELETE /users/:id" "FAIL" "HTTP $HTTP_STATUS"
    fi
fi

echo ""

# ==========================================
# TEST 3: MAHASISWA ROUTES
# ==========================================
echo -e "${MAGENTA}═══════════════════════════════════════════${NC}"
echo -e "${MAGENTA} TEST 3: MAHASISWA ROUTES                  ${NC}"
echo -e "${MAGENTA}═══════════════════════════════════════════${NC}"
echo ""

# GET all mahasiswa
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/mahasiswa?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN")
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_STATUS/d')

if [ "$HTTP_STATUS" == "200" ] && echo "$BODY" | grep -q '"total_data"'; then
    TOTAL=$(echo "$BODY" | grep -o '"total_data":[0-9]*' | grep -o '[0-9]*')
    log_test "GET /mahasiswa" "PASS" "Total: $TOTAL"
else
    log_test "GET /mahasiswa" "FAIL" "HTTP $HTTP_STATUS"
fi

# GET mahasiswa count
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/mahasiswa/count" \
  -H "Authorization: Bearer $TOKEN")
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)

if [ "$HTTP_STATUS" == "200" ]; then
    log_test "GET /mahasiswa/count" "PASS"
else
    log_test "GET /mahasiswa/count" "FAIL" "HTTP $HTTP_STATUS"
fi

# Create mahasiswa
RANDOM_NIM=$(cat /dev/urandom | tr -dc '0-9' | fold -w 10 | head -n 1)
RANDOM_EMAIL=$(cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 8 | head -n 1)
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$BASE_URL/mahasiswa" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"nim\":\"$RANDOM_NIM\",\"nama\":\"Test Mahasiswa\",\"jurusan\":\"Teknik Informatika\",\"angkatan\":2024,\"email\":\"mhs_${RANDOM_EMAIL}@student.ac.id\"}")
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_STATUS/d')

if [ "$HTTP_STATUS" == "201" ] || [ "$HTTP_STATUS" == "200" ]; then
    MHS_ID=$(echo "$BODY" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    log_test "POST /mahasiswa" "PASS" "Mahasiswa ID: $MHS_ID"
else
    log_test "POST /mahasiswa" "FAIL" "HTTP $HTTP_STATUS"
fi

# GET mahasiswa by ID
if [ -n "$MHS_ID" ]; then
    RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/mahasiswa/$MHS_ID" \
      -H "Authorization: Bearer $TOKEN")
    HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
    
    if [ "$HTTP_STATUS" == "200" ]; then
        log_test "GET /mahasiswa/:id" "PASS"
    else
        log_test "GET /mahasiswa/:id" "FAIL" "HTTP $HTTP_STATUS"
    fi
    
    # Update mahasiswa
    RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X PUT "$BASE_URL/mahasiswa/$MHS_ID" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{\"nama\":\"Updated Mahasiswa\",\"jurusan\":\"Sistem Informasi\",\"angkatan\":2024,\"email\":\"mhs_${RANDOM_EMAIL}@student.ac.id\"}")
    HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
    
    if [ "$HTTP_STATUS" == "200" ]; then
        log_test "PUT /mahasiswa/:id" "PASS"
    else
        log_test "PUT /mahasiswa/:id" "FAIL" "HTTP $HTTP_STATUS"
    fi
    
    # Delete mahasiswa
    RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X DELETE "$BASE_URL/mahasiswa/$MHS_ID" \
      -H "Authorization: Bearer $TOKEN")
    HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
    
    if [ "$HTTP_STATUS" == "200" ]; then
        log_test "DELETE /mahasiswa/:id" "PASS"
    else
        log_test "DELETE /mahasiswa/:id" "FAIL" "HTTP $HTTP_STATUS"
    fi
fi

echo ""

# ==========================================
# TEST 4: ALUMNI ROUTES
# ==========================================
echo -e "${MAGENTA}═══════════════════════════════════════════${NC}"
echo -e "${MAGENTA} TEST 4: ALUMNI ROUTES                     ${NC}"
echo -e "${MAGENTA}═══════════════════════════════════════════${NC}"
echo ""

# GET all alumni
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/alumni?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN")
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_STATUS/d')

if [ "$HTTP_STATUS" == "200" ] && echo "$BODY" | grep -q '"total_data"'; then
    TOTAL=$(echo "$BODY" | grep -o '"total_data":[0-9]*' | grep -o '[0-9]*')
    log_test "GET /alumni" "PASS" "Total: $TOTAL"
else
    log_test "GET /alumni" "FAIL" "HTTP $HTTP_STATUS"
fi

# GET alumni count
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/alumni/count" \
  -H "Authorization: Bearer $TOKEN")
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)

if [ "$HTTP_STATUS" == "200" ]; then
    log_test "GET /alumni/count" "PASS"
else
    log_test "GET /alumni/count" "FAIL" "HTTP $HTTP_STATUS"
fi

# Create alumni
RANDOM_NIM=$(cat /dev/urandom | tr -dc '0-9' | fold -w 10 | head -n 1)
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$BASE_URL/alumni" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{\"nim\":\"$RANDOM_NIM\",\"nama\":\"Test Alumni\",\"jurusan\":\"Teknik Informatika\",\"angkatan\":2020,\"tahun_lulus\":2024,\"user_id\":1,\"no_telepon\":\"081234567890\",\"alamat\":\"Test Address\"}")
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_STATUS/d')

if [ "$HTTP_STATUS" == "201" ] || [ "$HTTP_STATUS" == "200" ]; then
    ALUMNI_ID=$(echo "$BODY" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
    log_test "POST /alumni" "PASS" "Alumni ID: $ALUMNI_ID"
else
    log_test "POST /alumni" "FAIL" "HTTP $HTTP_STATUS"
fi

# GET alumni by ID
if [ -n "$ALUMNI_ID" ]; then
    RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/alumni/$ALUMNI_ID" \
      -H "Authorization: Bearer $TOKEN")
    HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
    
    if [ "$HTTP_STATUS" == "200" ]; then
        log_test "GET /alumni/:id" "PASS"
    else
        log_test "GET /alumni/:id" "FAIL" "HTTP $HTTP_STATUS"
    fi
    
    # Update alumni
    RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X PUT "$BASE_URL/alumni/$ALUMNI_ID" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{\"nama\":\"Updated Alumni\",\"jurusan\":\"Sistem Informasi\",\"angkatan\":2020,\"tahun_lulus\":2024,\"no_telepon\":\"081234567891\",\"alamat\":\"Updated Address\"}")
    HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
    
    if [ "$HTTP_STATUS" == "200" ]; then
        log_test "PUT /alumni/:id" "PASS"
    else
        log_test "PUT /alumni/:id" "FAIL" "HTTP $HTTP_STATUS"
    fi
    
    # Delete alumni
    RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X DELETE "$BASE_URL/alumni/$ALUMNI_ID" \
      -H "Authorization: Bearer $TOKEN")
    HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
    
    if [ "$HTTP_STATUS" == "200" ]; then
        log_test "DELETE /alumni/:id" "PASS"
    else
        log_test "DELETE /alumni/:id" "FAIL" "HTTP $HTTP_STATUS"
    fi
fi

echo ""

# ==========================================
# TEST 5: PEKERJAAN ALUMNI ROUTES
# ==========================================
echo -e "${MAGENTA}═══════════════════════════════════════════${NC}"
echo -e "${MAGENTA} TEST 5: PEKERJAAN ALUMNI ROUTES           ${NC}"
echo -e "${MAGENTA}═══════════════════════════════════════════${NC}"
echo ""

# GET all pekerjaan
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/pekerjaan?page=1&limit=10" \
  -H "Authorization: Bearer $TOKEN")
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$RESPONSE" | sed '/HTTP_STATUS/d')

if [ "$HTTP_STATUS" == "200" ] && echo "$BODY" | grep -q '"total_data"'; then
    TOTAL=$(echo "$BODY" | grep -o '"total_data":[0-9]*' | grep -o '[0-9]*')
    log_test "GET /pekerjaan" "PASS" "Total: $TOTAL"
else
    log_test "GET /pekerjaan" "FAIL" "HTTP $HTTP_STATUS"
fi

# GET pekerjaan count
RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/pekerjaan/count" \
  -H "Authorization: Bearer $TOKEN")
HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)

if [ "$HTTP_STATUS" == "200" ]; then
    log_test "GET /pekerjaan/count" "PASS"
else
    log_test "GET /pekerjaan/count" "FAIL" "HTTP $HTTP_STATUS"
fi

# Get an existing alumni ID for pekerjaan test
EXISTING_ALUMNI=$(curl -s -X GET "$BASE_URL/alumni?page=1&limit=1" -H "Authorization: Bearer $TOKEN" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')

if [ -n "$EXISTING_ALUMNI" ]; then
    # Create pekerjaan
    RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$BASE_URL/pekerjaan" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $TOKEN" \
      -d "{\"alumni_id\":$EXISTING_ALUMNI,\"nama_perusahaan\":\"Test Company\",\"posisi_jabatan\":\"Software Engineer\",\"bidang_industri\":\"IT\",\"lokasi_kerja\":\"Jakarta\",\"gaji_range\":\"10-15 juta\",\"tanggal_mulai_kerja\":\"2024-01-01T00:00:00Z\",\"status_pekerjaan\":\"aktif\",\"deskripsi_pekerjaan\":\"Test description\"}")
    HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
    BODY=$(echo "$RESPONSE" | sed '/HTTP_STATUS/d')
    
    if [ "$HTTP_STATUS" == "201" ] || [ "$HTTP_STATUS" == "200" ]; then
        PEKERJAAN_ID=$(echo "$BODY" | grep -o '"id":[0-9]*' | head -1 | grep -o '[0-9]*')
        log_test "POST /pekerjaan" "PASS" "Pekerjaan ID: $PEKERJAAN_ID"
    else
        log_test "POST /pekerjaan" "FAIL" "HTTP $HTTP_STATUS"
    fi
    
    # GET pekerjaan by ID
    if [ -n "$PEKERJAAN_ID" ]; then
        RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$BASE_URL/pekerjaan/$PEKERJAAN_ID" \
          -H "Authorization: Bearer $TOKEN")
        HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
        
        if [ "$HTTP_STATUS" == "200" ]; then
            log_test "GET /pekerjaan/:id" "PASS"
        else
            log_test "GET /pekerjaan/:id" "FAIL" "HTTP $HTTP_STATUS"
        fi
        
        # Update pekerjaan
        RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X PUT "$BASE_URL/pekerjaan/$PEKERJAAN_ID" \
          -H "Content-Type: application/json" \
          -H "Authorization: Bearer $TOKEN" \
          -d "{\"nama_perusahaan\":\"Updated Company\",\"posisi_jabatan\":\"Senior Engineer\",\"bidang_industri\":\"IT\",\"lokasi_kerja\":\"Bandung\",\"gaji_range\":\"15-20 juta\",\"tanggal_mulai_kerja\":\"2024-01-01T00:00:00Z\",\"status_pekerjaan\":\"aktif\",\"deskripsi_pekerjaan\":\"Updated description\"}")
        HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
        
        if [ "$HTTP_STATUS" == "200" ]; then
            log_test "PUT /pekerjaan/:id" "PASS"
        else
            log_test "PUT /pekerjaan/:id" "FAIL" "HTTP $HTTP_STATUS"
        fi
        
        # Delete pekerjaan
        RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X DELETE "$BASE_URL/pekerjaan/$PEKERJAAN_ID" \
          -H "Authorization: Bearer $TOKEN")
        HTTP_STATUS=$(echo "$RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
        
        if [ "$HTTP_STATUS" == "200" ]; then
            log_test "DELETE /pekerjaan/:id" "PASS"
        else
            log_test "DELETE /pekerjaan/:id" "FAIL" "HTTP $HTTP_STATUS"
        fi
    fi
else
    log_test "POST /pekerjaan" "FAIL" "No alumni available for testing"
fi

echo ""

# ==========================================
# FINAL SUMMARY
# ==========================================
echo -e "${CYAN}============================================${NC}"
echo -e "${CYAN}           TEST SUMMARY                     ${NC}"
echo -e "${CYAN}============================================${NC}"
echo -e "${BLUE}Total Tests:${NC}  $TOTAL_TESTS"
echo -e "${GREEN}Passed:${NC}       $PASSED_TESTS"
echo -e "${RED}Failed:${NC}       $FAILED_TESTS"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✅ ALL TESTS PASSED!${NC}"
    SUCCESS_RATE=100
else
    SUCCESS_RATE=$((PASSED_TESTS * 100 / TOTAL_TESTS))
    echo -e "${YELLOW}⚠️  SOME TESTS FAILED${NC}"
fi

echo -e "${BLUE}Success Rate: ${SUCCESS_RATE}%${NC}"
echo -e "${CYAN}============================================${NC}"
echo ""

# Exit with appropriate code
if [ $FAILED_TESTS -eq 0 ]; then
    exit 0
else
    exit 1
fi
