package repositories

import (
	"fmt"
	"warehouse-inventory-server/models"

	"gorm.io/gorm"
)

type BarangRepository struct {
	db *gorm.DB
}

func NewBarangRepository(db *gorm.DB) *BarangRepository {
	return &BarangRepository{db: db}
}

func (r *BarangRepository) Create(b *models.MasterBarang) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(b).Error; err != nil {
			return err
		}
		// Auto generate KodeBarang: BRG + ID (e.g. BRG001)
		b.KodeBarang = fmt.Sprintf("BRG%03d", b.ID)
		return tx.Save(b).Error
	})
}

func (r *BarangRepository) Update(b *models.MasterBarang) error {
	return r.db.Save(b).Error
}

func (r *BarangRepository) Delete(id uint) error {
	return r.db.Delete(&models.MasterBarang{}, id).Error
}

func (r *BarangRepository) GetByID(id uint) (*models.MasterBarang, error) {
	var b models.MasterBarang
	if err := r.db.First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BarangRepository) GetByKode(kode string) (*models.MasterBarang, error) {
	var b models.MasterBarang
	if err := r.db.Where("kode_barang = ?", kode).First(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BarangRepository) List(search string, page, limit int) ([]models.MasterBarang, int64, error) {
	var items []models.MasterBarang
	var total int64
	q := r.db.Model(&models.MasterBarang{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("kode_barang ILIKE ? OR nama_barang ILIKE ?", like, like)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit
	if err := q.Order("nama_barang ASC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
