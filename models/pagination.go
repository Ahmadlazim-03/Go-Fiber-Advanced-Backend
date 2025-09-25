package models

// PaginationRequest untuk request pagination, sorting, dan search
type PaginationRequest struct {
	Page      int    `query:"page" json:"page"`           // Halaman yang diminta (default 1)
	Limit     int    `query:"limit" json:"limit"`         // Jumlah data per halaman (default 10)
	Search    string `query:"search" json:"search"`       // Keyword untuk search
	SortBy    string `query:"sort_by" json:"sort_by"`     // Field untuk sorting (default "id")
	SortOrder string `query:"sort_order" json:"sort_order"` // ASC atau DESC (default "ASC")
}

// PaginationResponse untuk response pagination
type PaginationResponse struct {
	Data         interface{} `json:"data"`          // Data hasil query
	CurrentPage  int         `json:"current_page"`  // Halaman saat ini
	PerPage      int         `json:"per_page"`      // Jumlah data per halaman
	TotalData    int64       `json:"total_data"`    // Total data keseluruhan
	TotalPages   int         `json:"total_pages"`   // Total halaman
	HasNext      bool        `json:"has_next"`      // Apakah ada halaman selanjutnya
	HasPrevious  bool        `json:"has_previous"`  // Apakah ada halaman sebelumnya
	NextPage     *int        `json:"next_page"`     // Nomor halaman selanjutnya
	PreviousPage *int        `json:"previous_page"` // Nomor halaman sebelumnya
}

// SetDefaults mengatur default values untuk PaginationRequest
func (p *PaginationRequest) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.SortBy == "" {
		p.SortBy = "id"
	}
	if p.SortOrder == "" {
		p.SortOrder = "ASC"
	}
}

// GetOffset menghitung offset untuk query database
func (p *PaginationRequest) GetOffset() int {
	return (p.Page - 1) * p.Limit
}

// ValidateSortOrder memvalidasi sort order
func (p *PaginationRequest) ValidateSortOrder() {
	// Convert to uppercase for consistency
	if p.SortOrder == "asc" {
		p.SortOrder = "ASC"
	} else if p.SortOrder == "desc" {
		p.SortOrder = "DESC"
	} else if p.SortOrder != "ASC" && p.SortOrder != "DESC" {
		p.SortOrder = "ASC"
	}
}

// NewPaginationResponse membuat response pagination
func NewPaginationResponse(data interface{}, pagination *PaginationRequest, totalData int64) *PaginationResponse {
	totalPages := int((totalData + int64(pagination.Limit) - 1) / int64(pagination.Limit))
	
	response := &PaginationResponse{
		Data:        data,
		CurrentPage: pagination.Page,
		PerPage:     pagination.Limit,
		TotalData:   totalData,
		TotalPages:  totalPages,
		HasNext:     pagination.Page < totalPages,
		HasPrevious: pagination.Page > 1,
	}
	
	// Set next page
	if response.HasNext {
		nextPage := pagination.Page + 1
		response.NextPage = &nextPage
	}
	
	// Set previous page
	if response.HasPrevious {
		previousPage := pagination.Page - 1
		response.PreviousPage = &previousPage
	}
	
	return response
}