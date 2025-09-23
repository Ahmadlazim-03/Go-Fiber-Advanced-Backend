package models

import "time"

// Model Alumni

type Alumni struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	UserID     int               `gorm:"not null;index" json:"user_id"`
	User       User              `gorm:"foreignKey:UserID" json:"user"`
	NIM        string            `gorm:"type:varchar(20);unique;not null" json:"nim"`
	Nama       string            `gorm:"type:varchar(100);not null" json:"nama"`
	Jurusan    string            `gorm:"type:varchar(50);not null" json:"jurusan"`
	Angkatan   int               `gorm:"not null" json:"angkatan"`
	TahunLulus int               `gorm:"not null" json:"tahun_lulus"`
	NoTelepon  string            `gorm:"type:varchar(15)" json:"no_telepon"`
	Alamat     string            `gorm:"type:text" json:"alamat"`
	CreatedAt  time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	Pekerjaan  []PekerjaanAlumni `gorm:"foreignKey:AlumniID" json:"pekerjaan_alumni"`
}

// Model Pekerjaan Alumni

type PekerjaanAlumni struct {
	ID                  uint       `gorm:"primaryKey" json:"id"`
	AlumniID            uint       `gorm:"not null" json:"alumni_id"`
	Alumni              Alumni     `gorm:"foreignKey:AlumniID" json:"alumni"`
	NamaPerusahaan      string     `gorm:"type:varchar(100);not null" json:"nama_perusahaan"`
	PosisiJabatan       string     `gorm:"type:varchar(100);not null" json:"posisi_jabatan"`
	BidangIndustri      string     `gorm:"type:varchar(50);not null" json:"bidang_industri"`
	LokasiKerja         string     `gorm:"type:varchar(100);not null" json:"lokasi_kerja"`
	GajiRange           string     `gorm:"type:varchar(50)" json:"gaji_range"`
	TanggalMulaiKerja   time.Time  `gorm:"type:date;not null" json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja *time.Time `gorm:"type:date" json:"tanggal_selesai_kerja"`
	StatusPekerjaan     string     `gorm:"type:varchar(20);default:'aktif'" json:"status_pekerjaan"`
	DeskripsiPekerjaan  string     `gorm:"type:text" json:"deskripsi_pekerjaan"`
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt           *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

// Request struct untuk Alumni
type CreateAlumniRequest struct {
	UserID     int    `json:"user_id"`
	NIM        string `json:"nim"`
	Nama       string `json:"nama"`
	Jurusan    string `json:"jurusan"`
	Angkatan   int    `json:"angkatan"`
	TahunLulus int    `json:"tahun_lulus"`
	NoTelepon  string `json:"no_telepon"`
	Alamat     string `json:"alamat"`
}

type UpdateAlumniRequest struct {
	Nama       string `json:"nama"`
	Jurusan    string `json:"jurusan"`
	Angkatan   int    `json:"angkatan"`
	TahunLulus int    `json:"tahun_lulus"`
	NoTelepon  string `json:"no_telepon"`
	Alamat     string `json:"alamat"`
}

// Request struct untuk PekerjaanAlumni
type CreatePekerjaanAlumniRequest struct {
	AlumniID               uint      `json:"alumni_id"`
	NamaPerusahaan         string    `json:"nama_perusahaan"`
	PosisiJabatan          string    `json:"posisi_jabatan"`
	BidangIndustri         string    `json:"bidang_industri"`
	LokasiKerja            string    `json:"lokasi_kerja"`
	GajiRange              string    `json:"gaji_range"`
	TanggalMulaiKerja      time.Time `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja    *time.Time `json:"tanggal_selesai_kerja"`
	StatusPekerjaan        string    `json:"status_pekerjaan"`
	DeskripsiPekerjaan     string    `json:"deskripsi_pekerjaan"`
}

type UpdatePekerjaanAlumniRequest struct {
	NamaPerusahaan         string     `json:"nama_perusahaan"`
	PosisiJabatan          string     `json:"posisi_jabatan"`
	BidangIndustri         string     `json:"bidang_industri"`
	LokasiKerja            string     `json:"lokasi_kerja"`
	GajiRange              string     `json:"gaji_range"`
	TanggalMulaiKerja      time.Time  `json:"tanggal_mulai_kerja"`
	TanggalSelesaiKerja    *time.Time `json:"tanggal_selesai_kerja"`
	StatusPekerjaan        string     `json:"status_pekerjaan"`
	DeskripsiPekerjaan     string     `json:"deskripsi_pekerjaan"`
}
