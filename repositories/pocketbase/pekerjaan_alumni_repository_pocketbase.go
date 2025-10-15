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

type PekerjaanAlumniRepositoryPocketBase struct {
	baseURL string
	client  *http.Client
}

func NewPekerjaanAlumniRepository(baseURL string) *PekerjaanAlumniRepositoryPocketBase {
	return &PekerjaanAlumniRepositoryPocketBase{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *PekerjaanAlumniRepositoryPocketBase) Create(pekerjaan *models.PekerjaanAlumni) error {
	url := r.baseURL + "/api/collections/pekerjaan_alumnis/records"
	
	payload := map[string]interface{}{
		"alumni_id":               pekerjaan.AlumniID,
		"nama_perusahaan":         pekerjaan.NamaPerusahaan,
		"posisi_jabatan":          pekerjaan.PosisiJabatan,
		"bidang_industri":         pekerjaan.BidangIndustri,
		"lokasi_kerja":            pekerjaan.LokasiKerja,
		"gaji_range":              pekerjaan.GajiRange,
		"tanggal_mulai_kerja":     pekerjaan.TanggalMulaiKerja,
		"tanggal_selesai_kerja":   pekerjaan.TanggalSelesaiKerja,
		"status_pekerjaan":        pekerjaan.StatusPekerjaan,
		"deskripsi_pekerjaan":     pekerjaan.DeskripsiPekerjaan,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := r.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create pekerjaan: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create pekerjaan failed (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Set ID from response
	if id, ok := result["id"].(string); ok {
		fmt.Printf("Created pekerjaan with ID: %s\n", id)
	}

	return nil
}

func (r *PekerjaanAlumniRepositoryPocketBase) GetByID(id uint) (*models.PekerjaanAlumni, error) {
	url := fmt.Sprintf("%s/api/collections/pekerjaan_alumnis/records/%d", r.baseURL, id)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get pekerjaan: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get pekerjaan failed (status %d)", resp.StatusCode)
	}

	var pekerjaan models.PekerjaanAlumni
	if err := json.NewDecoder(resp.Body).Decode(&pekerjaan); err != nil {
		return nil, err
	}

	return &pekerjaan, nil
}

func (r *PekerjaanAlumniRepositoryPocketBase) GetByAlumniID(alumniID uint) ([]models.PekerjaanAlumni, error) {
	url := fmt.Sprintf("%s/api/collections/pekerjaan_alumnis/records?filter=(alumni_id=%d)", r.baseURL, alumniID)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get pekerjaan by alumni_id: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get pekerjaan failed (status %d)", resp.StatusCode)
	}

	var result struct {
		Items []models.PekerjaanAlumni `json:"items"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (r *PekerjaanAlumniRepositoryPocketBase) GetByUserID(userID int) ([]models.PekerjaanAlumni, error) {
	// First, get alumni by user_id
	alumniURL := fmt.Sprintf("%s/api/collections/alumnis/records?filter=(user_id=%d)", r.baseURL, userID)
	
	resp, err := r.client.Get(alumniURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get alumni: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get alumni failed (status %d)", resp.StatusCode)
	}

	var alumniResult struct {
		Items []models.Alumni `json:"items"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&alumniResult); err != nil {
		return nil, err
	}

	if len(alumniResult.Items) == 0 {
		return []models.PekerjaanAlumni{}, nil
	}

	// Then get pekerjaan by alumni_id
	return r.GetByAlumniID(alumniResult.Items[0].ID)
}

func (r *PekerjaanAlumniRepositoryPocketBase) Update(pekerjaan *models.PekerjaanAlumni) error {
	url := fmt.Sprintf("%s/api/collections/pekerjaan_alumnis/records/%d", r.baseURL, pekerjaan.ID)
	
	payload := map[string]interface{}{
		"alumni_id":               pekerjaan.AlumniID,
		"nama_perusahaan":         pekerjaan.NamaPerusahaan,
		"posisi_jabatan":          pekerjaan.PosisiJabatan,
		"bidang_industri":         pekerjaan.BidangIndustri,
		"lokasi_kerja":            pekerjaan.LokasiKerja,
		"gaji_range":              pekerjaan.GajiRange,
		"tanggal_mulai_kerja":     pekerjaan.TanggalMulaiKerja,
		"tanggal_selesai_kerja":   pekerjaan.TanggalSelesaiKerja,
		"status_pekerjaan":        pekerjaan.StatusPekerjaan,
		"deskripsi_pekerjaan":     pekerjaan.DeskripsiPekerjaan,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update pekerjaan: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update pekerjaan failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *PekerjaanAlumniRepositoryPocketBase) Delete(id uint) error {
	url := fmt.Sprintf("%s/api/collections/pekerjaan_alumnis/records/%d", r.baseURL, id)
	
	req, _ := http.NewRequest("DELETE", url, nil)
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete pekerjaan: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete pekerjaan failed (status %d)", resp.StatusCode)
	}

	return nil
}

// Soft delete in PocketBase - using deleted_at field
func (r *PekerjaanAlumniRepositoryPocketBase) SoftDelete(id uint) error {
	url := fmt.Sprintf("%s/api/collections/pekerjaan_alumnis/records/%d", r.baseURL, id)
	
	now := time.Now()
	payload := map[string]interface{}{
		"deleted_at": now.Format(time.RFC3339),
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to soft delete pekerjaan: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("soft delete pekerjaan failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *PekerjaanAlumniRepositoryPocketBase) SoftDeleteByAlumniID(alumniID uint) error {
	// Get all pekerjaan for this alumni
	pekerjaans, err := r.GetByAlumniID(alumniID)
	if err != nil {
		return err
	}

	// Soft delete each one
	for _, p := range pekerjaans {
		if err := r.SoftDelete(p.ID); err != nil {
			return err
		}
	}

	return nil
}

func (r *PekerjaanAlumniRepositoryPocketBase) Restore(id uint) error {
	url := fmt.Sprintf("%s/api/collections/pekerjaan_alumnis/records/%d", r.baseURL, id)
	
	payload := map[string]interface{}{
		"deleted_at": nil,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to restore pekerjaan: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("restore pekerjaan failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *PekerjaanAlumniRepositoryPocketBase) GetDeleted() ([]models.PekerjaanAlumni, error) {
	url := fmt.Sprintf("%s/api/collections/pekerjaan_alumnis/records?filter=(deleted_at!=null)", r.baseURL)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get deleted pekerjaan: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get deleted pekerjaan failed (status %d)", resp.StatusCode)
	}

	var result struct {
		Items []models.PekerjaanAlumni `json:"items"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (r *PekerjaanAlumniRepositoryPocketBase) GetDeletedByUserID(userID int) ([]models.PekerjaanAlumni, error) {
	// First, get alumni by user_id
	alumniURL := fmt.Sprintf("%s/api/collections/alumnis/records?filter=(user_id=%d)", r.baseURL, userID)
	
	resp, err := r.client.Get(alumniURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get alumni: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get alumni failed (status %d)", resp.StatusCode)
	}

	var alumniResult struct {
		Items []models.Alumni `json:"items"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&alumniResult); err != nil {
		return nil, err
	}

	if len(alumniResult.Items) == 0 {
		return []models.PekerjaanAlumni{}, nil
	}

	// Then get deleted pekerjaan by alumni_id
	url := fmt.Sprintf("%s/api/collections/pekerjaan_alumnis/records?filter=(alumni_id=%d&&deleted_at!=null)", 
		r.baseURL, alumniResult.Items[0].ID)
	
	resp2, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get deleted pekerjaan: %v", err)
	}
	defer resp2.Body.Close()

	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get deleted pekerjaan failed (status %d)", resp2.StatusCode)
	}

	var result struct {
		Items []models.PekerjaanAlumni `json:"items"`
	}
	
	if err := json.NewDecoder(resp2.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (r *PekerjaanAlumniRepositoryPocketBase) GetAll() ([]models.PekerjaanAlumni, error) {
	url := fmt.Sprintf("%s/api/collections/pekerjaan_alumnis/records?perPage=500&filter=(deleted_at=null||deleted_at='')", r.baseURL)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get pekerjaans: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get pekerjaans failed (status %d)", resp.StatusCode)
	}

	var result struct {
		Items []models.PekerjaanAlumni `json:"items"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (r *PekerjaanAlumniRepositoryPocketBase) GetWithPagination(pagination *models.PaginationRequest) ([]models.PekerjaanAlumni, int64, error) {
	page := pagination.Page
	if page < 1 {
		page = 1
	}
	
	url := fmt.Sprintf("%s/api/collections/pekerjaan_alumnis/records?perPage=%d&page=%d&filter=(deleted_at=null||deleted_at='')", 
		r.baseURL, pagination.Limit, page)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get pekerjaans: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("get pekerjaans failed (status %d)", resp.StatusCode)
	}

	var result struct {
		Items      []models.PekerjaanAlumni `json:"items"`
		TotalItems int64                    `json:"totalItems"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	return result.Items, result.TotalItems, nil
}

func (r *PekerjaanAlumniRepositoryPocketBase) Count() (int64, error) {
	url := fmt.Sprintf("%s/api/collections/pekerjaan_alumnis/records?perPage=1&filter=(deleted_at=null||deleted_at='')", r.baseURL)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to count pekerjaans: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("count pekerjaans failed (status %d)", resp.StatusCode)
	}

	var result struct {
		TotalItems int64 `json:"totalItems"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.TotalItems, nil
}

func (r *PekerjaanAlumniRepositoryPocketBase) GetAlumniCountByCompany(namaPerusahaan string) (int64, error) {
	url := fmt.Sprintf("%s/api/collections/pekerjaan_alumnis/records?perPage=1&filter=(nama_perusahaan='%s'&&(deleted_at=null||deleted_at=''))", 
		r.baseURL, namaPerusahaan)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to count alumni by company: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("count alumni by company failed (status %d)", resp.StatusCode)
	}

	var result struct {
		TotalItems int64 `json:"totalItems"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.TotalItems, nil
}
