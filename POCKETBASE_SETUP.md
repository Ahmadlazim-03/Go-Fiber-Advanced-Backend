# 🚀 PocketBase Setup Guide

## Step-by-Step Setup

### Step 1: Access PocketBase Admin Panel

Go to your PocketBase URL admin panel:
```
https://pocketbase-production-521e.up.railway.app/_/
```

**IMPORTANT:** PocketBase admin dashboard uses `/_/` endpoint!

### Step 2: Create Admin Account

**CRITICAL:** You MUST create an admin account first before running migrations!

1. When you first access the admin panel (`/_/`), you'll see a setup screen
2. Fill in the admin credentials (use credentials from your .env file):
   ```
   Email: admin@gmail.com
   Password: 12345678
   ```
3. Click "Create" or "Set up admin account"
4. You'll be automatically logged into the admin dashboard

### Step 3: Verify Admin Account

After creating the admin account, you should be able to login to the admin dashboard.

### Step 4: Run Application

Now you can run the application and migrations will work:
```bash
go run main.go
```

The application will:
1. ✅ Connect to PocketBase
2. ✅ Authenticate with admin credentials
3. ✅ Create collections automatically
4. ✅ Start the server

## If You See Authentication Error

Error message:
```
Error authenticating with PocketBase: authentication failed (status 404)
```

**Cause:** No admin account exists in PocketBase yet

**Solution:**
1. Go to `https://pocketbase-production-521e.up.railway.app/_/`
2. Create admin account with credentials from `.env`
3. Restart application

## Manual Collection Creation (Alternative)

If you prefer to create collections manually:

### 1. Login to Admin Panel
```
https://pocketbase-production-521e.up.railway.app/_/
```

### 2. Create Collections

Click "Collections" → "New Collection"

#### Users Collection
- **Name:** users
- **Type:** Auth (for authentication)
- Fields are auto-created by PocketBase

#### Mahasiswas Collection
- **Name:** mahasiswas
- **Type:** Base
- **Schema:**
  - `nim` (Text, Required, Max: 20)
  - `nama` (Text, Required, Max: 100)
  - `jurusan` (Text, Required, Max: 50)
  - `angkatan` (Number, Required)
  - `email` (Email, Required)

#### Alumnis Collection
- **Name:** alumnis
- **Type:** Base
- **Schema:**
  - `user_id` (Number, Required)
  - `nim` (Text, Required, Max: 20)
  - `nama` (Text, Required, Max: 100)
  - `jurusan` (Text, Required, Max: 50)
  - `angkatan` (Number, Required)
  - `tahun_lulus` (Number, Required)
  - `no_telepon` (Text, Max: 15)
  - `alamat` (Text)

#### Pekerjaan Alumnis Collection
- **Name:** pekerjaan_alumnis
- **Type:** Base
- **Schema:**
  - `alumni_id` (Number, Required)
  - `nama_perusahaan` (Text, Required, Max: 100)
  - `posisi_jabatan` (Text, Required, Max: 100)
  - `bidang_industri` (Text, Required, Max: 50)
  - `lokasi_kerja` (Text, Required, Max: 100)
  - `gaji_range` (Text, Max: 50)
  - `tanggal_mulai_kerja` (Date, Required)
  - `tanggal_selesai_kerja` (Date)
  - `status_pekerjaan` (Text, Max: 20)
  - `deskripsi_pekerjaan` (Text)
  - `deleted_at` (Date) - for soft delete

### 3. Set API Rules

For each collection, set API rules:

**Development (Open Access):**
```
List Rule: ""
View Rule: ""
Create Rule: ""
Update Rule: ""
Delete Rule: ""
```

**Production (Authenticated Only):**
```
List Rule: "@request.auth.id != ''"
View Rule: "@request.auth.id != ''"
Create Rule: "@request.auth.id != ''"
Update Rule: "@request.auth.id != ''"
Delete Rule: "@request.auth.id != ''"
```

## Troubleshooting

### 1. Cannot Access Admin Panel

**Problem:** `https://pocketbase-production-521e.up.railway.app/_/` returns 404

**Solutions:**
- Check if PocketBase instance is running on Railway
- Verify the URL is correct
- Check Railway deployment logs

### 2. Admin Already Exists

**Problem:** "Admin with this email already exists"

**Solutions:**
- Use the existing admin credentials
- Or reset admin password via Railway console:
  ```bash
  pocketbase admin update email@example.com newpassword
  ```

### 3. Migrations Still Fail After Creating Admin

**Problem:** Authentication still fails after creating admin

**Solutions:**
- Double-check email and password in `.env` match admin account
- Ensure no extra spaces in `.env` file
- Try logging in to admin panel manually to verify credentials
- Check if PocketBase version is compatible

### 4. Collections Not Created

**Problem:** Migrations run but collections don't appear

**Solutions:**
- Create collections manually via admin panel
- Check PocketBase logs for errors
- Verify admin token is valid

## Testing PocketBase Connection

### Test Admin Authentication
```bash
curl -X POST https://pocketbase-production-521e.up.railway.app/api/admins/auth-with-password \
  -H "Content-Type: application/json" \
  -d '{"identity":"ahmadlazim422@gmail.com","password":"Pembelajaranjarakjauh@123"}'
```

Expected response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "admin": {
    "id": "...",
    "email": "ahmadlazim422@gmail.com",
    ...
  }
}
```

### Test Collections
```bash
curl https://pocketbase-production-521e.up.railway.app/api/collections
```

## Environment Variables Reference

```env
# PocketBase Configuration
POCKETBASE_URL=https://pocketbase-production-521e.up.railway.app
POCKETBASE_ADMIN_EMAIL=ahmadlazim422@gmail.com
POCKETBASE_ADMIN_PASSWORD=Pembelajaranjarakjauh@123
```

**Important Notes:**
- ⚠️ Admin password must be at least 10 characters
- ⚠️ Admin email must be valid format
- ⚠️ No spaces before or after values in .env
- ⚠️ Password is case-sensitive

## Success Indicators

When everything is set up correctly, you should see:

```
Connecting to PocketBase at: https://pocketbase-production-521e.up.railway.app
✓ PocketBase URL configured successfully
Note: PocketBase uses HTTP API - no persistent connection needed
Database type: pocketbase
Running PocketBase database migrations...
✓ Authenticated with PocketBase successfully
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

## Next Steps After Setup

1. ✅ Test API endpoints
2. ✅ Create test data
3. ✅ Configure API rules for production
4. ✅ Set up backup strategy
5. ✅ Monitor PocketBase logs

## Getting Help

If you still have issues:
1. Check Railway deployment logs
2. Review PocketBase documentation: https://pocketbase.io/docs/
3. Verify all environment variables are correct
4. Test admin login manually in browser

---

**Remember:** Admin account MUST be created via admin panel (`/_/`) before running migrations!
