package handlers

import (
	"fmt"
	"log"
	"time"

	"warehouse-inventory-server/middleware"
	"warehouse-inventory-server/models"
	"warehouse-inventory-server/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type PembelianHandler struct {
	repo       *repositories.PembelianRepository
	stokRepo   *repositories.StokRepository
	barangRepo *repositories.BarangRepository
}

func NewPembelianHandler(repo *repositories.PembelianRepository, stokRepo *repositories.StokRepository, barangRepo *repositories.BarangRepository) *PembelianHandler {
	return &PembelianHandler{
		repo:       repo,
		stokRepo:   stokRepo,
		barangRepo: barangRepo,
	}
}

// RegisterRoute mendaftarkan seluruh endpoint "/api/pembelian"
func (h *PembelianHandler) RegisterRoute(r fiber.Router) {
	r.Post("/", h.CreatePembelian)
	r.Get("/", h.GetAllPembelian)
	r.Get("/:id", h.GetPembelianByID)
}

// CreatePembelian godoc
// @Summary Create new purchase
// @Description Create a new purchase transaction
// @Tags Pembelian
// @Accept json
// @Produce json
// @Param body body models.BeliHeaderRequest true "Purchase Request"
// @Success 201 {object} models.PembelianResponse "Created"
// @Failure 400 {object} middleware.ValidationError "Bad Request"
// @Failure 422 {object} middleware.ValidationError "Unprocessable Entity"
// @Failure 500 {object} middleware.ValidationError "Internal Server Error"
// @Security BearerAuth
// @Router /api/pembelian [post]
func (h *PembelianHandler) CreatePembelian(c *fiber.Ctx) error {
	var req models.BeliHeaderRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Data input tidak valid")
	}

	errMap := make(map[string]string)

	switch {
	case req.Supplier == "":
		errMap["supplier"] = "Nama supplier tidak boleh kosong"
	case len(req.Details) == 0:
		errMap["details"] = "details tidak boleh kosong"
	}

	if len(errMap) > 0 {
		return &middleware.ValidationError{
			Message: "validation error",
			Errors:  errMap,
		}
	}

	var userID uint
	if claims, ok := c.Locals("user").(jwt.MapClaims); ok {
		if sub, ok := claims["id"]; ok {
			switch v := sub.(type) {
			case float64:
				userID = uint(v)
			case int:
				userID = uint(v)
			}
		}
	}

	header := models.BeliHeader{
		Supplier:  req.Supplier,
		UserID:    userID,
		Status:    "selesai",
		CreatedAt: time.Now(),
	}

	var details []models.BeliDetail
	total := 0.0
	for _, d := range req.Details {
		if d.Qty <= 0 || d.Harga <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "qty dan harga tidak boleh kurang dari sama dengan 0")
		}

		// Validasi harga beli sesuai dengan harga di master barang
		barang, err := h.barangRepo.GetByID(d.BarangID)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Barang tidak ditemukan")
		}
		if d.Harga != barang.HargaBeli {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("Harga beli %s tidak sesuai (Expected: %.2f)", barang.NamaBarang, barang.HargaJual))
		}

		subtotal := float64(d.Qty) * d.Harga
		total += subtotal
		detail := models.BeliDetail{
			BarangID: d.BarangID,
			Qty:      d.Qty,
			Harga:    d.Harga,
			Subtotal: subtotal,
		}
		details = append(details, detail)
	}
	header.Total = total

	if err := h.repo.CreatePembelian(&header, details); err != nil {
		log.Println("Error CreatePembelian:", err.Error(), "pembelian_handler.go:CreatePembelian", "Error at line 119")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	created, err := h.repo.GetPembelianByID(header.ID)
	if err != nil {
		log.Println("Error fetching created pembelian:", err.Error(), "pembelian_handler.go:CreatePembelian", "Error at line 124")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	response := mapToPembelianResponse(created)
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetAllPembelian godoc
// @Summary Get all purchases
// @Description Get a list of all purchase transactions
// @Tags Pembelian
// @Produce json
// @Success 200 {object} models.PembelianResponse "OK"
// @Failure 500 {object} middleware.ValidationError "Internal Server Error"
// @Security BearerAuth
// @Router /api/pembelian [get]
func (h *PembelianHandler) GetAllPembelian(c *fiber.Ctx) error {
	data, err := h.repo.GetAllPembelian()
	if err != nil {
		log.Println("Error fetching all pembelian:", err.Error(), "pembelian_handler.go:GetAllPembelian", "Error at line 144")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	var response []models.PembelianResponse
	for _, p := range data {
		response = append(response, mapToPembelianResponse(&p))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}

// GetPembelianByID godoc
// @Summary Get purchase by ID
// @Description Get details of a specific purchase transaction
// @Tags Pembelian
// @Produce json
// @Param id path int true "Purchase ID"
// @Success 200 {object} models.PembelianResponse
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/pembelian/{id} [get]
func (h *PembelianHandler) GetPembelianByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Data input tidak valid")
	}
	data, err := h.repo.GetPembelianByID(uint(id))
	if err != nil {
		log.Println("Error fetching pembelian by ID:", err.Error(), "pembelian_handler.go:GetPembelianByID", "Error at line 176")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	response := mapToPembelianResponse(data)
	return c.Status(fiber.StatusOK).JSON(response)
}

// Private helper functions untuk mapping struct response
func mapToPembelianResponse(p *models.BeliHeader) models.PembelianResponse {
	details := make([]models.BeliDetailResponse, len(p.Details))
	for i, d := range p.Details {
		details[i] = models.BeliDetailResponse{
			ID:       d.ID,
			BarangID: d.BarangID,
			Barang: models.BarangPembelianResponse{
				KodeBarang: d.MasterBarang.KodeBarang,
				NamaBarang: d.MasterBarang.NamaBarang,
				Satuan:     d.MasterBarang.Satuan,
			},
			Qty:      d.Qty,
			Harga:    d.Harga,
			Subtotal: d.Subtotal,
		}
	}

	return models.PembelianResponse{
		Header: models.BeliHeaderResponse{
			ID:        p.ID,
			NoFaktur:  p.NoFaktur,
			UserID:    p.UserID,
			Supplier:  p.Supplier,
			Status:    p.Status,
			User:      models.UserSimpleResponse{Username: p.User.Username, FullName: p.User.FullName},
			Total:     p.Total,
			CreatedAt: p.CreatedAt,
		},
		Details: details,
	}
}
