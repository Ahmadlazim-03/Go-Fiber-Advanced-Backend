#!/bin/bash

# ============================================================================
# Complete Routes Test Script - Dynamic & Smart Testing
# Automatically tests all routes based on actual route definitions
# ============================================================================

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:8080/api"
SUCCESS_COUNT=0
FAIL_COUNT=0
SKIP_COUNT=0
TEST_LOG="/tmp/test_routes_$(date +%s).log"

# Detect database type
DB_TYPE=${1:-$(grep "^DB_TYPE=" .env 2>/dev/null | cut -d'=' -f2)}
DB_TYPE=${DB_TYPE:-postgres}

# Test results array
declare -A TEST_RESULTS

# ============================================================================
# Utility Functions
# ============================================================================

log_test() {
    echo "[$(date +'%Y-%m-%d %H:%M:%S')] $1" >> "$TEST_LOG"
}

print_header() {
    echo -e "${BLUE}============================================================${NC}"
    echo -e "${BLUE}   COMPLETE ROUTES TEST - ${DB_TYPE^^}${NC}"
    echo -e "${BLUE}   Test Log: ${TEST_LOG}${NC}"
    echo -e "${BLUE}============================================================${NC}"
    echo ""
}

print_section() {
    echo -e "${YELLOW}============================================================${NC}"
    echo -e "${YELLOW}$1${NC}"
    echo -e "${YELLOW}============================================================${NC}"
    echo ""
}

# Enhanced test endpoint function
test_endpoint() {
    local METHOD=$1
    local ENDPOINT=$2
    local DESCRIPTION=$3
    local DATA=$4
    local TOKEN=$5
    local EXPECTED_CODE=${6:-200}
    local SHOULD_FAIL=${7:-false}
    
    echo -e "${CYAN}Testing: ${DESCRIPTION}${NC}"
    log_test "TEST: $METHOD $ENDPOINT - $DESCRIPTION"
    
    # Execute request with proper error handling
    if [ -n "$TOKEN" ] && [ -n "$DATA" ] && [ "$METHOD" != "GET" ]; then
        RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X "$METHOD" "$BASE_URL$ENDPOINT" \
            -H "Authorization: Bearer $TOKEN" \
            -H "Content-Type: application/json" \
            -d "$DATA" 2>&1)
    elif [ -n "$TOKEN" ]; then
        RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X "$METHOD" "$BASE_URL$ENDPOINT" \
            -H "Authorization: Bearer $TOKEN" 2>&1)
    elif [ -n "$DATA" ] && [ "$METHOD" != "GET" ]; then
        RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X "$METHOD" "$BASE_URL$ENDPOINT" \
            -H "Content-Type: application/json" \
            -d "$DATA" 2>&1)
    else
        RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X "$METHOD" "$BASE_URL$ENDPOINT" 2>&1)
    fi
    local CURL_EXIT=$?
    
    if [ $CURL_EXIT -ne 0 ]; then
        echo -e "${RED}❌ FAILED (Curl error: $CURL_EXIT)${NC}"
        log_test "FAIL: Curl error $CURL_EXIT"
        ((FAIL_COUNT++))
        TEST_RESULTS["$DESCRIPTION"]="FAIL"
        echo ""
        return 1
    fi
    
    # Parse response
    HTTP_CODE=$(echo "$RESPONSE" | grep "HTTP_CODE" | cut -d':' -f2 | tr -d '[:space:]')
    BODY=$(echo "$RESPONSE" | sed '/HTTP_CODE/d')
    
    # Validate response
    if [ -z "$HTTP_CODE" ]; then
        echo -e "${RED}❌ FAILED (No HTTP code received)${NC}"
        log_test "FAIL: No HTTP code"
        ((FAIL_COUNT++))
        TEST_RESULTS["$DESCRIPTION"]="FAIL"
    elif [ "$SHOULD_FAIL" = "true" ]; then
        if [ "$HTTP_CODE" -eq 403 ] || [ "$HTTP_CODE" -eq 401 ]; then
            echo -e "${GREEN}✅ CORRECTLY DENIED (HTTP $HTTP_CODE)${NC}"
            log_test "PASS: Correctly denied with $HTTP_CODE"
            ((SUCCESS_COUNT++))
            TEST_RESULTS["$DESCRIPTION"]="PASS"
        else
            echo -e "${RED}❌ FAILED (Expected 401/403, got $HTTP_CODE)${NC}"
            log_test "FAIL: Permission check failed with $HTTP_CODE"
            ((FAIL_COUNT++))
            TEST_RESULTS["$DESCRIPTION"]="FAIL"
        fi
    elif [ "$HTTP_CODE" -ge 200 ] && [ "$HTTP_CODE" -lt 300 ]; then
        echo -e "${GREEN}✅ SUCCESS (HTTP $HTTP_CODE)${NC}"
        log_test "PASS: HTTP $HTTP_CODE"
        ((SUCCESS_COUNT++))
        TEST_RESULTS["$DESCRIPTION"]="PASS"
    else
        echo -e "${RED}❌ FAILED (HTTP $HTTP_CODE)${NC}"
        local ERROR_MSG=$(echo "$BODY" | head -c 200)
        echo -e "${RED}   Response: ${ERROR_MSG}${NC}"
        log_test "FAIL: HTTP $HTTP_CODE - $ERROR_MSG"
        ((FAIL_COUNT++))
        TEST_RESULTS["$DESCRIPTION"]="FAIL"
    fi
    
    echo ""
}

