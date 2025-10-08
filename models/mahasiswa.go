package models

import "time"

type CreateMahasiswaRequest struct {
	NIM      string `json:"nim"`
	Nama     string `json:"nama"`
	Jurusan  string `json:"jurusan"`
	Angkatan int    `json:"angkatan"`
	Email    string `json:"email"`
}

type UpdateMahasiswaRequest struct {
	Nama     string `json:"nama"`
	Jurusan  string `json:"jurusan"`
	Angkatan int    `json:"angkatan"`
	Email    string `json:"email"`
}

type Mahasiswa struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NIM       string    `gorm:"type:varchar(20);unique;not null" json:"nim"`
	Nama      string    `gorm:"type:varchar(100);not null" json:"nama"`
	Jurusan   string    `gorm:"type:varchar(50);not null" json:"jurusan"`
	Angkatan  int       `gorm:"not null" json:"angkatan"`
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
