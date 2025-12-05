package repositories

import (
	"errors"

	"warehouse-inventory-server/models"

	"gorm.io/gorm"
)

type StokRepository struct {
	db *gorm.DB
}

func NewStokRepository(db *gorm.DB) *StokRepository {
	return &StokRepository{db: db}
}

// GetOrCreateByBarangID mengambil stok berdasarkan barangID, atau membuat entri baru jika tidak ada
func (r *StokRepository) GetOrCreateByBarangID(barangID uint) (*models.Mstok, error) {
	var stok models.Mstok

	err := r.db.Preload("MasterBarang").Where("barang_id = ?", barangID).First(&stok).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		stok = models.Mstok{
			BarangID:  barangID,
			StokAkhir: 0,
		}
		if errCreate := r.db.Create(&stok).Error; errCreate != nil {
			return nil, errCreate
		}
		// Load relation for the newly created record
		if errLoad := r.db.Preload("MasterBarang").First(&stok, stok.ID).Error; errLoad != nil {
			// If loading fails, just return the stok without relation (or handle error)
			// For now, we can ignore or log, but returning is fine.
		}
		return &stok, nil
	}
	if err != nil {
		return nil, err
	}

	return &stok, nil
}

// UpdateStok memperbarui data stok di database
func (r *StokRepository) UpdateStok(stok *models.Mstok) error {
	return r.db.Save(stok).Error
}

// CreateHistory menambahkan entri history stok baru
func (r *StokRepository) CreateHistory(history *models.HistoryStok) error {
	return r.db.Create(history).Error
}

// GetAllStok mengambil semua data stok beserta relasi MasterBarang
func (r *StokRepository) GetAllStok() ([]models.Mstok, error) {
	var list []models.Mstok
	if err := r.db.Preload("MasterBarang").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetHistory mengambil data history stok dan total count
func (r *StokRepository) GetHistory(barangID uint, limit, offset int) ([]models.HistoryStok, int64, error) {
	var list []models.HistoryStok
	var total int64

	base := r.db.Model(&models.HistoryStok{})
	if barangID != 0 {
		base = base.Where("barang_id = ?", barangID)
	}
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	q := r.db.Preload("MasterBarang").Preload("Users").Order("created_at DESC")
	if barangID != 0 {
		q = q.Where("barang_id = ?", barangID)
	}
	if err := q.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