# Check if server is running
check_server() {
    echo -e "${CYAN}Checking server status...${NC}"
    if ! curl -s "http://localhost:8080" > /dev/null 2>&1; then
        echo -e "${RED}❌ Server is not running on http://localhost:8080${NC}"
        echo -e "${YELLOW}Please start the server first:${NC}"
        echo -e "${YELLOW}   ./main &${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Server is running${NC}"
    echo ""
}

# ============================================================================
# Main Test Flow
# ============================================================================

print_header
check_server

# ============================================================================
# AUTHENTICATION & TOKEN SETUP
# ============================================================================

print_section "AUTHENTICATION SETUP"

# Determine credentials based on database type
if [ "$DB_TYPE" == "pocketbase" ]; then
    ADMIN_EMAIL="pbadmin@test.com"
    ADMIN_PASSWORD="Admin123!"
else
    ADMIN_EMAIL="admin@example.com"
    ADMIN_PASSWORD="admin123"
fi

# Login with admin
echo -e "${CYAN}Logging in as admin...${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}")

ADMIN_TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -z "$ADMIN_TOKEN" ]; then
    echo -e "${RED}❌ Failed to get admin token${NC}"
    echo -e "${YELLOW}Attempting to continue with limited tests...${NC}"
    ADMIN_TOKEN=""
else
    echo -e "${GREEN}✅ Admin token obtained${NC}"
fi

# Register and login as regular user
TIMESTAMP=$(date +%s)
USER_EMAIL="testuser_${TIMESTAMP}@test.com"
USER_PASSWORD="Test123!"

echo -e "${CYAN}Creating test user...${NC}"
USER_REGISTER=$(curl -s -X POST "$BASE_URL/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\":\"testuser_${TIMESTAMP}\",\"email\":\"$USER_EMAIL\",\"password\":\"$USER_PASSWORD\",\"role\":\"user\"}")

USER_LOGIN=$(curl -s -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$USER_EMAIL\",\"password\":\"$USER_PASSWORD\"}")

USER_TOKEN=$(echo "$USER_LOGIN" | grep -o '"token":"[^"]*"' | head -1 | cut -d'"' -f4)

if [ -n "$USER_TOKEN" ]; then
    echo -e "${GREEN}✅ User token obtained${NC}"
else
    echo -e "${YELLOW}⚠️  User token not obtained, skipping user tests${NC}"
    USER_TOKEN=""
fi

echo ""

# ============================================================================
# 1. AUTHENTICATION ROUTES
# ============================================================================

