package postgre

import (
	"fmt"
	"modul4crud/models"
	repo "modul4crud/repositories/interface"
	"time"

	"gorm.io/gorm"
)

type pekerjaanAlumniRepository struct {
	db *gorm.DB
}

func NewPekerjaanAlumniRepository(db *gorm.DB) repo.PekerjaanAlumniRepository {
	return &pekerjaanAlumniRepository{db: db}
}

func (r *pekerjaanAlumniRepository) GetAll() ([]models.PekerjaanAlumni, error) {
	var pekerjaans []models.PekerjaanAlumni

	query := `
		SELECT 
			pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, 
			pa.bidang_industri, pa.lokasi_kerja, pa.gaji_range, 
			pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja, 
			pa.status_pekerjaan, pa.deskripsi_pekerjaan, 
			pa.created_at, pa.updated_at, pa.deleted_at,
			a.id as "Alumni__id", a.user_id as "Alumni__user_id", 
			a.nim as "Alumni__nim", a.nama as "Alumni__nama", 
			a.jurusan as "Alumni__jurusan", a.angkatan as "Alumni__angkatan", 
			a.tahun_lulus as "Alumni__tahun_lulus", a.no_telepon as "Alumni__no_telepon", 
			a.alamat as "Alumni__alamat", a.created_at as "Alumni__created_at", 
			a.updated_at as "Alumni__updated_at",
			u.id as "Alumni__User__id", u.username as "Alumni__User__username", 
			u.email as "Alumni__User__email", u.role as "Alumni__User__role", 
			u.is_active as "Alumni__User__is_active", u.created_at as "Alumni__User__created_at", 
			u.updated_at as "Alumni__User__updated_at"
		FROM pekerjaan_alumnis pa
		LEFT JOIN alumnis a ON pa.alumni_id = a.id
		LEFT JOIN users u ON a.user_id = u.id
		WHERE pa.deleted_at IS NULL
		ORDER BY pa.id DESC
	`

	err := r.db.Raw(query).Scan(&pekerjaans).Error
	return pekerjaans, err
}

