package repositories

import (
	"warehouse-inventory-server/models"

	"gorm.io/gorm"
)

type PenjualanRepository struct {
	db *gorm.DB
}

func NewPenjualanRepository(db *gorm.DB) *PenjualanRepository {
	return &PenjualanRepository{db: db}
}

// CreatePenjualan menyimpan header + detail penjualan dalam satu transaksi
func (r *PenjualanRepository) CreatePenjualan(header *models.JualHeader, details []models.JualDetail) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(header).Error; err != nil {
			return err
		}
		for i := range details {
			details[i].JualHeaderID = header.ID
		}
		if len(details) > 0 {
			if err := tx.Create(&details).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *PenjualanRepository) GetAllPenjualan() ([]models.JualHeader, error) {
	var headers []models.JualHeader
	if err := r.db.Preload("Details").Find(&headers).Error; err != nil {
		return nil, err
	}
	return headers, nil
}

func (r *PenjualanRepository) GetPenjualanByID(id uint) (*models.JualHeader, error) {
	var header models.JualHeader
	if err := r.db.Preload("Details").First(&header, id).Error; err != nil {
		return nil, err
	}
	return &header, nil
}