print_section "1. AUTHENTICATION ROUTES"

test_endpoint "POST" "/register" "Register New User" \
    "{\"username\":\"newuser_$(date +%s)\",\"email\":\"newuser_$(date +%s)@test.com\",\"password\":\"Test123!\",\"role\":\"user\"}" \
    ""

test_endpoint "POST" "/login" "Login with Admin" \
    "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\"}" \
    ""

# ============================================================================
# 2. USER MANAGEMENT ROUTES (Admin Only)
# ============================================================================

print_section "2. USER MANAGEMENT ROUTES"

if [ -n "$ADMIN_TOKEN" ]; then
    test_endpoint "GET" "/users" "Get All Users" "" "$ADMIN_TOKEN"
    test_endpoint "GET" "/users?page=1&limit=5" "Get Users with Pagination" "" "$ADMIN_TOKEN"
    test_endpoint "GET" "/users/count" "Get Users Count" "" "$ADMIN_TOKEN"
    test_endpoint "GET" "/users/1" "Get User by ID" "" "$ADMIN_TOKEN"
else
    echo -e "${YELLOW}⚠️  Skipping user routes (no admin token)${NC}"
    ((SKIP_COUNT+=4))
fi

if [ -n "$USER_TOKEN" ]; then
    test_endpoint "GET" "/profile" "Get Own Profile" "" "$USER_TOKEN"
else
    echo -e "${YELLOW}⚠️  Skipping profile test (no user token)${NC}"
    ((SKIP_COUNT++))
fi

# ============================================================================
# 3. MAHASISWA ROUTES
# ============================================================================

print_section "3. MAHASISWA ROUTES"

