package handlers

import (
	"fmt"
	"log"
	"strconv"

	"warehouse-inventory-server/middleware"
	"warehouse-inventory-server/models"
	"warehouse-inventory-server/repositories"

	"github.com/gofiber/fiber/v2"
)

type BarangHandler struct {
	repo *repositories.BarangRepository
}

func NewBarangHandler(repo *repositories.BarangRepository) *BarangHandler {
	return &BarangHandler{repo: repo}
}

func (h *BarangHandler) RegisterRoute(r fiber.Router) {
	r.Get("/", h.GetBarang)
	r.Get("/:id", h.GetBarangByID)
	r.Post("/", middleware.GuardAdmin(), h.CreateBarang)
	r.Put("/:id", middleware.GuardAdmin(), h.UpdateBarangByID)
	r.Delete("/:id", middleware.GuardAdmin(), h.DeleteBarangByID)
}

func (h *BarangHandler) GetBarang(c *fiber.Ctx) error {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	items, total, err := h.repo.List(search, page, limit)
	if err != nil {
		log.Println("Error fetching barang list:", err.Error(), "barang_handler.go:GetBarang", "Error at line 36")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"data": items,
		"meta": fiber.Map{"total": total, "page": page, "limit": limit},
	})
}

func (h *BarangHandler) GetBarangByID(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid id"})
	}
	barang, err := h.repo.GetByID(uint(id64))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Barang tidak Ditemukan"})
	}
	return c.Status(200).JSON(barang)
}

func (h *BarangHandler) CreateBarang(c *fiber.Ctx) error {
	var payload models.MasterBarang
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   "body parsing error",
			"message": "Data input tidak valid",
		})
	}

	errMap := make(map[string][]string)

	switch {
	case payload.NamaBarang == "":
		errMap["nama_barang"] = append(errMap["nama_barang"], "nama barang tidak boleh kosong")
	case payload.Satuan == "":
		errMap["satuan"] = append(errMap["satuan"], "satuan tidak boleh kosong")
	case payload.HargaBeli <= 0:
		errMap["harga_beli"] = append(errMap["harga_beli"], "harga beli harus lebih dari 0")
	case payload.HargaJual <= 0:
		errMap["harga_jual"] = append(errMap["harga_jual"], "harga jual harus lebih dari 0")
	}

	if len(errMap) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "validation error",
			"message": errMap,
		})
	}

	if err := h.repo.Create(&payload); err != nil {
		log.Println("Error creating barang:", err.Error(), "barang_handler.go:CreateBarang", "Error at line 89")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server error",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(payload)
}

func (h *BarangHandler) UpdateBarangByID(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   "Data input tidak valid",
			"message": "Parameter ID tidak valid",
		})
	}
	barang, err := h.repo.GetByID(uint(id64))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "database error",
			"message": "Barang tidak Ditemukan",
		})
	}
	var payload models.MasterBarang
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   "body parsing error",
			"message": "Data input tidak valid",
		})
	}

	barang.KodeBarang = payload.KodeBarang
	barang.NamaBarang = payload.NamaBarang
	barang.Satuan = payload.Satuan
	barang.HargaBeli = payload.HargaBeli
	barang.HargaJual = payload.HargaJual

	if err := h.repo.Update(barang); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "database error",
			"message": err.Error(),
		})
	}
	return c.Status(200).JSON(barang)
}

func (h *BarangHandler) DeleteBarangByID(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   "Data input tidak valid",
			"message": "Parameter ID tidak valid",
		})
	}
	if err := h.repo.Delete(uint(id64)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "database error",
			"message": "Barang tidak ditemukan",
		})
	}

	message := fmt.Sprintf("Barang dengan ID %d berhasil dihapus", id64)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": message})
}
