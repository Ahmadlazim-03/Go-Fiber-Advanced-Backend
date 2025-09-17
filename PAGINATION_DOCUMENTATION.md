# Documentation: Pagination, Sorting & Search Features

## Overview

This documentation covers the implementation of pagination, sorting, and search features in the CRUD Go Fiber PostgreSQL application. These features provide efficient data management capabilities for Mahasiswa (Students), Alumni, and Pekerjaan Alumni (Alumni Jobs) sections.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Backend Implementation](#backend-implementation)
3. [Frontend Implementation](#frontend-implementation)
4. [API Endpoints](#api-endpoints)
5. [Usage Examples](#usage-examples)
6. [Testing Guide](#testing-guide)
7. [Troubleshooting](#troubleshooting)

---

## Architecture Overview

The pagination system follows a clean architecture pattern with clear separation of concerns:

```
Frontend (JavaScript/HTML) 
    ↓ AJAX Requests
Service Layer (Go Fiber)
    ↓ Business Logic
Repository Layer (GORM)
    ↓ Database Queries
PostgreSQL Database
```

### Key Components

- **PaginationRequest**: Handles incoming pagination parameters
- **PaginationResponse**: Standardized response format
- **Search**: ILIKE-based search across multiple fields
- **Sorting**: Dynamic ORDER BY with ASC/DESC options
- **UI Controls**: Bootstrap-based responsive interface

---

## Backend Implementation

### 1. Pagination Models (`models/pagination.go`)

#### PaginationRequest Struct
```go
type PaginationRequest struct {
    Page      int    `query:"page" json:"page"`
    Limit     int    `query:"limit" json:"limit"`
    Search    string `query:"search" json:"search"`
    SortBy    string `query:"sort_by" json:"sort_by"`
    SortOrder string `query:"sort_order" json:"sort_order"`
}
```

**Methods:**
- `SetDefaults()`: Sets default values (Page: 1, Limit: 10, SortOrder: "asc")
- `GetOffset()`: Calculates database offset based on page and limit
- `ValidateSortOrder()`: Ensures sort order is either "asc" or "desc"

#### PaginationResponse Struct
```go
type PaginationResponse struct {
    Data       interface{} `json:"data"`
    TotalItems int64       `json:"total_items"`
    TotalPages int         `json:"total_pages"`
    Page       int         `json:"page"`
    Limit      int         `json:"limit"`
    HasNext    bool        `json:"has_next"`
    HasPrev    bool        `json:"has_prev"`
}
```

### 2. Repository Layer Updates

All repositories now include the `GetWithPagination` method:

```go
GetWithPagination(pagination *PaginationRequest) ([]Model, int64, error)
```

#### Example Implementation (Mahasiswa Repository)
```go
func (r *mahasiswaRepository) GetWithPagination(pagination *PaginationRequest) ([]Mahasiswa, int64, error) {
    var mahasiswas []Mahasiswa
    var total int64
    
    // Build base query
    query := r.db
    
    // Apply search if provided
    if pagination.Search != "" {
        searchTerm := "%" + pagination.Search + "%"
        query = query.Where(
            "nim ILIKE ? OR nama ILIKE ? OR jurusan ILIKE ? OR email ILIKE ?",
            searchTerm, searchTerm, searchTerm, searchTerm,
        )
    }
    
    // Count total items
    if err := query.Model(&Mahasiswa{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Apply sorting
    if pagination.SortBy != "" {
        orderClause := pagination.SortBy + " " + pagination.SortOrder
        query = query.Order(orderClause)
    }
    
    // Apply pagination
    offset := pagination.GetOffset()
    if err := query.Limit(pagination.Limit).Offset(offset).Find(&mahasiswas).Error; err != nil {
        return nil, 0, err
    }
    
    return mahasiswas, total, nil
}
```

### 3. Service Layer Updates

Services now support both paginated and legacy endpoints:

#### Paginated Endpoints
- `GET /api/mahasiswa` - With pagination parameters
- `GET /api/alumni` - With pagination parameters  
- `GET /api/pekerjaan-alumni` - With pagination parameters

#### Legacy Endpoints (Backward Compatibility)
- `GET /api/mahasiswa/all` - Returns all records
- `GET /api/alumni/all` - Returns all records
- `GET /api/pekerjaan-alumni/all` - Returns all records

#### Query Parameter Parsing
```go
func (s *mahasiswaService) GetMahasiswaWithPagination(c *fiber.Ctx) error {
    pagination := &models.PaginationRequest{}
    
    // Parse query parameters
    if err := c.QueryParser(pagination); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid pagination parameters"})
    }
    
    pagination.SetDefaults()
    pagination.ValidateSortOrder()
    
    // Get paginated data
    mahasiswas, total, err := s.mahasiswaRepo.GetWithPagination(pagination)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch mahasiswa"})
    }
    
    // Create response
    response := models.NewPaginationResponse(mahasiswas, total, pagination.Page, pagination.Limit)
    return c.JSON(response)
}
```

---

## Frontend Implementation

### 1. UI Components

Each section includes:

#### Search Controls
```html
<div class="row g-3">
    <div class="col-md-3">
        <input type="text" class="form-control" id="mahasiswaSearchNim" placeholder="Cari NIM...">
    </div>
    <div class="col-md-3">
        <input type="text" class="form-control" id="mahasiswaSearchNama" placeholder="Cari Nama...">
    </div>
    <!-- Additional search fields -->
</div>
```

#### Sort Controls
```html
<div class="col-md-3">
    <select class="form-select" id="mahasiswaSortBy">
        <option value="">Urutkan berdasarkan...</option>
        <option value="nim">NIM</option>
        <option value="nama">Nama</option>
        <option value="jurusan">Jurusan</option>
        <option value="tahun_masuk">Tahun Masuk</option>
    </select>
</div>
<div class="col-md-3">
    <select class="form-select" id="mahasiswaSortOrder">
        <option value="asc">A-Z / Lama-Baru</option>
        <option value="desc">Z-A / Baru-Lama</option>
    </select>
</div>
```

#### Pagination Controls
```html
<div class="d-flex justify-content-between align-items-center mt-3">
    <div>
        <span class="text-muted" id="mahasiswaInfo">Menampilkan 0 - 0 dari 0 data</span>
    </div>
    <nav aria-label="Mahasiswa pagination">
        <ul class="pagination pagination-sm mb-0" id="mahasiswaPagination">
            <!-- Pagination items dynamically generated -->
        </ul>
    </nav>
</div>
```

### 2. JavaScript Implementation

#### Global State Variables
```javascript
// Pagination state
let mahasiswaCurrentPage = 1;
let alumniCurrentPage = 1;
let pekerjaanCurrentPage = 1;
const itemsPerPage = 10;
```

#### Pagination Helper Functions
```javascript
function generatePagination(currentPage, totalPages, paginationId) {
    const pagination = document.getElementById(paginationId);
    pagination.innerHTML = '';
    
    // Previous button
    const prevButton = document.createElement('li');
    prevButton.className = `page-item ${currentPage === 1 ? 'disabled' : ''}`;
    prevButton.innerHTML = `<a class="page-link" href="#" onclick="changePage('${paginationId}', ${currentPage - 1})">Previous</a>`;
    pagination.appendChild(prevButton);
    
    // Page numbers
    for (let i = Math.max(1, currentPage - 2); i <= Math.min(totalPages, currentPage + 2); i++) {
        const pageItem = document.createElement('li');
        pageItem.className = `page-item ${i === currentPage ? 'active' : ''}`;
        pageItem.innerHTML = `<a class="page-link" href="#" onclick="changePage('${paginationId}', ${i})">${i}</a>`;
        pagination.appendChild(pageItem);
    }
    
    // Next button
    const nextButton = document.createElement('li');
    nextButton.className = `page-item ${currentPage === totalPages ? 'disabled' : ''}`;
    nextButton.innerHTML = `<a class="page-link" href="#" onclick="changePage('${paginationId}', ${currentPage + 1})">Next</a>`;
    pagination.appendChild(nextButton);
}

function updatePaginationInfo(start, end, total, infoId) {
    document.getElementById(infoId).textContent = `Menampilkan ${start} - ${end} dari ${total} data`;
}
```

#### Data Loading Functions
```javascript
async function loadMahasiswa(page = 1, search = '', sortBy = '', sortOrder = 'asc') {
    try {
        // Build URL with parameters
        const params = new URLSearchParams({
            page: page,
            limit: itemsPerPage,
            search: search,
            sort_by: sortBy,
            sort_order: sortOrder
        });
        
        const response = await authorizedFetch(`/api/mahasiswa?${params}`);
        const result = await response.json();
        
        if (response.ok) {
            // Update table
            displayMahasiswaData(result.data);
            
            // Update pagination
            generatePagination(result.page, result.total_pages, 'mahasiswaPagination');
            
            // Update info
            const start = ((result.page - 1) * result.limit) + 1;
            const end = Math.min(result.page * result.limit, result.total_items);
            updatePaginationInfo(start, end, result.total_items, 'mahasiswaInfo');
            
            mahasiswaCurrentPage = result.page;
        }
    } catch (error) {
        console.error('Error loading mahasiswa:', error);
    }
}
```

#### Search Functions
```javascript
function searchMahasiswa() {
    const nim = document.getElementById('mahasiswaSearchNim').value;
    const nama = document.getElementById('mahasiswaSearchNama').value;
    const jurusan = document.getElementById('mahasiswaSearchJurusan').value;
    const tahunMasuk = document.getElementById('mahasiswaSearchTahunMasuk').value;
    const email = document.getElementById('mahasiswaSearchEmail').value;
    
    // Combine all search terms
    const searchTerms = [nim, nama, jurusan, tahunMasuk, email].filter(term => term.trim() !== '');
    const combinedSearch = searchTerms.join(' ');
    
    const sortBy = document.getElementById('mahasiswaSortBy').value;
    const sortOrder = document.getElementById('mahasiswaSortOrder').value;
    
    loadMahasiswa(1, combinedSearch, sortBy, sortOrder);
}
```

#### Real-time Search Event Listeners
```javascript
// Add event listeners for real-time search
document.addEventListener('DOMContentLoaded', function() {
    // Mahasiswa search inputs
    ['mahasiswaSearchNim', 'mahasiswaSearchNama', 'mahasiswaSearchJurusan', 'mahasiswaSearchTahunMasuk', 'mahasiswaSearchEmail'].forEach(id => {
        const input = document.getElementById(id);
        if (input) {
            input.addEventListener('input', debounce(searchMahasiswa, 300));
        }
    });
    
    // Sort dropdowns
    ['mahasiswaSortBy', 'mahasiswaSortOrder'].forEach(id => {
        const select = document.getElementById(id);
        if (select) {
            select.addEventListener('change', searchMahasiswa);
        }
    });
});

// Debounce function to limit API calls
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}
```

---

## API Endpoints

### Paginated Endpoints

#### GET /api/mahasiswa
**Parameters:**
- `page` (int, optional): Page number (default: 1)
- `limit` (int, optional): Items per page (default: 10)
- `search` (string, optional): Search term
- `sort_by` (string, optional): Field to sort by
- `sort_order` (string, optional): "asc" or "desc" (default: "asc")

**Response:**
```json
{
    "data": [
        {
            "id": 1,
            "nim": "12345678",
            "nama": "John Doe",
            "jurusan": "Teknik Informatika",
            "tahun_masuk": 2020,
            "email": "john@example.com"
        }
    ],
    "total_items": 100,
    "total_pages": 10,
    "page": 1,
    "limit": 10,
    "has_next": true,
    "has_prev": false
}
```

#### GET /api/alumni
Similar structure to mahasiswa endpoint with alumni-specific fields.

#### GET /api/pekerjaan-alumni
Similar structure with job-specific fields including alumni relationship.

### Legacy Endpoints (Backward Compatibility)

#### GET /api/mahasiswa/all
Returns all mahasiswa records without pagination.

#### GET /api/alumni/all
Returns all alumni records without pagination.

#### GET /api/pekerjaan-alumni/all
Returns all job records without pagination.

---

## Usage Examples

### 1. Basic Pagination
```javascript
// Load first page
loadMahasiswa(1);

// Load specific page
loadMahasiswa(3);
```

### 2. Search with Pagination
```javascript
// Search for students with "John" in any field
loadMahasiswa(1, 'John');

// Search with sorting
loadMahasiswa(1, 'Teknik', 'nama', 'asc');
```

### 3. Sorting Only
```javascript
// Sort by name ascending
loadMahasiswa(1, '', 'nama', 'asc');

// Sort by year descending
loadMahasiswa(1, '', 'tahun_masuk', 'desc');
```

### 4. Combined Usage
```javascript
// Search "John", sort by name, page 2
loadMahasiswa(2, 'John', 'nama', 'asc');
```

---

## Testing Guide

### 1. Manual Testing Steps

#### Test Pagination
1. Load the application
2. Navigate to Mahasiswa section
3. Verify data loads with pagination controls
4. Click "Next" and "Previous" buttons
5. Click specific page numbers
6. Verify page info updates correctly

#### Test Search
1. Enter search term in any search field
2. Verify results filter correctly
3. Test multiple search fields simultaneously
4. Test with no results
5. Clear search and verify all data returns

#### Test Sorting
1. Select different sort fields
2. Toggle between ascending/descending
3. Verify data order changes correctly
4. Test sorting with search active

#### Test Real-time Features
1. Type in search fields slowly
2. Verify results update automatically (with debounce)
3. Change sort options
4. Verify immediate updates

### 2. API Testing with Postman/curl

#### Test Pagination Endpoint
```bash
curl -X GET "http://localhost:3000/api/mahasiswa?page=1&limit=5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Test Search
```bash
curl -X GET "http://localhost:3000/api/mahasiswa?search=John&page=1&limit=10" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Test Sorting
```bash
curl -X GET "http://localhost:3000/api/mahasiswa?sort_by=nama&sort_order=desc" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 3. Edge Cases to Test

- Empty search results
- Invalid page numbers (negative, zero, beyond total)
- Invalid sort fields
- Invalid sort orders
- Large search terms
- Special characters in search
- Very large page numbers
- Simultaneous requests

---

## Troubleshooting

### Common Issues

#### 1. Pagination Not Working
**Symptoms:** Pagination controls not appearing or not functional
**Solutions:**
- Check browser console for JavaScript errors
- Verify API endpoints are returning correct response format
- Ensure JWT token is valid and included in requests

#### 2. Search Not Filtering
**Symptoms:** Search input doesn't filter results
**Solutions:**
- Check if search term is being passed to API
- Verify ILIKE queries in repository implementation
- Check for SQL syntax errors in logs

#### 3. Sorting Not Working
**Symptoms:** Data order doesn't change when selecting sort options
**Solutions:**
- Verify sort field names match database columns
- Check sort order validation in backend
- Ensure ORDER BY clause is properly formatted

#### 4. Real-time Search Too Aggressive
**Symptoms:** Too many API calls while typing
**Solutions:**
- Increase debounce delay (currently 300ms)
- Check if debounce function is properly implemented
- Monitor network tab for excessive requests

### Debugging Tips

#### Backend Debugging
```go
// Add logging to repository methods
log.Printf("Pagination query: page=%d, limit=%d, search=%s, sortBy=%s, sortOrder=%s", 
    pagination.Page, pagination.Limit, pagination.Search, pagination.SortBy, pagination.SortOrder)
```

#### Frontend Debugging
```javascript
// Add console logging to track state
console.log('Loading page:', page, 'Search:', search, 'Sort:', sortBy, sortOrder);
console.log('API Response:', result);
```

#### Database Query Debugging
Enable GORM logging to see generated SQL queries:
```go
db.Logger = logger.Default.LogMode(logger.Info)
```

### Performance Considerations

1. **Database Indexing**: Ensure indexes on frequently searched/sorted columns
2. **Limit Size**: Consider maximum limit to prevent large result sets
3. **Search Optimization**: Use full-text search for complex search requirements
4. **Caching**: Implement caching for frequently accessed data
5. **Connection Pooling**: Optimize database connection settings

---

## Conclusion

This pagination, sorting, and search implementation provides a robust foundation for data management in the CRUD Go Fiber PostgreSQL application. The system is designed to be:

- **Scalable**: Handles large datasets efficiently
- **User-friendly**: Intuitive interface with real-time feedback
- **Maintainable**: Clean architecture with separation of concerns
- **Extensible**: Easy to add new search fields or sorting options
- **Compatible**: Maintains backward compatibility with existing endpoints

For additional features or modifications, refer to the respective code sections and follow the established patterns.