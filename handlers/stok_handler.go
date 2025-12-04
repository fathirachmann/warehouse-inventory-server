package handlers

import (
	"log"
	"strconv"

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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"data": data,
	})
}

// GetStokByBarangID adalah handler untuk mendapatkan data stok berdasarkan barang_id
func (h *StokHandler) GetStokByBarangID(c *fiber.Ctx) error {
	barangIDStr := c.Params("barang_id")
	barangID64, err := strconv.ParseUint(barangIDStr, 10, 64)
	if err != nil {
		return c.Status(200).Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Data input tidak valid"})
	}

	stok, err := h.repo.GetOrCreateByBarangID(uint(barangID64))
	if err != nil {
		log.Println("Error fetching stok by barang ID:", err.Error(), "stok_handler.go:GetStokByBarangID", "Error at line 56")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": []any{stok},
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data":  data,
		"total": total,
	})
}

// GetHistoryByBarangID adalah handler untuk mendapatkan data history stok berdasarkan barang_id dengan pagination
func (h *StokHandler) GetHistoryByBarangID(c *fiber.Ctx) error {
	barangIDStr := c.Params("barang_id")
	barangID64, err := strconv.ParseUint(barangIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid barang_id"})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":  data,
		"total": total,
	})
}