func (r *pekerjaanAlumniRepository) GetWithPagination(pagination *models.PaginationRequest) ([]models.PekerjaanAlumni, int64, error) {
	var pekerjaans []models.PekerjaanAlumni
	var total int64

	// Set default values
	pagination.SetDefaults()
	pagination.ValidateSortOrder()

	// Count query
	countQuery := `
		SELECT COUNT(*) 
		FROM pekerjaan_alumnis pa
		LEFT JOIN alumnis a ON pa.alumni_id = a.id
		WHERE pa.deleted_at IS NULL
	`

	// Search filter
	searchCondition := ""
	searchArgs := []interface{}{}
	if pagination.Search != "" {
		searchPattern := "%" + pagination.Search + "%"
		searchCondition = ` AND (
			pa.posisi_jabatan ILIKE ? OR 
			pa.nama_perusahaan ILIKE ? OR 
			pa.bidang_industri ILIKE ? OR 
			pa.lokasi_kerja ILIKE ? OR 
			a.nama ILIKE ?
		)`
		searchArgs = []interface{}{searchPattern, searchPattern, searchPattern, searchPattern, searchPattern}
	}

	// Execute count query
	err := r.db.Raw(countQuery+searchCondition, searchArgs...).Scan(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Data query
	dataQuery := `
		SELECT 
			pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, 
			pa.bidang_industri, pa.lokasi_kerja, pa.gaji_range, 
			pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja, 
			pa.status_pekerjaan, pa.deskripsi_pekerjaan, 
			pa.created_at, pa.updated_at, pa.deleted_at,
			a.id as "Alumni__id", a.user_id as "Alumni__user_id", 
			a.nim as "Alumni__nim", a.nama as "Alumni__nama", 
			a.jurusan as "Alumni__jurusan", a.angkatan as "Alumni__angkatan", 
			a.tahun_lulus as "Alumni__tahun_lulus", a.no_telepon as "Alumni__no_telepon", 
			a.alamat as "Alumni__alamat", a.created_at as "Alumni__created_at", 
			a.updated_at as "Alumni__updated_at",
			u.id as "Alumni__User__id", u.username as "Alumni__User__username", 
			u.email as "Alumni__User__email", u.role as "Alumni__User__role", 
			u.is_active as "Alumni__User__is_active", u.created_at as "Alumni__User__created_at", 
			u.updated_at as "Alumni__User__updated_at"
		FROM pekerjaan_alumnis pa
		LEFT JOIN alumnis a ON pa.alumni_id = a.id
		LEFT JOIN users u ON a.user_id = u.id
		WHERE pa.deleted_at IS NULL
	`

	// Add search condition to data query
	dataQuery += searchCondition

	// Add sorting and pagination
	dataQuery += fmt.Sprintf(" ORDER BY pa.%s %s LIMIT ? OFFSET ?", pagination.SortBy, pagination.SortOrder)

	// Prepare arguments for data query
	dataArgs := append(searchArgs, pagination.Limit, pagination.GetOffset())

	err = r.db.Raw(dataQuery, dataArgs...).Scan(&pekerjaans).Error
	return pekerjaans, total, err
}

func (r *pekerjaanAlumniRepository) GetByID(id uint) (*models.PekerjaanAlumni, error) {
	var pekerjaan models.PekerjaanAlumni

	query := `
		SELECT 
			pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, 
			pa.bidang_industri, pa.lokasi_kerja, pa.gaji_range, 
			pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja, 
			pa.status_pekerjaan, pa.deskripsi_pekerjaan, 
			pa.created_at, pa.updated_at, pa.deleted_at,
			a.id as "Alumni__id", a.user_id as "Alumni__user_id", 
			a.nim as "Alumni__nim", a.nama as "Alumni__nama", 
			a.jurusan as "Alumni__jurusan", a.angkatan as "Alumni__angkatan", 
			a.tahun_lulus as "Alumni__tahun_lulus", a.no_telepon as "Alumni__no_telepon", 
			a.alamat as "Alumni__alamat", a.created_at as "Alumni__created_at", 
			a.updated_at as "Alumni__updated_at",
			u.id as "Alumni__User__id", u.username as "Alumni__User__username", 
			u.email as "Alumni__User__email", u.role as "Alumni__User__role", 
			u.is_active as "Alumni__User__is_active", u.created_at as "Alumni__User__created_at", 
			u.updated_at as "Alumni__User__updated_at"
		FROM pekerjaan_alumnis pa
		LEFT JOIN alumnis a ON pa.alumni_id = a.id
		LEFT JOIN users u ON a.user_id = u.id
		WHERE pa.id = ? AND pa.deleted_at IS NULL
	`

	err := r.db.Raw(query, id).Scan(&pekerjaan).Error
	if err != nil {
		return nil, err
	}
	return &pekerjaan, nil
}

func (r *pekerjaanAlumniRepository) GetByAlumniID(alumniID uint) ([]models.PekerjaanAlumni, error) {
	var pekerjaans []models.PekerjaanAlumni

	query := `
		SELECT 
			pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, 
			pa.bidang_industri, pa.lokasi_kerja, pa.gaji_range, 
			pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja, 
			pa.status_pekerjaan, pa.deskripsi_pekerjaan, 
			pa.created_at, pa.updated_at, pa.deleted_at,
			a.id as "Alumni__id", a.user_id as "Alumni__user_id", 
			a.nim as "Alumni__nim", a.nama as "Alumni__nama", 
			a.jurusan as "Alumni__jurusan", a.angkatan as "Alumni__angkatan", 
			a.tahun_lulus as "Alumni__tahun_lulus", a.no_telepon as "Alumni__no_telepon", 
			a.alamat as "Alumni__alamat", a.created_at as "Alumni__created_at", 
			a.updated_at as "Alumni__updated_at",
			u.id as "Alumni__User__id", u.username as "Alumni__User__username", 
			u.email as "Alumni__User__email", u.role as "Alumni__User__role", 
			u.is_active as "Alumni__User__is_active", u.created_at as "Alumni__User__created_at", 
			u.updated_at as "Alumni__User__updated_at"
		FROM pekerjaan_alumnis pa
		LEFT JOIN alumnis a ON pa.alumni_id = a.id
		LEFT JOIN users u ON a.user_id = u.id
		WHERE pa.alumni_id = ? AND pa.deleted_at IS NULL
		ORDER BY pa.id DESC
	`

	err := r.db.Raw(query, alumniID).Scan(&pekerjaans).Error
	return pekerjaans, err
}

func (r *pekerjaanAlumniRepository) GetByUserID(userID int) ([]models.PekerjaanAlumni, error) {
	var pekerjaans []models.PekerjaanAlumni

	query := `
		SELECT 
			pa.id, pa.alumni_id, pa.nama_perusahaan, pa.posisi_jabatan, 
			pa.bidang_industri, pa.lokasi_kerja, pa.gaji_range, 
			pa.tanggal_mulai_kerja, pa.tanggal_selesai_kerja, 
			pa.status_pekerjaan, pa.deskripsi_pekerjaan, 
			pa.created_at, pa.updated_at, pa.deleted_at,
			a.id as "Alumni__id", a.user_id as "Alumni__user_id", 
			a.nim as "Alumni__nim", a.nama as "Alumni__nama", 
			a.jurusan as "Alumni__jurusan", a.angkatan as "Alumni__angkatan", 
			a.tahun_lulus as "Alumni__tahun_lulus", a.no_telepon as "Alumni__no_telepon", 
			a.alamat as "Alumni__alamat", a.created_at as "Alumni__created_at", 
			a.updated_at as "Alumni__updated_at",
			u.id as "Alumni__User__id", u.username as "Alumni__User__username", 
			u.email as "Alumni__User__email", u.role as "Alumni__User__role", 
			u.is_active as "Alumni__User__is_active", u.created_at as "Alumni__User__created_at", 
			u.updated_at as "Alumni__User__updated_at"
		FROM pekerjaan_alumnis pa
		LEFT JOIN alumnis a ON pa.alumni_id = a.id
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.user_id = ? AND pa.deleted_at IS NULL
		ORDER BY pa.id DESC
	`

	err := r.db.Raw(query, userID).Scan(&pekerjaans).Error
	return pekerjaans, err
}

func (r *pekerjaanAlumniRepository) Create(pekerjaan *models.PekerjaanAlumni) error {
	query := `
		INSERT INTO pekerjaan_alumnis 
		(alumni_id, nama_perusahaan, posisi_jabatan, bidang_industri, lokasi_kerja, 
		 gaji_range, tanggal_mulai_kerja, tanggal_selesai_kerja, status_pekerjaan, 
		 deskripsi_pekerjaan, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	return r.db.Raw(query,
		pekerjaan.AlumniID,
		pekerjaan.NamaPerusahaan,
		pekerjaan.PosisiJabatan,
		pekerjaan.BidangIndustri,
		pekerjaan.LokasiKerja,
		pekerjaan.GajiRange,
		pekerjaan.TanggalMulaiKerja,
		pekerjaan.TanggalSelesaiKerja,
		pekerjaan.StatusPekerjaan,
		pekerjaan.DeskripsiPekerjaan,
	).Scan(pekerjaan).Error
}

func (r *pekerjaanAlumniRepository) Update(pekerjaan *models.PekerjaanAlumni) error {
	query := `
		UPDATE pekerjaan_alumnis 
		SET nama_perusahaan = ?, posisi_jabatan = ?, bidang_industri = ?, 
		    lokasi_kerja = ?, gaji_range = ?, tanggal_mulai_kerja = ?, 
		    tanggal_selesai_kerja = ?, status_pekerjaan = ?, 
		    deskripsi_pekerjaan = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
		RETURNING updated_at
	`

	return r.db.Raw(query,
		pekerjaan.NamaPerusahaan,
		pekerjaan.PosisiJabatan,
		pekerjaan.BidangIndustri,
		pekerjaan.LokasiKerja,
		pekerjaan.GajiRange,
		pekerjaan.TanggalMulaiKerja,
		pekerjaan.TanggalSelesaiKerja,
		pekerjaan.StatusPekerjaan,
		pekerjaan.DeskripsiPekerjaan,
		pekerjaan.ID,
	).Scan(pekerjaan).Error
}

func (r *pekerjaanAlumniRepository) Delete(id uint) error {
	var deletedAt *time.Time
	checkQuery := `SELECT deleted_at FROM pekerjaan_alumnis WHERE id = ?`
	err := r.db.Raw(checkQuery, id).Scan(&deletedAt).Error
	if err != nil {
		return fmt.Errorf("data pekerjaan alumni tidak ditemukan")
	}

	if deletedAt == nil {
		return fmt.Errorf("tidak bisa hard delete: data belum di-soft delete terlebih dahulu")
	}

	query := `DELETE FROM pekerjaan_alumnis WHERE id = ? AND deleted_at IS NOT NULL`
	result := r.db.Exec(query, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("tidak ada data yang dihapus - pastikan data sudah di-soft delete")
	}
	return result.Error
}

func (r *pekerjaanAlumniRepository) Count() (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM pekerjaan_alumnis WHERE deleted_at IS NULL`
	err := r.db.Raw(query).Scan(&count).Error
	return count, err
}

func (r *pekerjaanAlumniRepository) GetAlumniCountByCompany(namaPerusahaan string) (int64, error) {
	var count int64
	query := `SELECT COUNT(DISTINCT alumni_id) FROM pekerjaan_alumnis WHERE nama_perusahaan = ? AND deleted_at IS NULL`
	err := r.db.Raw(query, namaPerusahaan).Scan(&count).Error
	return count, err
}

// Soft Delete methods
func (r *pekerjaanAlumniRepository) SoftDelete(id uint) error {
	query := `UPDATE pekerjaan_alumnis SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	result := r.db.Exec(query, id)
	return result.Error
}

func (r *pekerjaanAlumniRepository) SoftDeleteByAlumniID(alumniID uint) error {
	query := `UPDATE pekerjaan_alumnis SET deleted_at = NOW() WHERE alumni_id = ? AND deleted_at IS NULL`
	result := r.db.Exec(query, alumniID)
	return result.Error
}

func (r *pekerjaanAlumniRepository) Restore(id uint) error {
	query := `UPDATE pekerjaan_alumnis SET deleted_at = NULL WHERE id = ?`
	result := r.db.Exec(query, id)
	return result.Error
}

func (r *pekerjaanAlumniRepository) GetDeleted() ([]models.PekerjaanAlumni, error) {
	var pekerjaans []models.PekerjaanAlumni

	query := `
		SELECT 
			pekerjaan_alumnis.*, 
			alumnis.*, 
			users.*
		FROM pekerjaan_alumnis
		LEFT JOIN alumnis ON pekerjaan_alumnis.alumni_id = alumnis.id
		LEFT JOIN users ON alumnis.user_id = users.id
		WHERE pekerjaan_alumnis.deleted_at IS NOT NULL
		ORDER BY pekerjaan_alumnis.deleted_at DESC;
	`

	err := r.db.Raw(query).Scan(&pekerjaans).Error
	return pekerjaans, err
}

func (r *pekerjaanAlumniRepository) GetDeletedByUserID(userID int) ([]models.PekerjaanAlumni, error) {
	var pekerjaans []models.PekerjaanAlumni

	query := `
		SELECT 
			pekerjaan_alumnis.*, 
			alumnis.*, 
			users.*
		FROM pekerjaan_alumnis
		LEFT JOIN alumnis ON pekerjaan_alumnis.alumni_id = alumnis.id
		LEFT JOIN users ON alumnis.user_id = users.id
		WHERE pekerjaan_alumnis.deleted_at IS NOT NULL 
		AND alumnis.user_id = ?
		ORDER BY pekerjaan_alumnis.deleted_at DESC;
	`

	err := r.db.Raw(query, userID).Scan(&pekerjaans).Error
	return pekerjaans, err
}
