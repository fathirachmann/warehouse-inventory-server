package repositories

import (
	"errors"
	"fmt"

	"warehouse-inventory-server/models"

	"gorm.io/gorm"
)

type PenjualanRepository struct {
	db *gorm.DB
}

func NewPenjualanRepository(db *gorm.DB) *PenjualanRepository {
	return &PenjualanRepository{db: db}
}

// CreatePenjualan adalah method untuk menyimpan header + detail penjualan dalam satu transaksi
func (r *PenjualanRepository) CreatePenjualan(header *models.JualHeader, details []models.JualDetail) error {
	// Mulai transaksi
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Validasi stok sebelum melakukan perubahan
	for _, d := range details {
		var stok models.Mstok
		if err := tx.Where("barang_id = ?", d.BarangID).First(&stok).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return errors.New("stok tidak mencukupi")
			}
			tx.Rollback()
			return err
		}
		if stok.StokAkhir < d.Qty {
			tx.Rollback()
			return errors.New("stok tidak mencukupi")
		}
	}

	// Buat header penjualan untuk mendapatkan ID
	if err := tx.Create(header).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Generate NoFaktur berdasarkan ID: JUAL + 3 digit (misal JUAL001)
	header.NoFaktur = fmt.Sprintf("JUAL%03d", header.ID)
	if err := tx.Model(header).Update("no_faktur", header.NoFaktur).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update stok & buat history untuk setiap detail (stok keluar)
	for i := range details {
		details[i].JualHeaderID = header.ID
		// ambil stok terbaru dalam transaksi
		var stok models.Mstok
		if err := tx.Where("barang_id = ?", details[i].BarangID).First(&stok).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return errors.New("stok tidak ditemukan")
			}
			tx.Rollback()
			return err
		}
		stokSebelum := stok.StokAkhir
		stokSesudah := stokSebelum - details[i].Qty
		if stokSesudah < 0 {
			tx.Rollback()
			return errors.New("stok tidak mencukupi")
		}
		stok.StokAkhir = stokSesudah
		if err := tx.Save(&stok).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Buat history stok penjualan
		history := models.HistoryStok{
			BarangID:       details[i].BarangID,
			UserID:         header.UserID,
			JenisTransaksi: "keluar",
			Jumlah:         details[i].Qty,
			StokSebelum:    stokSebelum,
			StokSesudah:    stokSesudah,
			Keterangan:     "Penjualan " + header.NoFaktur,
		}
		if err := tx.Create(&history).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Buat detail penjualan
	if len(details) > 0 {
		if err := tx.Create(&details).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Jika semua operasi berhasil, commit transaksi
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// GetAllPenjualan mengambil semua data penjualan beserta detailnya
func (r *PenjualanRepository) GetAllPenjualan() ([]models.JualHeader, error) {
	var headers []models.JualHeader
	if err := r.db.Preload("Details.MasterBarang").Preload("User").Order("created_at desc").Find(&headers).Error; err != nil {
		return nil, err
	}
	return headers, nil
}

// GetPenjualanByID mengambil data penjualan berdasarkan ID beserta detailnya
func (r *PenjualanRepository) GetPenjualanByID(id uint) (*models.JualHeader, error) {
	var header models.JualHeader
	if err := r.db.Preload("Details.MasterBarang").Preload("User").First(&header, id).Error; err != nil {
		return nil, err
	}
	return &header, nil
}
