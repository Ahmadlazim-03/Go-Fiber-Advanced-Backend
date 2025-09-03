package models

import "time"

// Model Alumni

type Alumni struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	NIM        string            `gorm:"type:varchar(20);unique;not null" json:"nim"`
	Nama       string            `gorm:"type:varchar(100);not null" json:"nama"`
	Jurusan    string            `gorm:"type:varchar(50);not null" json:"jurusan"`
	Angkatan   int               `gorm:"not null" json:"angkatan"`
	TahunLulus int               `gorm:"not null" json:"tahun_lulus"`
	Email      string            `gorm:"type:varchar(100);unique;not null" json:"email"`
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
}
