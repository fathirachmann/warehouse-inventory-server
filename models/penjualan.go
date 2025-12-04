package models

import "time"

// JualHeader merepresentasikan tabel 'jual_header' (Info Faktur Penjualan)
type JualHeader struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NoFaktur  string    `gorm:"type:varchar(100);unique;not null" json:"no_faktur"` //
	Customer  string    `gorm:"type:varchar(200);not null" json:"customer"`         // Nama pembeli
	Total     float64   `gorm:"type:decimal(15,2);default:0" json:"total"`          //
	UserID    uint      `gorm:"not null" json:"user_id"`                            //
	Status    string    `gorm:"type:varchar(50);default:'selesai'" json:"status"`   //
	CreatedAt time.Time `json:"created_at"`                                         //

	// Associations
	Details []JualDetail `gorm:"foreignKey:JualHeaderID" json:"details,omitempty"` // JualHeader one to many JualDetail
	User    *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`          // JualHeader many to one User
}

func (JualHeader) TableName() string {
	return "jual_header"
}

// JualDetail merepresentasikan tabel 'jual_detail' (Barang yang dijual)
type JualDetail struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	JualHeaderID uint    `gorm:"not null" json:"jual_header_id"`              // FK ke Header
	BarangID     uint    `gorm:"not null" json:"barang_id"`                   // FK ke Barang
	Qty          int     `gorm:"not null" json:"qty"`                         // Jumlah dijual
	Harga        float64 `gorm:"type:decimal(15,2);not null" json:"harga"`    // Harga saat transaksi
	Subtotal     float64 `gorm:"type:decimal(15,2);not null" json:"subtotal"` // Qty * Harga

	// Associations
	MasterBarang *MasterBarang `gorm:"foreignKey:BarangID" json:"barang,omitempty"` // JualDetail many to one MasterBarang
}

func (JualDetail) TableName() string {
	return "jual_detail"
}

// Request structs for penjualan API
type JualDetailRequest struct {
	BarangID uint    `json:"barang_id"`
	Qty      int     `json:"qty"`
	Harga    float64 `json:"harga"`
}

type JualHeaderRequest struct {
	Customer string              `json:"customer"`
	Details  []JualDetailRequest `json:"details"`
}
