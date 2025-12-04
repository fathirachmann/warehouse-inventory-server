package repositories

import (
	"errors"
	"fmt"

	"warehouse-inventory-server/models"

	"gorm.io/gorm"
)

type PembelianRepository struct {
	db *gorm.DB
}

func NewPembelianRepository(db *gorm.DB) *PembelianRepository {
	return &PembelianRepository{db: db}
}

// CreatePembelian membuat pembelian baru beserta update stok dan history
func (r *PembelianRepository) CreatePembelian(header *models.BeliHeader, details []models.BeliDetail) error {
	// Mulai transaksi
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Update stok dan buat history untuk setiap detail pembelian
	for i := range details {
		var stok models.Mstok
		if err := tx.Where("barang_id = ?", details[i].BarangID).First(&stok).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				stok = models.Mstok{
					BarangID:  details[i].BarangID,
					StokAkhir: 0,
				}
				if errCreate := tx.Create(&stok).Error; errCreate != nil {
					tx.Rollback()
					return errCreate
				}
			} else {
				tx.Rollback()
				return err
			}
		}

		// Update stok
		stokSebelum := stok.StokAkhir
		stokSesudah := stokSebelum + details[i].Qty
		stok.StokAkhir = stokSesudah
		if err := tx.Save(&stok).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Buat history stok
		history := models.HistoryStok{
			BarangID:       details[i].BarangID,
			UserID:         header.UserID,
			JenisTransaksi: "masuk",
			Jumlah:         details[i].Qty,
			StokSebelumnya: stokSebelum,
			StokSesudah:    stokSesudah,
			Keterangan:     "Pembelian " + header.NoFaktur,
		}
		if err := tx.Create(&history).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Simpan header pembelian terlebih dahulu untuk mendapatkan ID
	if err := tx.Create(header).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Generate NoFaktur berdasarkan ID: BLI + 3 digit (misal BLI001)
	header.NoFaktur = fmt.Sprintf("BLI%03d", header.ID)
	if err := tx.Model(header).Update("no_faktur", header.NoFaktur).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Set BeliHeaderID untuk setiap detail pembelian
	for i := range details {
		details[i].BeliHeaderID = header.ID
	}
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

// GetAllPembelian mengambil semua data pembelian beserta detailnya
func (r *PembelianRepository) GetAllPembelian() ([]models.BeliHeader, error) {
	var headers []models.BeliHeader
	if err := r.db.Preload("Details.MasterBarang").Preload("User").Find(&headers).Error; err != nil {
		return nil, err
	}
	return headers, nil
}

// GetPembelianByID mengambil data pembelian berdasarkan ID beserta detailnya
func (r *PembelianRepository) GetPembelianByID(id uint) (*models.BeliHeader, error) {
	var header models.BeliHeader
	if err := r.db.Preload("Details.MasterBarang").Preload("User").First(&header, id).Error; err != nil {
		return nil, err
	}
	return &header, nil
}
