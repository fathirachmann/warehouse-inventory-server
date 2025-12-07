package repositories

import (
	"warehouse-inventory-server/models"

	"gorm.io/gorm"
)

type StokRepository struct {
	db *gorm.DB
}

func NewStokRepository(db *gorm.DB) *StokRepository {
	return &StokRepository{db: db}
}

// GetByBarangID mengambil stok berdasarkan barangID
func (r *StokRepository) GetByBarangID(barangID uint) (*models.Mstok, error) {
	var stok models.Mstok
	err := r.db.Preload("MasterBarang").Where("barang_id = ?", barangID).First(&stok).Error
	if err != nil {
		return nil, err
	}
	return &stok, nil
}

// CreateStok membuat entri stok baru untuk barangID
func (r *StokRepository) CreateStok(barangID uint) (*models.Mstok, error) {
	stok := models.Mstok{
		BarangID:  barangID,
		StokAkhir: 0,
	}
	if err := r.db.Create(&stok).Error; err != nil {
		return nil, err
	}
	// Load relation for the newly created record
	if err := r.db.Preload("MasterBarang").First(&stok, stok.ID).Error; err != nil {
		// If loading fails, just return the stok without relation
		return &stok, nil
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
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
