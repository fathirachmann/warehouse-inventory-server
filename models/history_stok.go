package models

import "time"

type HistoryStok struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	BarangID       uint      `gorm:"not null" json:"barang_id"`
	UserID         uint      `gorm:"not null" json:"user_id"`
	JenisTransaksi string    `gorm:"not null" json:"jenis_transaksi"` // "masuk" or "keluar"
	Jumlah         int       `gorm:"not null" json:"jumlah"`
	StokSebelumnya int       `gorm:"not null" json:"stok_sebelumnya"`
	StokSesudah    int       `gorm:"not null" json:"stok_sesudah"`
	Keterangan     string    `json:"keterangan"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Associations
	MasterBarang MasterBarang `gorm:"foreignKey:BarangID;references:ID" json:"master_barang"` // HistoryStok many to one MasterBarang
	Users        User         `gorm:"foreignKey:UserID;references:ID" json:"user"`            // HistoryStok many to one User
}

func (HistoryStok) TableName() string {
	return "history_stok"
}
