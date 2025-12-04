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

func (r *StokRepository) GetOrCreateByBarangID(barangID uint) (*models.Mstok, error) {
	var stok models.Mstok

	err := r.db.Where("barang_id = ?", barangID).First(&stok).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		stok = models.Mstok{
			BarangID:  barangID,
			StokAkhir: 0,
		}
		if errCreate := r.db.Create(&stok).Error; errCreate != nil {
			return nil, errCreate
		}
		return &stok, nil
	}
	if err != nil {
		return nil, err
	}

	return &stok, nil
}

func (r *StokRepository) UpdateStok(stok *models.Mstok) error {
	return r.db.Save(stok).Error
}

func (r *StokRepository) CreateHistory(history *models.HistoryStok) error {
	return r.db.Create(history).Error
}

func (r *StokRepository) GetAllStok() ([]models.Mstok, error) {
	var list []models.Mstok
	if err := r.db.Preload("MasterBarang").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *StokRepository) GetHistory(barangID uint, limit, offset int) ([]models.HistoryStok, error) {
	var list []models.HistoryStok
	q := r.db.Preload("MasterBarang").Preload("Users").Order("created_at DESC")
	if barangID != 0 {
		q = q.Where("barang_id = ?", barangID)
	}
	if err := q.Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
