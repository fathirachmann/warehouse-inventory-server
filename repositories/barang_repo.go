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
		if err := tx.Save(b).Error; err != nil {
			return err
		}

		// Create Mstok with default 0
		stok := models.Mstok{
			BarangID:  b.ID,
			StokAkhir: 0,
		}
		return tx.Create(&stok).Error
	})
}

func (r *BarangRepository) Update(b *models.MasterBarang) error {
	return r.db.Save(b).Error
}

func (r *BarangRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var stok models.Mstok
		// Check if stock exists
		if err := tx.Where("barang_id = ?", id).First(&stok).Error; err == nil {
			if stok.StokAkhir > 0 {
				return fmt.Errorf("stok barang berjumlah (%d), tidak dapat dihapus. Hanya bisa menghapus barang yang stok-nya sudah habis", stok.StokAkhir)
			}
			// Delete stock
			if err := tx.Delete(&stok).Error; err != nil {
				return err
			}
		}

		// Delete barang
		result := tx.Delete(&models.MasterBarang{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("barang tidak ditemukan")
		}
		return nil
	})
}

func (r *BarangRepository) GetByID(id uint) (*models.MasterBarang, error) {
	var b models.MasterBarang
	if err := r.db.First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BarangRepository) GetDetailByID(id uint) (*models.BarangWithStock, error) {
	var b models.BarangWithStock
	err := r.db.Table("master_barang").
		Select("master_barang.*, mstok.stok_akhir").
		Joins("JOIN mstok ON mstok.barang_id = master_barang.id").
		Where("master_barang.id = ?", id).
		First(&b).Error
	if err != nil {
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

func (r *BarangRepository) List(search string, page, limit int) ([]models.BarangWithStock, int64, error) {
	var items []models.BarangWithStock
	var total int64

	// Base query with join
	q := r.db.Table("master_barang").
		Select("master_barang.*, mstok.stok_akhir").
		Joins("JOIN mstok ON mstok.barang_id = master_barang.id")

	if search != "" {
		like := "%" + search + "%"
		q = q.Where("master_barang.kode_barang ILIKE ? OR master_barang.nama_barang ILIKE ?", like, like)
	}

	// Count total matching records
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

	if err := q.Order("master_barang.kode_barang ASC").Limit(limit).Offset(offset).Scan(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
