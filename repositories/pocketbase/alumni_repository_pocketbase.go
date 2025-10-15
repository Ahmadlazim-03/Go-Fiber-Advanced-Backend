package pocketbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"modul4crud/models"
	"net/http"
	"time"
)

type AlumniRepositoryPocketBase struct {
	baseURL string
	client  *http.Client
}

func NewAlumniRepository(baseURL string) *AlumniRepositoryPocketBase {
	return &AlumniRepositoryPocketBase{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *AlumniRepositoryPocketBase) Create(alumni *models.Alumni) error {
	url := r.baseURL + "/api/collections/alumnis/records"
	
	payload := map[string]interface{}{
		"user_id":     alumni.UserID,
		"nim":         alumni.NIM,
		"nama":        alumni.Nama,
		"jurusan":     alumni.Jurusan,
		"angkatan":    alumni.Angkatan,
		"tahun_lulus": alumni.TahunLulus,
		"no_telepon":  alumni.NoTelepon,
		"alamat":      alumni.Alamat,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := r.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create alumni: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create alumni failed (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Set ID from response
	if id, ok := result["id"].(string); ok {
		fmt.Printf("Created alumni with ID: %s\n", id)
	}

	return nil
}

func (r *AlumniRepositoryPocketBase) GetByID(id uint) (*models.Alumni, error) {
	url := fmt.Sprintf("%s/api/collections/alumnis/records/%d", r.baseURL, id)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get alumni: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get alumni failed (status %d)", resp.StatusCode)
	}

	var alumni models.Alumni
	if err := json.NewDecoder(resp.Body).Decode(&alumni); err != nil {
		return nil, err
	}

	return &alumni, nil
}

func (r *AlumniRepositoryPocketBase) GetByUserID(userID int) (*models.Alumni, error) {
	url := fmt.Sprintf("%s/api/collections/alumnis/records?filter=(user_id=%d)", r.baseURL, userID)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get alumni by user_id: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get alumni failed (status %d)", resp.StatusCode)
	}

	var result struct {
		Items []models.Alumni `json:"items"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	return &result.Items[0], nil
}

func (r *AlumniRepositoryPocketBase) Update(alumni *models.Alumni) error {
	url := fmt.Sprintf("%s/api/collections/alumnis/records/%d", r.baseURL, alumni.ID)
	
	payload := map[string]interface{}{
		"user_id":     alumni.UserID,
		"nim":         alumni.NIM,
		"nama":        alumni.Nama,
		"jurusan":     alumni.Jurusan,
		"angkatan":    alumni.Angkatan,
		"tahun_lulus": alumni.TahunLulus,
		"no_telepon":  alumni.NoTelepon,
		"alamat":      alumni.Alamat,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update alumni: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update alumni failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *AlumniRepositoryPocketBase) Delete(id uint) error {
	url := fmt.Sprintf("%s/api/collections/alumnis/records/%d", r.baseURL, id)
	
	req, _ := http.NewRequest("DELETE", url, nil)
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete alumni: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete alumni failed (status %d)", resp.StatusCode)
	}

	return nil
}

func (r *AlumniRepositoryPocketBase) GetAll() ([]models.Alumni, error) {
	url := fmt.Sprintf("%s/api/collections/alumnis/records?perPage=500&expand=user", r.baseURL)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get alumnis: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get alumnis failed (status %d)", resp.StatusCode)
	}

	var result struct {
		Items []models.Alumni `json:"items"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (r *AlumniRepositoryPocketBase) GetWithPagination(pagination *models.PaginationRequest) ([]models.Alumni, int64, error) {
	page := pagination.Page
	if page < 1 {
		page = 1
	}
	
	url := fmt.Sprintf("%s/api/collections/alumnis/records?perPage=%d&page=%d&expand=user", 
		r.baseURL, pagination.Limit, page)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get alumnis: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("get alumnis failed (status %d)", resp.StatusCode)
	}

	var result struct {
		Items      []models.Alumni `json:"items"`
		TotalItems int64           `json:"totalItems"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	return result.Items, result.TotalItems, nil
}

func (r *AlumniRepositoryPocketBase) Count() (int64, error) {
	url := fmt.Sprintf("%s/api/collections/alumnis/records?perPage=1", r.baseURL)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to count alumnis: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("count alumnis failed (status %d)", resp.StatusCode)
	}

	var result struct {
		TotalItems int64 `json:"totalItems"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.TotalItems, nil
}
