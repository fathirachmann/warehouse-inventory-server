package models

import "time"

type BeliHeader struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NoFaktur  string    `gorm:"type:varchar(100);unique;not null" json:"no_faktur"`
	Supplier  string    `gorm:"type:varchar(200);not null" json:"supplier"`
	Total     float64   `gorm:"type:decimal(15,2);default:0" json:"total"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Status    string    `gorm:"type:varchar(50);default:'selesai'" json:"status"`
	CreatedAt time.Time `json:"created_at"`

	// Associations
	Details []BeliDetail `gorm:"foreignKey:BeliHeaderID" json:"details,omitempty"` // BeliHeader one to many BeliDetail
	User    *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`          // BeliHeader many to one User
}

func (BeliHeader) TableName() string {
	return "beli_header"
}

type BeliDetail struct {
	ID           uint    `gorm:"primaryKey" json:"id"`
	BeliHeaderID uint    `gorm:"not null" json:"beli_header_id"`
	BarangID     uint    `gorm:"not null" json:"barang_id"`
	Qty          int     `gorm:"not null" json:"qty"`
	Harga        float64 `gorm:"type:decimal(15,2);not null" json:"harga"`
	Subtotal     float64 `gorm:"type:decimal(15,2);not null" json:"subtotal"`

	// Associations
	MasterBarang *MasterBarang `gorm:"foreignKey:BarangID" json:"barang,omitempty"` // BeliDetail many to one MasterBarang
}

func (BeliDetail) TableName() string {
	return "beli_detail"
}

// Request structs for pembelian API
type BeliDetailRequest struct {
	BarangID uint    `json:"barang_id"`
	Qty      int     `json:"qty"`
	Harga    float64 `json:"harga"`
}

type BeliHeaderRequest struct {
	Supplier string              `json:"supplier"`
	Details  []BeliDetailRequest `json:"details"`
}

// Response structs for pembelian API
type BeliHeaderResponse struct {
	ID        uint               `json:"id"`
	NoFaktur  string             `json:"no_faktur"`
	Supplier  string             `json:"supplier"`
	Total     float64            `json:"total"`
	UserID    uint               `json:"user_id"`
	Status    string             `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	User      UserSimpleResponse `json:"user"`
}

type BeliDetailResponse struct {
	ID       uint                    `json:"id"`
	BarangID uint                    `json:"barang_id"`
	Qty      int                     `json:"qty"`
	Harga    float64                 `json:"harga"`
	Subtotal float64                 `json:"subtotal"`
	Barang   BarangPembelianResponse `json:"barang"`
}

type PembelianResponse struct {
	Header  BeliHeaderResponse   `json:"header"`
	Details []BeliDetailResponse `json:"details"`
}

type BarangPembelianResponse struct {
	KodeBarang string `json:"kode_barang"`
	NamaBarang string `json:"nama_barang"`
	Satuan     string `json:"satuan"`
}
