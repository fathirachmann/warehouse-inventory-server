package models

import "time"

type MasterBarang struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	KodeBarang string    `gorm:"size:50;not null" json:"kode_barang"`
	NamaBarang string    `gorm:"size:255;not null" json:"nama_barang"`
	Deskripsi  string    `gorm:"size:512" json:"deskripsi"`
	Satuan     string    `gorm:"size:50;not null" json:"satuan"`
	HargaBeli  float64   `gorm:"not null" json:"harga_beli"`
	HargaJual  float64   `gorm:"not null" json:"harga_jual"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (MasterBarang) TableName() string {
	return "master_barang"
}
