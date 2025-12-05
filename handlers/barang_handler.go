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
// @Description Mendapatkan daftar seluruh barang
// @Tags Barang
// @Accept json
// @Produce json
// @Param search query string false "Search by name"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} models.BarangResponse "OK"
// @Failure 500 {object} middleware.ValidationError "Internal Server Error"
// @Router /api/barang [get]
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

	var response []models.BarangResponse
	for _, item := range items {
		response = append(response, models.BarangResponse{
			ID:         item.ID,
			KodeBarang: item.KodeBarang,
			NamaBarang: item.NamaBarang,
			Deskripsi:  item.Deskripsi,
			Satuan:     item.Satuan,
			HargaBeli:  item.HargaBeli,
			HargaJual:  item.HargaJual,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": response,
		"meta": fiber.Map{"total": total, "page": page, "limit": limit},
	})
}

// GetBarangByID godoc
// @Summary Get barang by ID
// @Description Mendapatkan detail barang berdasarkan ID
// @Tags Barang
// @Accept json
// @Produce json
// @Param id path int true "Barang ID"
// @Success 200 {object} models.BarangResponse "OK"
// @Failure 404 {object} middleware.SpecificErrorResponse "Not Found"
// @Router /api/barang/{id} [get]
// @Security BearerAuth
func (h *BarangHandler) GetBarangByID(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}
	barang, err := h.repo.GetByID(uint(id64))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Barang tidak Ditemukan")
	}

	response := models.BarangResponse{
		ID:         barang.ID,
		KodeBarang: barang.KodeBarang,
		NamaBarang: barang.NamaBarang,
		Deskripsi:  barang.Deskripsi,
		Satuan:     barang.Satuan,
		HargaBeli:  barang.HargaBeli,
		HargaJual:  barang.HargaJual,
	}

	return c.Status(200).JSON(response)
}

// CreateBarang godoc
// @Summary Create new barang
// @Description Membuat barang baru
// @Tags Barang
// @Accept json
// @Produce json
// @Param body body models.BarangRequest true "Barang Request"
// @Success 201 {object} models.CreatedBarangResponse "Created"
// @Failure 400 {object} middleware.SpecificErrorResponse "Bad Request"
// @Failure 422 {object} middleware.ValidationError "Unprocessable Entity"
// @Failure 500 {object} middleware.ErrorResponse "Internal Server Error"
// @Router /api/barang [post]
// @Security BearerAuth
func (h *BarangHandler) CreateBarang(c *fiber.Ctx) error {
	var req models.BarangRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Data input tidak valid")
	}

	errMap := make(map[string]string)

	switch {
	case req.NamaBarang == "":
		errMap["nama_barang"] = "nama barang tidak boleh kosong"
	case req.Deskripsi == "":
		errMap["deskripsi"] = "deskripsi tidak boleh kosong"
	case req.Satuan == "":
		errMap["satuan"] = "satuan tidak boleh kosong"
	case req.HargaBeli <= 0:
		errMap["harga_beli"] = "harga beli harus lebih dari 0"
	case req.HargaJual <= 0:
		errMap["harga_jual"] = "harga jual harus lebih dari 0"
	}

	if len(errMap) > 0 {
		return &middleware.ValidationError{
			Message: "validation error",
			Errors:  errMap,
		}
	}

	barang := models.MasterBarang{
		NamaBarang: req.NamaBarang,
		Deskripsi:  req.Deskripsi,
		Satuan:     req.Satuan,
		HargaBeli:  req.HargaBeli,
		HargaJual:  req.HargaJual,
	}

	if err := h.repo.Create(&barang); err != nil {
		log.Println("Error creating barang:", err.Error(), "barang_handler.go:CreateBarang", "Error at line 89")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	response := models.CreatedBarangResponse{
		ID:         barang.ID,
		KodeBarang: barang.KodeBarang,
		NamaBarang: barang.NamaBarang,
		Deskripsi:  barang.Deskripsi,
		Satuan:     barang.Satuan,
		HargaBeli:  barang.HargaBeli,
		HargaJual:  barang.HargaJual,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// UpdateBarangByID godoc
// @Summary Update barang by ID
// @Description Memperbarui detail barang berdasarkan ID
// @Tags Barang
// @Accept json
// @Produce json
// @Param id path int true "Barang ID"
// @Param body body models.BarangRequest true "Barang Request"
// @Success 200 {object} models.BarangResponse "OK"
// @Failure 400 {object} middleware.SpecificErrorResponse "Bad Request"
// @Failure 404 {object} middleware.SpecificErrorResponse "Not Found"
// @Failure 422 {object} middleware.ValidationError "Unprocessable Entity"
// @Failure 500 {object} middleware.ErrorResponse "Internal Server Error"
// @Router /api/barang/{id} [put]
// @Security BearerAuth
func (h *BarangHandler) UpdateBarangByID(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Parameter ID tidak valid")
	}
	barang, err := h.repo.GetByID(uint(id64))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Barang tidak Ditemukan")
	}

	var req models.BarangRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Data input tidak valid")
	}

	errMap := make(map[string]string)

	switch {
	case req.NamaBarang == "":
		errMap["nama_barang"] = "nama barang tidak boleh kosong"
	case req.Deskripsi == "":
		errMap["deskripsi"] = "deskripsi tidak boleh kosong"
	case req.Satuan == "":
		errMap["satuan"] = "satuan tidak boleh kosong"
	case req.HargaBeli <= 0:
		errMap["harga_beli"] = "harga beli harus lebih dari 0"
	case req.HargaJual <= 0:
		errMap["harga_jual"] = "harga jual harus lebih dari 0"
	}

	if len(errMap) > 0 {
		return &middleware.ValidationError{
			Message: "validation error",
			Errors:  errMap,
		}
	}

	barang.NamaBarang = req.NamaBarang
	barang.Deskripsi = req.Deskripsi
	barang.Satuan = req.Satuan
	barang.HargaBeli = req.HargaBeli
	barang.HargaJual = req.HargaJual

	if err := h.repo.Update(barang); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := models.BarangResponse{
		ID:         barang.ID,
		KodeBarang: barang.KodeBarang,
		NamaBarang: barang.NamaBarang,
		Deskripsi:  barang.Deskripsi,
		Satuan:     barang.Satuan,
		HargaBeli:  barang.HargaBeli,
		HargaJual:  barang.HargaJual,
	}

	return c.Status(200).JSON(response)
}

// DeleteBarangByID godoc
// @Summary Delete barang by ID
// @Description Menghapus barang berdasarkan ID
// @Tags Barang
// @Accept json
// @Produce json
// @Param id path int true "Barang ID"
// @Success 200 {object} models.DeleteBarangResponse "OK"
// @Failure 422 {object} middleware.ErrorResponse "Unprocessable Entity"
// @Failure 404 {object} middleware.ErrorResponse "Not Found"
// @Failure 500 {object} middleware.ErrorResponse "Internal Server Error"
// @Router /api/barang/{id} [delete]
// @Security BearerAuth
func (h *BarangHandler) DeleteBarangByID(c *fiber.Ctx) error {
	id64, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Parameter ID tidak valid")
	}
	if err := h.repo.Delete(uint(id64)); err != nil {
		return fiber.NewError(fiber.StatusNotFound, "Barang tidak ditemukan")
	}

	message := fmt.Sprintf("Barang dengan ID %d berhasil dihapus", id64)

	return c.Status(fiber.StatusOK).JSON(models.DeleteBarangResponse{
		Message: message,
	})
}