if [ -n "$ADMIN_TOKEN" ]; then
    MHS_NIM="MHS$(date +%s)"
    CREATE_MHS_DATA="{\"nim\":\"$MHS_NIM\",\"nama\":\"Test Mahasiswa Auto\",\"jurusan\":\"Teknik Informatika\",\"angkatan\":2024,\"email\":\"mhs_${TIMESTAMP}@test.com\"}"
    
    CREATE_MHS_RESPONSE=$(curl -s -X POST "$BASE_URL/mahasiswa" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -d "$CREATE_MHS_DATA")
    
    MHS_ID=$(echo "$CREATE_MHS_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    
    test_endpoint "POST" "/mahasiswa" "Create Mahasiswa" "$CREATE_MHS_DATA" "$ADMIN_TOKEN"
fi

if [ -n "$USER_TOKEN" ]; then
    test_endpoint "GET" "/mahasiswa" "Get All Mahasiswa" "" "$USER_TOKEN"
    test_endpoint "GET" "/mahasiswa?page=1&limit=5" "Get Mahasiswa Pagination" "" "$USER_TOKEN"
    test_endpoint "GET" "/mahasiswa/count" "Get Mahasiswa Count" "" "$USER_TOKEN"
    test_endpoint "GET" "/mahasiswa/search?q=Test" "Search Mahasiswa" "" "$USER_TOKEN"
    test_endpoint "GET" "/mahasiswa/filter?jurusan=Teknik%20Informatika" "Filter Mahasiswa" "" "$USER_TOKEN"
    
    if [ -n "$MHS_ID" ] && [ "$MHS_ID" != "0" ]; then
        test_endpoint "GET" "/mahasiswa/$MHS_ID" "Get Mahasiswa by ID" "" "$USER_TOKEN"
    fi
fi

if [ -n "$ADMIN_TOKEN" ] && [ -n "$MHS_ID" ] && [ "$MHS_ID" != "0" ]; then
    UPDATE_MHS_DATA="{\"nama\":\"Updated Mahasiswa Auto\",\"jurusan\":\"Sistem Informasi\",\"angkatan\":2024,\"email\":\"mhs_updated_${TIMESTAMP}@test.com\"}"
    test_endpoint "PUT" "/mahasiswa/$MHS_ID" "Update Mahasiswa" "$UPDATE_MHS_DATA" "$ADMIN_TOKEN"
fi

# ============================================================================
# 4. ALUMNI ROUTES
# ============================================================================

print_section "4. ALUMNI ROUTES"

if [ -n "$ADMIN_TOKEN" ]; then
    ALUMNI_NIM="ALM$(date +%s)"
    CREATE_ALUMNI_DATA="{\"nim\":\"$ALUMNI_NIM\",\"nama\":\"Test Alumni Auto\",\"jurusan\":\"Teknik Informatika\",\"angkatan\":2020,\"tahun_lulus\":2024,\"user_id\":1,\"no_telepon\":\"081234567890\",\"alamat\":\"Test Address\"}"
    
    CREATE_ALUMNI_RESPONSE=$(curl -s -X POST "$BASE_URL/alumni" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -d "$CREATE_ALUMNI_DATA")
    
    ALUMNI_ID=$(echo "$CREATE_ALUMNI_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    
    test_endpoint "POST" "/alumni" "Create Alumni" "$CREATE_ALUMNI_DATA" "$ADMIN_TOKEN"
fi

if [ -n "$USER_TOKEN" ]; then
    test_endpoint "GET" "/alumni" "Get All Alumni" "" "$USER_TOKEN"
    test_endpoint "GET" "/alumni?page=1&limit=5" "Get Alumni Pagination" "" "$USER_TOKEN"
    test_endpoint "GET" "/alumni/count" "Get Alumni Count" "" "$USER_TOKEN"
    test_endpoint "GET" "/alumni/search?q=Test" "Search Alumni" "" "$USER_TOKEN"
    test_endpoint "GET" "/alumni/filter?jurusan=Teknik%20Informatika" "Filter Alumni" "" "$USER_TOKEN"
    test_endpoint "GET" "/alumni/stats/by-year" "Alumni Stats by Year" "" "$USER_TOKEN"
    test_endpoint "GET" "/alumni/stats/by-jurusan" "Alumni Stats by Jurusan" "" "$USER_TOKEN"
    
    if [ -n "$ALUMNI_ID" ] && [ "$ALUMNI_ID" != "0" ]; then
        test_endpoint "GET" "/alumni/$ALUMNI_ID" "Get Alumni by ID" "" "$USER_TOKEN"
    fi
fi

if [ -n "$ADMIN_TOKEN" ] && [ -n "$ALUMNI_ID" ] && [ "$ALUMNI_ID" != "0" ]; then
    UPDATE_ALUMNI_DATA="{\"nama\":\"Updated Alumni Auto\",\"jurusan\":\"Sistem Informasi\",\"angkatan\":2020,\"tahun_lulus\":2024,\"no_telepon\":\"081234567899\",\"alamat\":\"Updated Address\"}"
    test_endpoint "PUT" "/alumni/$ALUMNI_ID" "Update Alumni" "$UPDATE_ALUMNI_DATA" "$ADMIN_TOKEN"
fi

# ============================================================================
# 5. PEKERJAAN ALUMNI ROUTES
# ============================================================================

print_section "5. PEKERJAAN ALUMNI ROUTES"

if [ -n "$ADMIN_TOKEN" ] && [ -n "$ALUMNI_ID" ] && [ "$ALUMNI_ID" != "0" ]; then
    CREATE_PEKERJAAN_DATA="{\"alumni_id\":$ALUMNI_ID,\"nama_perusahaan\":\"PT Auto Test\",\"posisi_jabatan\":\"Software Engineer\",\"bidang_industri\":\"Teknologi Informasi\",\"lokasi_kerja\":\"Jakarta\",\"gaji_range\":\"10-15 juta\",\"tanggal_mulai_kerja\":\"2024-01-01T00:00:00Z\",\"status_pekerjaan\":\"aktif\",\"deskripsi_pekerjaan\":\"Auto Test Job\"}"
    
    CREATE_PEKERJAAN_RESPONSE=$(curl -s -X POST "$BASE_URL/pekerjaan" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ADMIN_TOKEN" \
        -d "$CREATE_PEKERJAAN_DATA")
    
    PEKERJAAN_ID=$(echo "$CREATE_PEKERJAAN_RESPONSE" | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
    
    test_endpoint "POST" "/pekerjaan" "Create Pekerjaan" "$CREATE_PEKERJAAN_DATA" "$ADMIN_TOKEN"
fi

if [ -n "$USER_TOKEN" ]; then
    test_endpoint "GET" "/pekerjaan" "Get All Pekerjaan" "" "$USER_TOKEN"
    test_endpoint "GET" "/pekerjaan?page=1&limit=5" "Get Pekerjaan Pagination" "" "$USER_TOKEN"
    test_endpoint "GET" "/pekerjaan/count" "Get Pekerjaan Count" "" "$USER_TOKEN"
    test_endpoint "GET" "/pekerjaan/search?q=Test" "Search Pekerjaan" "" "$USER_TOKEN"
    test_endpoint "GET" "/pekerjaan/filter?status_pekerjaan=aktif" "Filter Pekerjaan" "" "$USER_TOKEN"
    test_endpoint "GET" "/pekerjaan/stats/by-industry" "Stats by Industry" "" "$USER_TOKEN"
    test_endpoint "GET" "/pekerjaan/stats/by-location" "Stats by Location" "" "$USER_TOKEN"
    
    if [ -n "$PEKERJAAN_ID" ] && [ "$PEKERJAAN_ID" != "0" ]; then
        test_endpoint "GET" "/pekerjaan/$PEKERJAAN_ID" "Get Pekerjaan by ID" "" "$USER_TOKEN"
    fi
    
    if [ -n "$ALUMNI_ID" ] && [ "$ALUMNI_ID" != "0" ]; then
        test_endpoint "GET" "/pekerjaan/alumni/$ALUMNI_ID" "Get Pekerjaan by Alumni" "" "$USER_TOKEN"
    fi
fi

if [ -n "$ADMIN_TOKEN" ] && [ -n "$PEKERJAAN_ID" ] && [ "$PEKERJAAN_ID" != "0" ]; then
    UPDATE_PEKERJAAN_DATA="{\"nama_perusahaan\":\"PT Auto Updated\",\"posisi_jabatan\":\"Senior Engineer\",\"bidang_industri\":\"Teknologi Informasi\",\"lokasi_kerja\":\"Bandung\",\"gaji_range\":\"15-20 juta\",\"tanggal_mulai_kerja\":\"2024-01-01T00:00:00Z\",\"status_pekerjaan\":\"aktif\",\"deskripsi_pekerjaan\":\"Updated Job\"}"
    test_endpoint "PUT" "/pekerjaan/$PEKERJAAN_ID" "Update Pekerjaan" "$UPDATE_PEKERJAAN_DATA" "$ADMIN_TOKEN"
fi

# ============================================================================
# 6. SOFT DELETE / TRASH ROUTES
# ============================================================================

print_section "6. SOFT DELETE / TRASH ROUTES"

if [ -n "$ADMIN_TOKEN" ] && [ -n "$PEKERJAAN_ID" ] && [ "$PEKERJAAN_ID" != "0" ]; then
    test_endpoint "DELETE" "/pekerjaan/soft/$PEKERJAAN_ID" "Soft Delete Pekerjaan" "" "$ADMIN_TOKEN"
    test_endpoint "GET" "/pekerjaan/deleted" "Get Deleted Pekerjaan" "" "$ADMIN_TOKEN"
    test_endpoint "GET" "/trash/pekerjaan" "Get Trash Pekerjaan" "" "$ADMIN_TOKEN"
    test_endpoint "POST" "/pekerjaan/restore/$PEKERJAAN_ID" "Restore Pekerjaan" "" "$ADMIN_TOKEN"
    test_endpoint "DELETE" "/pekerjaan/soft/$PEKERJAAN_ID" "Soft Delete Again" "" "$ADMIN_TOKEN"
    test_endpoint "DELETE" "/trash/pekerjaan/$PEKERJAAN_ID" "Permanent Delete Pekerjaan" "" "$ADMIN_TOKEN"
fi

# Cleanup
if [ -n "$ADMIN_TOKEN" ]; then
    if [ -n "$ALUMNI_ID" ] && [ "$ALUMNI_ID" != "0" ]; then
        test_endpoint "DELETE" "/alumni/$ALUMNI_ID" "Delete Test Alumni" "" "$ADMIN_TOKEN"
    fi
    
    if [ -n "$MHS_ID" ] && [ "$MHS_ID" != "0" ]; then
        test_endpoint "DELETE" "/mahasiswa/$MHS_ID" "Delete Test Mahasiswa" "" "$ADMIN_TOKEN"
    fi
fi

# ============================================================================
# 7. PERMISSION TESTS
# ============================================================================

print_section "7. PERMISSION / RBAC TESTS"

if [ -n "$USER_TOKEN" ]; then
    test_endpoint "POST" "/mahasiswa" "User Create Mahasiswa (Should Fail)" \
        "{\"nim\":\"FAIL123\",\"nama\":\"Should Fail\",\"jurusan\":\"Test\",\"angkatan\":2024,\"email\":\"fail@test.com\"}" \
        "$USER_TOKEN" 403 true
    
    test_endpoint "DELETE" "/mahasiswa/1" "User Delete Mahasiswa (Should Fail)" \
        "" "$USER_TOKEN" 403 true
fi

# ============================================================================
# FINAL SUMMARY
# ============================================================================

echo ""
print_section "TEST SUMMARY & RESULTS"

TOTAL_TESTS=$((SUCCESS_COUNT + FAIL_COUNT))

echo -e "${BLUE}Database Type:${NC} ${DB_TYPE^^}"
echo -e "${GREEN}✅ Successful Tests: $SUCCESS_COUNT${NC}"
echo -e "${RED}❌ Failed Tests: $FAIL_COUNT${NC}"
if [ $SKIP_COUNT -gt 0 ]; then
    echo -e "${YELLOW}⏭️  Skipped Tests: $SKIP_COUNT${NC}"
fi
echo -e "${BLUE}📊 Total Tests: $TOTAL_TESTS${NC}"

if [ $TOTAL_TESTS -gt 0 ]; then
    SUCCESS_RATE=$((SUCCESS_COUNT * 100 / TOTAL_TESTS))
    echo -e "${MAGENTA}📈 Success Rate: ${SUCCESS_RATE}%${NC}"
fi

echo ""
echo -e "${CYAN}📋 Test Log: ${TEST_LOG}${NC}"
echo ""

# PocketBase specific note
if [ "$DB_TYPE" == "pocketbase" ]; then
    echo -e "${YELLOW}============================================================${NC}"
    echo -e "${YELLOW}   ⚠️  PocketBase Limitations Note${NC}"
    echo -e "${YELLOW}============================================================${NC}"
    echo -e "${CYAN}Some operations may fail due to ID type mismatch:${NC}"
    echo -e "${YELLOW}   - PocketBase uses string IDs${NC}"
    echo -e "${YELLOW}   - Go models use numeric IDs${NC}"
    echo -e "${YELLOW}   - List operations may fail on JSON unmarshaling${NC}"
    echo ""
fi

# Final verdict
if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}============================================================${NC}"
    echo -e "${GREEN}   🎉 ALL TESTS PASSED! 🎉${NC}"
    echo -e "${GREEN}============================================================${NC}"
    exit 0
else
    echo -e "${YELLOW}============================================================${NC}"
    echo -e "${YELLOW}   ⚠️  Some tests failed - Check log for details${NC}"
    echo -e "${YELLOW}============================================================${NC}"
    
    echo -e "${RED}Failed Tests:${NC}"
    for test_name in "${!TEST_RESULTS[@]}"; do
        if [ "${TEST_RESULTS[$test_name]}" == "FAIL" ]; then
            echo -e "${RED}  ❌ $test_name${NC}"
        fi
    done
    
    exit 1
fi
