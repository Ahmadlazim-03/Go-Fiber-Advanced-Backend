package migration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"modul4crud/database"
	"net/http"
	"os"
	"time"
)

// PocketBase collection schema definitions
type PBCollection struct {
	Name       string                   `json:"name"`
	Type       string                   `json:"type"` // base, auth, view
	Schema     []PBField                `json:"schema"`
	System     bool                     `json:"system"`
	ListRule   *string                  `json:"listRule"`
	ViewRule   *string                  `json:"viewRule"`
	CreateRule *string                  `json:"createRule"`
	UpdateRule *string                  `json:"updateRule"`
	DeleteRule *string                  `json:"deleteRule"`
	Options    map[string]interface{}   `json:"options,omitempty"`
}

type PBField struct {
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	Required bool                   `json:"required"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

// RunPocketBaseMigrations membuat collections di PocketBase jika belum ada
func RunPocketBaseMigrations() {
	if !database.IsPocketBase() {
		log.Println("Skipping PocketBase migrations (not using PocketBase)")
		return
	}

	log.Println("Running PocketBase database migrations...")

	// Try to authenticate as admin
	token, err := authenticatePocketBase()
	if err != nil {
		log.Printf("⚠️  Warning: Could not authenticate with PocketBase admin: %v", err)
		log.Println("⚠️  Please create an admin account first:")
		log.Printf("⚠️  1. Go to: %s/_/", database.PocketBaseURL)
		log.Println("⚠️  2. Create an admin account with the credentials in .env")
		log.Println("⚠️  3. Restart the application")
		log.Println("⚠️  Collections will need to be created manually or via admin panel")
		log.Println("PocketBase migrations skipped - admin authentication required")
		return
	}

	// Create collections
	createUsersCollection(token)
	createMahasiswasCollection(token)
	createAlumnisCollection(token)
	createPekerjaanAlumnisCollection(token)

	log.Println("PocketBase database migrations completed successfully!")
}

// authenticatePocketBase authenticates with PocketBase admin API
func authenticatePocketBase() (string, error) {
	// PocketBase v0.20+ uses _superusers collection instead of admins
	url := database.PocketBaseURL + "/api/collections/_superusers/auth-with-password"
	
	email := os.Getenv("POCKETBASE_ADMIN_EMAIL")
	password := os.Getenv("POCKETBASE_ADMIN_PASSWORD")

	if email == "" || password == "" {
		return "", fmt.Errorf("PocketBase admin credentials not set")
	}

	payload := map[string]string{
		"identity": email,
		"password": password,
	}

	jsonData, _ := json.Marshal(payload)
	
	// Use longer timeout for Railway (cloud hosting can be slower)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to authenticate: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("authentication failed (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %v", err)
	}

	token, ok := result["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	log.Println("✓ Authenticated with PocketBase successfully")
	return token, nil
}

// Helper function to create or update collection
func createOrUpdateCollection(token string, collection PBCollection) error {
	url := database.PocketBaseURL + "/api/collections"
	
	// Check if collection exists (with longer timeout for Railway)
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("GET", url+"/"+collection.Name, nil)
	req.Header.Set("Authorization", token)
	
	resp, err := client.Do(req)
	collectionExists := err == nil && resp.StatusCode == http.StatusOK
	if resp != nil {
		resp.Body.Close()
	}

	jsonData, _ := json.Marshal(collection)

	var method string
	var endpoint string
	if collectionExists {
		method = "PATCH"
		endpoint = url + "/" + collection.Name
		log.Printf("Updating collection: %s with schema fields: %d", collection.Name, len(collection.Schema))
	} else {
		method = "POST"
		endpoint = url
		log.Printf("Creating collection: %s with schema fields: %d", collection.Name, len(collection.Schema))
	}
	
	// Debug: print the schema being sent
	log.Printf("Schema payload: %s", string(jsonData))

	req, _ = http.NewRequest(method, endpoint, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to create/update collection: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed (status %d): %s", resp.StatusCode, string(body))
	}

	log.Printf("✓ Collection %s ready - Response: %s", collection.Name, string(body))
	return nil
}

// createUsersCollection updates users collection with custom fields
func createUsersCollection(token string) {
	// PocketBase built-in users collection (auth type) exists by default
	// We need to add our custom fields: username, role, is_active
	
	collection := PBCollection{
		Name: "users",
		Type: "auth", // Keep as auth type - don't change!
		Schema: []PBField{
			// Add custom fields (built-in fields like email, password are automatic)
			{Name: "username", Type: "text", Required: true, Options: map[string]interface{}{"min": 3, "max": 50}},
			{Name: "role", Type: "text", Required: true, Options: map[string]interface{}{"min": 1, "max": 20}},
			{Name: "is_active", Type: "bool", Required: false},
		},
		ListRule:   stringPtr(""),
		ViewRule:   stringPtr(""),
		CreateRule: stringPtr(""),
		UpdateRule: stringPtr(""),
		DeleteRule: stringPtr(""),
	}

	if err := createOrUpdateCollection(token, collection); err != nil {
		log.Printf("⚠️  Warning: Could not update users collection: %v", err)
		log.Println("   Please add custom fields manually in PocketBase admin:")
		log.Println("   - username (text, required)")
		log.Println("   - role (text, required)")
		log.Println("   - is_active (bool)")
	} else {
		log.Println("✓ Users collection updated with custom fields")
	}
}

// createMahasiswasCollection creates mahasiswas collection
func createMahasiswasCollection(token string) {
	collection := PBCollection{
		Name: "mahasiswas",
		Type: "base",
		Schema: []PBField{
			{Name: "nim", Type: "text", Required: true, Options: map[string]interface{}{"min": 1, "max": 20}},
			{Name: "nama", Type: "text", Required: true, Options: map[string]interface{}{"min": 1, "max": 100}},
			{Name: "jurusan", Type: "text", Required: true, Options: map[string]interface{}{"min": 1, "max": 50}},
			{Name: "angkatan", Type: "number", Required: true},
			{Name: "email", Type: "email", Required: true},
		},
		ListRule:   stringPtr(""),
		ViewRule:   stringPtr(""),
		CreateRule: stringPtr(""),
		UpdateRule: stringPtr(""),
		DeleteRule: stringPtr(""),
	}

	if err := createOrUpdateCollection(token, collection); err != nil {
		log.Printf("Error with mahasiswas collection: %v", err)
	}
}

// createAlumnisCollection creates alumnis collection
func createAlumnisCollection(token string) {
	collection := PBCollection{
		Name: "alumnis",
		Type: "base",
		Schema: []PBField{
			{Name: "user_id", Type: "number", Required: true},
			{Name: "nim", Type: "text", Required: true, Options: map[string]interface{}{"min": 1, "max": 20}},
			{Name: "nama", Type: "text", Required: true, Options: map[string]interface{}{"min": 1, "max": 100}},
			{Name: "jurusan", Type: "text", Required: true, Options: map[string]interface{}{"min": 1, "max": 50}},
			{Name: "angkatan", Type: "number", Required: true},
			{Name: "tahun_lulus", Type: "number", Required: true},
			{Name: "no_telepon", Type: "text", Required: false, Options: map[string]interface{}{"max": 15}},
			{Name: "alamat", Type: "text", Required: false},
		},
		ListRule:   stringPtr(""),
		ViewRule:   stringPtr(""),
		CreateRule: stringPtr(""),
		UpdateRule: stringPtr(""),
		DeleteRule: stringPtr(""),
	}

	if err := createOrUpdateCollection(token, collection); err != nil {
		log.Printf("Error with alumnis collection: %v", err)
	}
}

// createPekerjaanAlumnisCollection creates pekerjaan_alumnis collection
func createPekerjaanAlumnisCollection(token string) {
	collection := PBCollection{
		Name: "pekerjaan_alumnis",
		Type: "base",
		Schema: []PBField{
			{Name: "alumni_id", Type: "number", Required: true},
			{Name: "nama_perusahaan", Type: "text", Required: true, Options: map[string]interface{}{"max": 100}},
			{Name: "posisi_jabatan", Type: "text", Required: true, Options: map[string]interface{}{"max": 100}},
			{Name: "bidang_industri", Type: "text", Required: true, Options: map[string]interface{}{"max": 50}},
			{Name: "lokasi_kerja", Type: "text", Required: true, Options: map[string]interface{}{"max": 100}},
			{Name: "gaji_range", Type: "text", Required: false, Options: map[string]interface{}{"max": 50}},
			{Name: "tanggal_mulai_kerja", Type: "date", Required: true},
			{Name: "tanggal_selesai_kerja", Type: "date", Required: false},
			{Name: "status_pekerjaan", Type: "text", Required: false, Options: map[string]interface{}{"max": 20}},
			{Name: "deskripsi_pekerjaan", Type: "text", Required: false},
		},
		ListRule:   stringPtr(""),
		ViewRule:   stringPtr(""),
		CreateRule: stringPtr(""),
		UpdateRule: stringPtr(""),
		DeleteRule: stringPtr(""),
	}

	if err := createOrUpdateCollection(token, collection); err != nil {
		log.Printf("Error with pekerjaan_alumnis collection: %v", err)
	}
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}
