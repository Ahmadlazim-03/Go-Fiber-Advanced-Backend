#!/bin/bash

# Warna untuk output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   DUAL DATABASE TEST - Comprehensive  ${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# ==========================================
# TEST MONGODB
# ==========================================
echo -e "${MAGENTA}========================================${NC}"
echo -e "${MAGENTA}   TESTING MONGODB DATABASE            ${NC}"
echo -e "${MAGENTA}========================================${NC}"
echo ""

# Update .env to use MongoDB
echo -e "${YELLOW}Setting up MongoDB configuration...${NC}"
cat > .env << 'EOF'
# Database Configuration
DB_TYPE=mongodb

# MongoDB Configuration
MONGODB_URI=mongodb://mongo:pakgtnLdkcJlREVyWpuhiecIEQvnVOkh@caboose.proxy.rlwy.net:48828
MONGODB_DATABASE=railway

# JWT Secret
JWT_SECRET=your-secret-key-change-in-production
EOF

echo -e "${GREEN}✅ MongoDB configuration set${NC}"
echo ""

# Kill existing server
echo -e "${YELLOW}Stopping existing server...${NC}"
pkill -f "go run main.go" 2>/dev/null || true
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
sleep 2

# Start server with MongoDB
echo -e "${YELLOW}Starting server with MongoDB...${NC}"
cd /workspaces/Go-Fiber-Advanced-Backend
go run main.go > mongodb_server.log 2>&1 &
SERVER_PID=$!
echo -e "${GREEN}Server PID: $SERVER_PID${NC}"

# Wait for server to start
echo -e "${YELLOW}Waiting for server to start...${NC}"
for i in {1..20}; do
    if curl -s http://localhost:8080/ > /dev/null 2>&1; then
        echo -e "${GREEN}✅ MongoDB server is running${NC}"
        break
    fi
    echo -n "."
    sleep 1
done
echo ""

# Final check
if ! curl -s http://localhost:8080/ > /dev/null 2>&1; then
    echo -e "${RED}❌ MongoDB server failed to start${NC}"
    cat mongodb_server.log | tail -20
    exit 1
fi
echo ""

# Run MongoDB tests
echo -e "${YELLOW}Running MongoDB routes test...${NC}"
./test_all_routes_mongodb.sh > mongodb_test_result.log 2>&1

# Extract results
MONGODB_SUCCESS=$(grep "Successful Tests:" mongodb_test_result.log | grep -o '[0-9]*' | head -1)
MONGODB_FAIL=$(grep "Failed Tests:" mongodb_test_result.log | grep -o '[0-9]*' | head -1)
MONGODB_TOTAL=$(grep "Total Tests:" mongodb_test_result.log | grep -o '[0-9]*' | head -1)

echo -e "${BLUE}MongoDB Test Results:${NC}"
echo -e "  ${GREEN}Success: $MONGODB_SUCCESS${NC}"
echo -e "  ${RED}Failed: $MONGODB_FAIL${NC}"
echo -e "  ${CYAN}Total: $MONGODB_TOTAL${NC}"
echo ""

# Stop MongoDB server
echo -e "${YELLOW}Stopping MongoDB server...${NC}"
kill $SERVER_PID 2>/dev/null || true
sleep 2

# ==========================================
# TEST POSTGRESQL
# ==========================================
echo -e "${MAGENTA}========================================${NC}"
echo -e "${MAGENTA}   TESTING POSTGRESQL DATABASE         ${NC}"
echo -e "${MAGENTA}========================================${NC}"
echo ""

# Update .env to use PostgreSQL
echo -e "${YELLOW}Setting up PostgreSQL configuration...${NC}"
cat > .env << 'EOF'
# Database Configuration
DB_TYPE=postgres

# PostgreSQL Configuration
POSTGRES_HOST=autorack.proxy.rlwy.net
POSTGRES_PORT=35086
POSTGRES_USER=postgres
POSTGRES_PASSWORD=pNrNmMVIFkEOCndZvxPdhBtNjZhYRFze
POSTGRES_DB=railway

# JWT Secret
JWT_SECRET=your-secret-key-change-in-production
EOF

echo -e "${GREEN}✅ PostgreSQL configuration set${NC}"
echo ""

# Start server with PostgreSQL
echo -e "${YELLOW}Starting server with PostgreSQL...${NC}"
go run main.go > postgres_server.log 2>&1 &
SERVER_PID=$!
echo -e "${GREEN}Server PID: $SERVER_PID${NC}"

# Wait for server to start
echo -e "${YELLOW}Waiting for server to start...${NC}"
for i in {1..20}; do
    if curl -s http://localhost:8080/ > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PostgreSQL server is running${NC}"
        break
    fi
    echo -n "."
    sleep 1
done
echo ""

# Final check
if ! curl -s http://localhost:8080/ > /dev/null 2>&1; then
    echo -e "${RED}❌ PostgreSQL server failed to start${NC}"
    cat postgres_server.log | tail -20
    exit 1
fi
echo ""

# Run PostgreSQL tests
echo -e "${YELLOW}Running PostgreSQL routes test...${NC}"
./test_all_routes_mongodb.sh > postgres_test_result.log 2>&1

