package handlers

import (
	"log"
	"time"

	"warehouse-inventory-server/middleware"
	"warehouse-inventory-server/models"
	"warehouse-inventory-server/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type PenjualanHandler struct {
	repo     *repositories.PenjualanRepository
	stokRepo *repositories.StokRepository
}

func NewPenjualanHandler(repo *repositories.PenjualanRepository, stokRepo *repositories.StokRepository) *PenjualanHandler {
	return &PenjualanHandler{
		repo:     repo,
		stokRepo: stokRepo,
	}
}

func (h *PenjualanHandler) RegisterRoute(r fiber.Router) {
	r.Post("/", h.CreatePenjualan)
	r.Get("/", h.GetAllPenjualan)
	r.Get("/:id", h.GetPenjualanByID)
}

// CreatePenjualan godoc
// @Summary Create new sale
// @Description Create a new sale transaction
// @Tags Penjualan
// @Accept json
// @Produce json
// @Param body body models.JualHeaderRequest true "Sale Request"
// @Success 201 {object} models.PenjualanResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /penjualan [post]
func (h *PenjualanHandler) CreatePenjualan(c *fiber.Ctx) error {
	var req models.JualHeaderRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Data input tidak valid")
	}

	errMap := make(map[string]string)

	switch {
	case req.Customer == "":
		errMap["customer"] = "Nama customer tidak boleh kosong"
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

	header := models.JualHeader{
		Customer:  req.Customer,
		UserID:    userID,
		Status:    "selesai",
		CreatedAt: time.Now(),
	}

	var details []models.JualDetail
	total := 0.0
	for _, d := range req.Details {
		if d.Qty <= 0 || d.Harga <= 0 {
			return fiber.NewError(fiber.StatusBadRequest, "qty and harga harus lebih dari 0")
		}
		subtotal := float64(d.Qty) * d.Harga
		total += subtotal
		detail := models.JualDetail{
			BarangID: d.BarangID,
			Qty:      d.Qty,
			Harga:    d.Harga,
			Subtotal: subtotal,
		}
		details = append(details, detail)
	}
	header.Total = total

	if err := h.repo.CreatePenjualan(&header, details); err != nil {
		log.Println("Error CreatePenjualan:", err.Error(), "penjualan_handler.go:CreatePenjualan", "Error at line 82")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	created, err := h.repo.GetPenjualanByID(header.ID)
	if err != nil {
		log.Println("Error fetching created penjualan:", err.Error(), "penjualan_handler.go:CreatePenjualan", "Error at line 89")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	response := mapToPenjualanResponse(created)
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetAllPenjualan godoc
// @Summary Get all sales
// @Description Get a list of all sale transactions
// @Tags Penjualan
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /penjualan [get]
func (h *PenjualanHandler) GetAllPenjualan(c *fiber.Ctx) error {
	data, err := h.repo.GetAllPenjualan()
	if err != nil {
		log.Println("Error fetching all penjualan:", err.Error(), "penjualan_handler.go:GetAllPenjualan", "Error at line 104")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	var response []models.PenjualanResponse
	for _, p := range data {
		response = append(response, mapToPenjualanResponse(&p))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": response,
	})
}

// GetPenjualanByID godoc
// @Summary Get sale by ID
// @Description Get details of a specific sale transaction
// @Tags Penjualan
// @Produce json
// @Param id path int true "Sale ID"
// @Success 200 {object} models.PenjualanResponse
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /penjualan/{id} [get]
func (h *PenjualanHandler) GetPenjualanByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "parameter ID tidak valid")
	}
	data, err := h.repo.GetPenjualanByID(uint(id))
	if err != nil {
		log.Println("Error fetching penjualan by ID:", err.Error(), "penjualan_handler.go:GetPenjualanByID", "Error at line 121")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	response := mapToPenjualanResponse(data)
	return c.Status(fiber.StatusOK).JSON(response)
}

// Private helper functions untuk mapping struct response
func mapToPenjualanResponse(p *models.JualHeader) models.PenjualanResponse {
	details := make([]models.JualDetailResponse, len(p.Details))
	for i, d := range p.Details {
		details[i] = models.JualDetailResponse{
			ID:       d.ID,
			BarangID: d.BarangID,
			Barang: models.BarangPenjualanResponse{
				KodeBarang: d.MasterBarang.KodeBarang,
				NamaBarang: d.MasterBarang.NamaBarang,
				Satuan:     d.MasterBarang.Satuan,
			},
			Qty:      d.Qty,
			Harga:    d.Harga,
			Subtotal: d.Subtotal,
		}
	}

	return models.PenjualanResponse{
		Header: models.JualHeaderResponse{
			ID:        p.ID,
			NoFaktur:  p.NoFaktur,
			Customer:  p.Customer,
			UserID:    p.UserID,
			User:      models.UserSimpleResponse{Username: p.User.Username, FullName: p.User.FullName},
			Total:     p.Total,
			Status:    p.Status,
			CreatedAt: p.CreatedAt,
		},
		Details: details,
	}
}
