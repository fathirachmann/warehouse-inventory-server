package handlers

import (
	"fmt"
	"log"
	"strconv"

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
	r.Post("/", h.CreateBarang)
	r.Put("/:id", h.UpdateBarangByID)
	r.Delete("/:id", h.DeleteBarangByID)
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	barang, err := h.repo.GetByID(uint(id64))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Barang tidak Ditemukan"})
	}
	return c.Status(200).JSON(barang)
}

func (h *BarangHandler) CreateBarang(c *fiber.Ctx) error {
	var payload models.MasterBarang
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Data input tidak valid"})
	}

	switch {
	case payload.KodeBarang == "":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "kode barang tidak boleh kosong"})
	case payload.NamaBarang == "":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "nama barang tidak boleh kosong"})
	case payload.Satuan == "":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "satuan tidak boleh kosong"})
	}

	if err := h.repo.Create(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Data input tidak valid"})
	}
	return c.Status(fiber.StatusCreated).JSON(payload)
}

func (h *BarangHandler) UpdateBarangByID(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Data input tidak valid"})
	}
	barang, err := h.repo.GetByID(uint(id64))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Barang tidak Ditemukan"})
	}
	var payload models.MasterBarang
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Data input tidak valid"})
	}
	barang.KodeBarang = payload.KodeBarang
	barang.NamaBarang = payload.NamaBarang
	barang.Satuan = payload.Satuan
	barang.HargaBeli = payload.HargaBeli
	barang.HargaJual = payload.HargaJual
	if err := h.repo.Update(barang); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(200).JSON(barang)
}

func (h *BarangHandler) DeleteBarangByID(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Data input tidak valid"})
	}
	if err := h.repo.Delete(uint(id64)); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Barang tidak Ditemukan"})
	}

	message := fmt.Sprintf("Barang dengan ID %d berhasil dihapus", id64)

	return c.Status(200).JSON(fiber.Map{"message": message})
}
