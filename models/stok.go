package models

import (
	"time"
)

// Model struct for mstok table
type Mstok struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BarangID  uint      `gorm:"not null" json:"barang_id"`
	StokAkhir int       `gorm:"default:0" json:"stok_akhir"`
	UpdatedAt time.Time `json:"updated_at"`

	// Associations
	MasterBarang MasterBarang `gorm:"foreignKey:BarangID;references:ID" json:"barang"`
}

func (Mstok) TableName() string {
	return "mstok"
}

// Response struct for mstok API
type MstokResponse struct {
	ID        uint               `json:"id"`
	BarangID  uint               `json:"barang_id"`
	StokAkhir int                `json:"stok_akhir"`
	UpdatedAt time.Time          `json:"updated_at"`
	Barang    BarangStokResponse `json:"barang"`
}

type BarangStokResponse struct {
	KodeBarang string  `json:"kode_barang"`
	NamaBarang string  `json:"nama_barang"`
	Satuan     string  `json:"satuan"`
	HargaJual  float64 `json:"harga_jual"`
}
