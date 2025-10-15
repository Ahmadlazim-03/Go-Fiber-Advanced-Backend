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

type UserRepositoryPocketBase struct {
	baseURL string
	client  *http.Client
}

func NewUserRepository(baseURL string) *UserRepositoryPocketBase {
	return &UserRepositoryPocketBase{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// PocketBase uses built-in users collection with different structure
// We'll need to adapt our User model to PocketBase's auth collection

func (r *UserRepositoryPocketBase) Create(user *models.User) error {
	url := r.baseURL + "/api/collections/users/records"
	
	payload := map[string]interface{}{
		"username":       user.Username,
		"email":          user.Email,
		"password":       user.Password,
		"passwordConfirm": user.Password,
		"role":           user.Role,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := r.client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("create user failed (status %d): %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Set ID from response
	if id, ok := result["id"].(string); ok {
		// PocketBase uses string IDs, we'll need to handle this
		fmt.Printf("Created user with ID: %s\n", id)
	}

	return nil
}

func (r *UserRepositoryPocketBase) GetByID(id int) (*models.User, error) {
	url := fmt.Sprintf("%s/api/collections/users/records/%d", r.baseURL, id)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user failed (status %d)", resp.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepositoryPocketBase) GetByEmail(email string) (*models.User, error) {
	url := fmt.Sprintf("%s/api/collections/users/records?filter=(email='%s')", r.baseURL, email)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user failed (status %d)", resp.StatusCode)
	}

	var result struct {
		Items []models.User `json:"items"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	return &result.Items[0], nil
}

func (r *UserRepositoryPocketBase) GetByUsername(username string) (*models.User, error) {
	url := fmt.Sprintf("%s/api/collections/users/records?filter=(username='%s')", r.baseURL, username)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get user failed (status %d)", resp.StatusCode)
	}

	var result struct {
		Items []models.User `json:"items"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	return &result.Items[0], nil
}

func (r *UserRepositoryPocketBase) Update(user *models.User) error {
	url := fmt.Sprintf("%s/api/collections/users/records/%d", r.baseURL, user.ID)
	
	payload := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("update user failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (r *UserRepositoryPocketBase) Delete(id int) error {
	url := fmt.Sprintf("%s/api/collections/users/records/%d", r.baseURL, id)
	
	req, _ := http.NewRequest("DELETE", url, nil)
	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete user failed (status %d)", resp.StatusCode)
	}

	return nil
}

func (r *UserRepositoryPocketBase) GetAll() ([]models.User, error) {
	url := fmt.Sprintf("%s/api/collections/users/records?perPage=500", r.baseURL)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get users failed (status %d)", resp.StatusCode)
	}

	var result struct {
		Items []models.User `json:"items"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

func (r *UserRepositoryPocketBase) GetWithPagination(pagination *models.PaginationRequest) ([]models.User, int64, error) {
	page := pagination.Page
	if page < 1 {
		page = 1
	}
	
	url := fmt.Sprintf("%s/api/collections/users/records?perPage=%d&page=%d", 
		r.baseURL, pagination.Limit, page)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("get users failed (status %d)", resp.StatusCode)
	}

	var result struct {
		Items      []models.User `json:"items"`
		TotalItems int64         `json:"totalItems"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	return result.Items, result.TotalItems, nil
}

func (r *UserRepositoryPocketBase) Count() (int64, error) {
	url := fmt.Sprintf("%s/api/collections/users/records?perPage=1", r.baseURL)
	
	resp, err := r.client.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("count users failed (status %d)", resp.StatusCode)
	}

	var result struct {
		TotalItems int64 `json:"totalItems"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.TotalItems, nil
}
