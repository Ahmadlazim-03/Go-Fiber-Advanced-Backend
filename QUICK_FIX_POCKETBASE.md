# 🔧 Quick Fix: PocketBase 404 Authentication Error

## Problem
```
Error authenticating with PocketBase: authentication failed (status 404)
{"data":{},"message":"The requested resource wasn't found.","status":404}
```

## Root Cause
❌ **No admin account exists in PocketBase yet!**

The endpoint `/api/admins/auth-with-password` returns 404 because:
- PocketBase instance is running ✅
- But no admin account has been created yet ❌

## Solution (2 Steps)

### Step 1: Open PocketBase Admin Dashboard
```
https://pocketbase-production-521e.up.railway.app/_/
```
⚠️ **Important:** Must use `/_/` at the end!

### Step 2: Create Admin Account
When you open the URL, you'll see a setup screen. Fill in:
```
Email: admin@gmail.com
Password: 12345678
```
(These are from your `.env` file)

Click **"Create Admin"** or **"Set up"**

### Step 3: Restart Application
```powershell
go run main.go
```

## Expected Result After Fix

```
Connecting to PocketBase at: https://pocketbase-production-521e.up.railway.app
✓ PocketBase URL configured successfully
Note: PocketBase uses HTTP API - no persistent connection needed
Database type: pocketbase
Running PocketBase database migrations...
✓ Authenticated with PocketBase successfully  <-- Should work now!
✓ Using PocketBase built-in users collection (auth type)
Creating collection: mahasiswas...
✓ Collection mahasiswas ready
Creating collection: alumnis...
✓ Collection alumnis ready
Creating collection: pekerjaan_alumnis...
✓ Collection pekerjaan_alumnis ready
PocketBase database migrations completed successfully!
✓ All PocketBase repositories initialized successfully
Server running on http://localhost:8080
```

## Verification

### Test 1: Check PocketBase is Running
```powershell
Invoke-WebRequest -Uri "https://pocketbase-production-521e.up.railway.app/api/health"
```
Expected: `{"message":"API is healthy.","code":200,"data":{}}`

### Test 2: Check Admin Dashboard Accessible
Open in browser:
```
https://pocketbase-production-521e.up.railway.app/_/
```
Expected: Admin login page or setup page

### Test 3: Test Authentication (After Creating Admin)
```powershell
$body = @{
    identity = "admin@gmail.com"
    password = "12345678"
} | ConvertTo-Json

Invoke-WebRequest -Uri "https://pocketbase-production-521e.up.railway.app/api/admins/auth-with-password" `
    -Method POST `
    -ContentType "application/json" `
    -Body $body
```
Expected: Returns token and admin info

## Why This Happens

PocketBase is a self-hosted backend that requires manual admin setup:
1. ✅ Instance deployed on Railway
2. ❌ Admin account NOT created automatically
3. ❌ Cannot use Admin API without admin account
4. ✅ After creating admin → migrations work!

## Current Status

✅ PocketBase instance is running (verified with `/api/health`)
✅ Admin dashboard is accessible (`/_/`)
✅ Code is correct
❌ Just need to create admin account!

---

**Action Required:** Open `https://pocketbase-production-521e.up.railway.app/_/` now and create admin account!
