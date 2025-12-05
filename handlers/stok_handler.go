package handlers

import (
	"log"
	"strconv"

	"warehouse-inventory-server/models"
	"warehouse-inventory-server/repositories"

	"github.com/gofiber/fiber/v2"
)

type StokHandler struct {
	repo *repositories.StokRepository
}

func NewStokHandler(repo *repositories.StokRepository) *StokHandler {
	return &StokHandler{repo: repo}
}

// Route Handlers - Stock
func (h *StokHandler) RegisterStockRoute(r fiber.Router) {
	r.Get("/", h.GetAllStok)
	r.Get("/:barang_id", h.GetStokByBarangID)
}

// Route Handlers - History
func (h *StokHandler) RegisterHistoryRoute(r fiber.Router) {
	r.Get("/", h.GetHistoryAll)
	r.Get("/:barang_id", h.GetHistoryByBarangID)
}

// GetAllStok adalah handler untuk mendapatkan semua data stok
func (h *StokHandler) GetAllStok(c *fiber.Ctx) error {
	data, err := h.repo.GetAllStok()
	if err != nil {
		log.Println("Error fetching all stok:", err.Error(), "stok_handler.go:GetAllStok", "Error at line 36")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	var response []models.MstokResponse
	for _, item := range data {
		response = append(response, models.MstokResponse{
			ID:        item.ID,
			BarangID:  item.BarangID,
			StokAkhir: item.StokAkhir,
			UpdatedAt: item.UpdatedAt,
			Barang: models.BarangStokResponse{
				KodeBarang: item.MasterBarang.KodeBarang,
				NamaBarang: item.MasterBarang.NamaBarang,
				Satuan:     item.MasterBarang.Satuan,
				HargaJual:  item.MasterBarang.HargaJual,
			},
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": response,
	})
}

// GetStokByBarangID adalah handler untuk mendapatkan data stok berdasarkan barang_id
func (h *StokHandler) GetStokByBarangID(c *fiber.Ctx) error {
	barangIDStr := c.Params("barang_id")
	barangID64, err := strconv.ParseUint(barangIDStr, 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Parameter barang_id tidak valid")
	}

	stok, err := h.repo.GetOrCreateByBarangID(uint(barangID64))
	if err != nil {
		log.Println("Error fetching stok by barang ID:", err.Error(), "stok_handler.go:GetStokByBarangID", "Error at line 56")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	response := models.MstokResponse{
		ID:        stok.ID,
		BarangID:  stok.BarangID,
		StokAkhir: stok.StokAkhir,
		UpdatedAt: stok.UpdatedAt,
		Barang: models.BarangStokResponse{
			KodeBarang: stok.MasterBarang.KodeBarang,
			NamaBarang: stok.MasterBarang.NamaBarang,
			Satuan:     stok.MasterBarang.Satuan,
			HargaJual:  stok.MasterBarang.HargaJual,
		},
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": []models.MstokResponse{response},
	})
}

// GetHistoryAll adalah handler untuk mendapatkan semua data history stok dengan pagination
func (h *StokHandler) GetHistoryAll(c *fiber.Ctx) error {

	// Note: Pagination tidak ada di requirement. Namun, untuk frontend agar tidak load data terlalu banyak.

	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	data, total, err := h.repo.GetHistory(0, limit, offset)
	if err != nil {
		log.Println("Error fetching all history stok:", err.Error(), "stok_handler.go:GetHistoryAll", "Error at line 87")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	var response []models.HistoryStokResponse
	for _, item := range data {
		response = append(response, models.HistoryStokResponse{
			ID:             item.ID,
			BarangID:       item.BarangID,
			UserID:         item.UserID,
			JenisTransaksi: item.JenisTransaksi,
			Jumlah:         item.Jumlah,
			StokSebelumnya: item.StokSebelumnya,
			StokSesudah:    item.StokSesudah,
			Keterangan:     item.Keterangan,
			CreatedAt:      item.CreatedAt,
			Barang: models.BarangSimpleResponse{
				KodeBarang: item.MasterBarang.KodeBarang,
				NamaBarang: item.MasterBarang.NamaBarang,
			},
			User: models.UserSimpleResponse{
				Username: item.Users.Username,
				FullName: item.Users.FullName,
			},
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data":  response,
		"total": total,
	})
}

// GetHistoryByBarangID adalah handler untuk mendapatkan data history stok berdasarkan barang_id dengan pagination
func (h *StokHandler) GetHistoryByBarangID(c *fiber.Ctx) error {
	barangIDStr := c.Params("barang_id")
	barangID64, err := strconv.ParseUint(barangIDStr, 10, 64)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid barang_id")
	}

	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	data, total, err := h.repo.GetHistory(uint(barangID64), limit, offset)
	if err != nil {
		log.Println("Error fetching history by barang ID:", err.Error(), "stok_handler.go:GetHistoryByBarangID", "Error at line 123")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	var response []models.HistoryStokResponse
	for _, item := range data {
		response = append(response, models.HistoryStokResponse{
			ID:             item.ID,
			BarangID:       item.BarangID,
			UserID:         item.UserID,
			JenisTransaksi: item.JenisTransaksi,
			Jumlah:         item.Jumlah,
			StokSebelumnya: item.StokSebelumnya,
			StokSesudah:    item.StokSesudah,
			Keterangan:     item.Keterangan,
			CreatedAt:      item.CreatedAt,
			Barang: models.BarangSimpleResponse{
				KodeBarang: item.MasterBarang.KodeBarang,
				NamaBarang: item.MasterBarang.NamaBarang,
			},
			User: models.UserSimpleResponse{
				Username: item.Users.Username,
				FullName: item.Users.FullName,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
		"meta": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}
