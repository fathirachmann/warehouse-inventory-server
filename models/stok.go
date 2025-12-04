package models

import (
	"time"
)

type Mstok struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BarangID  uint      `gorm:"not null" json:"barang_id"`
	StokAkhir int       `gorm:"default:0" json:"stok_akhir"`
	UpdatedAt time.Time `json:"updated_at"`

	// Associations
	MasterBarang MasterBarang `gorm:"foreignKey:BarangID;references:ID" json:"master_barang"`
}

func (Mstok) TableName() string {
	return "mstok"
}
