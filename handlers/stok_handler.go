package handlers

import (
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

func (h *StokHandler) GetAllStok(c *fiber.Ctx) error {
	data, err := h.repo.GetAllStok()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
			"error":  err.Error(),
		})
	}
	return c.Status(200).JSON(data)
}

func (h *StokHandler) GetStokByBarangID(c *fiber.Ctx) error {
	barangIDStr := c.Params("barang_id")
	barangID64, err := strconv.ParseUint(barangIDStr, 10, 64)
	if err != nil {
		return c.Status(200).Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Data input tidak valid"})
	}

	stok, err := h.repo.GetOrCreateByBarangID(uint(barangID64))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
			"error":  err.Error(),
		})
	}

	return c.Status(200).JSON(stok)
}

// History handler methods
func (h *StokHandler) GetHistoryAll(c *fiber.Ctx) error {
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
	data, err := h.repo.GetHistory(0, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
			"error":  err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"page":  page,
		"limit": limit,
		"data":  data,
	})
}

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
	data, err := h.repo.GetHistory(uint(barangID64), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"page":      page,
		"limit":     limit,
		"barang_id": barangID64,
		"data":      data,
	})
}
