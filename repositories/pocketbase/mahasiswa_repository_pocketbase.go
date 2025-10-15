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

type MahasiswaRepositoryPocketBase struct {
	baseURL string
	client  *http.Client
}

func NewMahasiswaRepository(baseURL string) *MahasiswaRepositoryPocketBase {
	return &MahasiswaRepositoryPocketBase{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *MahasiswaRepositoryPocketBase) Create(mahasiswa *models.Mahasiswa) error {
	url := r.baseURL + "/api/collections/mahasiswas/records"
	
	payload := map[string]interface{}{
		"nim":      mahasiswa.NIM,
		"nama":     mahasiswa.Nama,
		"jurusan":  mahasiswa.Jurusan,
		"angkatan": mahasiswa.Angkatan,
		"email":    mahasiswa.Email,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := r.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create mahasiswa: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create mahasiswa failed (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Set ID from response
	if id, ok := result["id"].(string); ok {
		fmt.Printf("Created mahasiswa with ID: %s\n", id)
	}

	return nil
}

func (r *MahasiswaRepositoryPocketBase) GetByID(id uint) (*models.Mahasiswa, error) {
	url := fmt.Sprintf("%s/api/collections/mahasiswas/records/%d", r.baseURL, id)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get mahasiswa: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get mahasiswa failed (status %d)", resp.StatusCode)
	}

	var mahasiswa models.Mahasiswa
	if err := json.NewDecoder(resp.Body).Decode(&mahasiswa); err != nil {
		return nil, err
	}

	return &mahasiswa, nil
}

func (r *MahasiswaRepositoryPocketBase) Update(mahasiswa *models.Mahasiswa) error {
	url := fmt.Sprintf("%s/api/collections/mahasiswas/records/%d", r.baseURL, mahasiswa.ID)
	
	payload := map[string]interface{}{
		"nim":      mahasiswa.NIM,
		"nama":     mahasiswa.Nama,
		"jurusan":  mahasiswa.Jurusan,
		"angkatan": mahasiswa.Angkatan,
		"email":    mahasiswa.Email,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update mahasiswa: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update mahasiswa failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *MahasiswaRepositoryPocketBase) Delete(id uint) error {
	url := fmt.Sprintf("%s/api/collections/mahasiswas/records/%d", r.baseURL, id)
	
	req, _ := http.NewRequest("DELETE", url, nil)
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete mahasiswa: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete mahasiswa failed (status %d)", resp.StatusCode)
	}

	return nil
}

func (r *MahasiswaRepositoryPocketBase) GetAll() ([]models.Mahasiswa, error) {
	url := fmt.Sprintf("%s/api/collections/mahasiswas/records?perPage=500", r.baseURL)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get mahasiswas: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get mahasiswas failed (status %d)", resp.StatusCode)
	}

	var result struct {
		Items []models.Mahasiswa `json:"items"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (r *MahasiswaRepositoryPocketBase) GetWithPagination(pagination *models.PaginationRequest) ([]models.Mahasiswa, int64, error) {
	page := pagination.Page
	if page < 1 {
		page = 1
	}
	
	url := fmt.Sprintf("%s/api/collections/mahasiswas/records?perPage=%d&page=%d", 
		r.baseURL, pagination.Limit, page)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get mahasiswas: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("get mahasiswas failed (status %d)", resp.StatusCode)
	}

	var result struct {
		Items      []models.Mahasiswa `json:"items"`
		TotalItems int64              `json:"totalItems"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	return result.Items, result.TotalItems, nil
}

func (r *MahasiswaRepositoryPocketBase) Count() (int64, error) {
	url := fmt.Sprintf("%s/api/collections/mahasiswas/records?perPage=1", r.baseURL)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to count mahasiswas: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("count mahasiswas failed (status %d)", resp.StatusCode)
	}

	var result struct {
		TotalItems int64 `json:"totalItems"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.TotalItems, nil
}
