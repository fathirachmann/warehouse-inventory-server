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

// GetBarang godoc
// @Summary Get all barang
// @Description Get a list of barang with pagination and search
// @Tags Barang
// @Accept json
// @Produce json
// @Param search query string false "Search by name"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /barang [get]
// @Security BearerAuth
func (h *BarangHandler) GetBarang(c *fiber.Ctx) error {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	items, total, err := h.repo.List(search, page, limit)
	if err != nil {
		log.Println("Error fetching barang list:", err.Error(), "barang_handler.go:GetBarang", "Error at line 36")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}
	return c.Status(200).JSON(fiber.Map{
		"data": items,
		"meta": fiber.Map{"total": total, "page": page, "limit": limit},
	})
}

func (h *BarangHandler) GetBarangByID(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	barang, err := h.repo.GetByID(uint(id64))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Barang tidak Ditemukan")
	}
	return c.Status(200).JSON(barang)
}

func (h *BarangHandler) CreateBarang(c *fiber.Ctx) error {
	var payload models.MasterBarang
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Data input tidak valid")
	}

	errMap := make(map[string]string)

	switch {
	case payload.NamaBarang == "":
		errMap["nama_barang"] = "nama barang tidak boleh kosong"
	case payload.Satuan == "":
		errMap["satuan"] = "satuan tidak boleh kosong"
	case payload.HargaBeli <= 0:
		errMap["harga_beli"] = "harga beli harus lebih dari 0"
	case payload.HargaJual <= 0:
		errMap["harga_jual"] = "harga jual harus lebih dari 0"
	}

	if len(errMap) > 0 {
		return &middleware.ValidationError{
			Message: "validation error",
			Errors:  errMap,
		}
	}

	if err := h.repo.Create(&payload); err != nil {
		log.Println("Error creating barang:", err.Error(), "barang_handler.go:CreateBarang", "Error at line 89")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}
	return c.Status(fiber.StatusCreated).JSON(payload)
}

func (h *BarangHandler) UpdateBarangByID(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Parameter ID tidak valid")
	}
	barang, err := h.repo.GetByID(uint(id64))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Barang tidak Ditemukan")
	}
	var payload models.MasterBarang
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Data input tidak valid")
	}

	barang.KodeBarang = payload.KodeBarang
	barang.NamaBarang = payload.NamaBarang
	barang.Satuan = payload.Satuan
	barang.HargaBeli = payload.HargaBeli
	barang.HargaJual = payload.HargaJual

	if err := h.repo.Update(barang); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(200).JSON(barang)
}

func (h *BarangHandler) DeleteBarangByID(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Parameter ID tidak valid")
	}
	if err := h.repo.Delete(uint(id64)); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Barang tidak ditemukan")
	}

	message := fmt.Sprintf("Barang dengan ID %d berhasil dihapus", id64)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": message})
}
