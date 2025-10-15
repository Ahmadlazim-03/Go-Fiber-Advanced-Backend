package repositories

import "modul4crud/models"

// UserRepository interface untuk operasi user
type UserRepository interface {
	GetAll() ([]models.User, error)
	GetWithPagination(pagination *models.PaginationRequest) ([]models.User, int64, error)
	GetByID(id int) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id int) error
	Count() (int64, error)
	// AuthenticateWithPassword verifies credentials (PocketBase specific)
	// For PostgreSQL/MongoDB, this returns error since they use bcrypt
	AuthenticateWithPassword(email, password string) (*models.User, error)
}

// MahasiswaRepository interface untuk operasi mahasiswa
type MahasiswaRepository interface {
	GetAll() ([]models.Mahasiswa, error)
	GetWithPagination(pagination *models.PaginationRequest) ([]models.Mahasiswa, int64, error)
	GetByID(id uint) (*models.Mahasiswa, error)
	Create(mahasiswa *models.Mahasiswa) error
	Update(mahasiswa *models.Mahasiswa) error
	Delete(id uint) error
	Count() (int64, error)
}

// AlumniRepository interface untuk operasi alumni
type AlumniRepository interface {
	GetAll() ([]models.Alumni, error)
	GetWithPagination(pagination *models.PaginationRequest) ([]models.Alumni, int64, error)
	GetByID(id uint) (*models.Alumni, error)
	GetByUserID(userID int) (*models.Alumni, error)
	Create(alumni *models.Alumni) error
	Update(alumni *models.Alumni) error
	Delete(id uint) error
	Count() (int64, error)
}

// PekerjaanAlumniRepository interface untuk operasi pekerjaan alumni
type PekerjaanAlumniRepository interface {
	GetAll() ([]models.PekerjaanAlumni, error)
	GetWithPagination(pagination *models.PaginationRequest) ([]models.PekerjaanAlumni, int64, error)
	GetByID(id uint) (*models.PekerjaanAlumni, error)
	GetByAlumniID(alumniID uint) ([]models.PekerjaanAlumni, error)
	GetByUserID(userID int) ([]models.PekerjaanAlumni, error)
	Create(pekerjaan *models.PekerjaanAlumni) error
	Update(pekerjaan *models.PekerjaanAlumni) error
	Delete(id uint) error
	SoftDelete(id uint) error
	SoftDeleteByAlumniID(alumniID uint) error
	Restore(id uint) error
	GetDeleted() ([]models.PekerjaanAlumni, error)
	GetDeletedByUserID(userID int) ([]models.PekerjaanAlumni, error)
	Count() (int64, error)
	GetAlumniCountByCompany(namaPerusahaan string) (int64, error)
}