# Extract results
POSTGRES_SUCCESS=$(grep "Successful Tests:" postgres_test_result.log | grep -o '[0-9]*' | head -1)
POSTGRES_FAIL=$(grep "Failed Tests:" postgres_test_result.log | grep -o '[0-9]*' | head -1)
POSTGRES_TOTAL=$(grep "Total Tests:" postgres_test_result.log | grep -o '[0-9]*' | head -1)

echo -e "${BLUE}PostgreSQL Test Results:${NC}"
echo -e "  ${GREEN}Success: $POSTGRES_SUCCESS${NC}"
echo -e "  ${RED}Failed: $POSTGRES_FAIL${NC}"
echo -e "  ${CYAN}Total: $POSTGRES_TOTAL${NC}"
echo ""

# Stop PostgreSQL server
echo -e "${YELLOW}Stopping PostgreSQL server...${NC}"
kill $SERVER_PID 2>/dev/null || true
sleep 2

# ==========================================
# FINAL SUMMARY
# ==========================================
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}        DUAL DATABASE TEST SUMMARY     ${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

echo -e "${CYAN}MongoDB Results:${NC}"
echo -e "  ✅ Success: ${GREEN}$MONGODB_SUCCESS${NC}"
echo -e "  ❌ Failed:  ${RED}$MONGODB_FAIL${NC}"
echo -e "  📊 Total:   ${CYAN}$MONGODB_TOTAL${NC}"
echo -e "  📈 Success Rate: ${GREEN}$(awk "BEGIN {printf \"%.1f\", ($MONGODB_SUCCESS/$MONGODB_TOTAL)*100}")%${NC}"
echo ""

echo -e "${CYAN}PostgreSQL Results:${NC}"
echo -e "  ✅ Success: ${GREEN}$POSTGRES_SUCCESS${NC}"
echo -e "  ❌ Failed:  ${RED}$POSTGRES_FAIL${NC}"
echo -e "  📊 Total:   ${CYAN}$POSTGRES_TOTAL${NC}"
echo -e "  📈 Success Rate: ${GREEN}$(awk "BEGIN {printf \"%.1f\", ($POSTGRES_SUCCESS/$POSTGRES_TOTAL)*100}")%${NC}"
echo ""

# Compare results
if [ "$MONGODB_SUCCESS" -eq "$POSTGRES_SUCCESS" ]; then
    echo -e "${GREEN}✅ Both databases have the same success count!${NC}"
else
    echo -e "${YELLOW}⚠️  Different success counts between databases${NC}"
    if [ "$MONGODB_SUCCESS" -gt "$POSTGRES_SUCCESS" ]; then
        echo -e "  MongoDB has ${GREEN}$(($MONGODB_SUCCESS - $POSTGRES_SUCCESS))${NC} more successful tests"
    else
        echo -e "  PostgreSQL has ${GREEN}$(($POSTGRES_SUCCESS - $MONGODB_SUCCESS))${NC} more successful tests"
    fi
fi
echo ""

# Calculate total
TOTAL_SUCCESS=$(($MONGODB_SUCCESS + $POSTGRES_SUCCESS))
TOTAL_FAIL=$(($MONGODB_FAIL + $POSTGRES_FAIL))
TOTAL_TESTS=$(($MONGODB_TOTAL + $POSTGRES_TOTAL))

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}         OVERALL STATISTICS            ${NC}"
echo -e "${BLUE}========================================${NC}"
echo -e "  Total Tests Run:     ${CYAN}$TOTAL_TESTS${NC}"
echo -e "  Total Successful:    ${GREEN}$TOTAL_SUCCESS${NC}"
echo -e "  Total Failed:        ${RED}$TOTAL_FAIL${NC}"
echo -e "  Overall Success Rate: ${GREEN}$(awk "BEGIN {printf \"%.1f\", ($TOTAL_SUCCESS/$TOTAL_TESTS)*100}")%${NC}"
echo ""

if [ "$TOTAL_FAIL" -eq 0 ]; then
    echo -e "${GREEN}========================================${NC}"
    echo -e "${GREEN}   🎉 ALL TESTS PASSED! 🎉            ${NC}"
    echo -e "${GREEN}========================================${NC}"
else
    echo -e "${YELLOW}========================================${NC}"
    echo -e "${YELLOW}   ⚠️  Some tests failed. Review logs.${NC}"
    echo -e "${YELLOW}========================================${NC}"
fi

echo ""
echo -e "${CYAN}📄 Detailed logs:${NC}"
echo -e "  • MongoDB: mongodb_test_result.log"
echo -e "  • PostgreSQL: postgres_test_result.log"
echo -e "  • Server logs: mongodb_server.log, postgres_server.log"
echo ""

# Restore MongoDB as default
echo -e "${YELLOW}Restoring MongoDB as default database...${NC}"
cat > .env << 'EOF'
# Database Configuration
DB_TYPE=mongodb

# MongoDB Configuration
MONGODB_URI=mongodb://mongo:pakgtnLdkcJlREVyWpuhiecIEQvnVOkh@caboose.proxy.rlwy.net:48828
MONGODB_DATABASE=railway

# JWT Secret
JWT_SECRET=your-secret-key-change-in-production
EOF

echo -e "${GREEN}✅ Configuration restored to MongoDB${NC}"
echo ""
echo -e "${BLUE}Test completed!${NC}"
