package handlers

import (
	"log"
	"time"

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

func (h *PenjualanHandler) CreatePenjualan(c *fiber.Ctx) error {
	var req models.JualHeaderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Data input tidak valid"})
	}

	switch {
	case req.Customer == "":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "customer tidak boleh kosong"})
	case len(req.Details) == 0:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "details tidak boleh kosong"})
	}

	var userID uint
	if claims, ok := c.Locals("user").(jwt.MapClaims); ok {
		if sub, ok := claims["od"]; ok {
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "qty and harga must be positive"})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
		})
	}

	created, err := h.repo.GetPenjualanByID(header.ID)
	if err != nil {
		log.Println("Error fetching created penjualan:", err.Error(), "penjualan_handler.go:CreatePenjualan", "Error at line 89")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *PenjualanHandler) GetAllPenjualan(c *fiber.Ctx) error {
	data, err := h.repo.GetAllPenjualan()
	if err != nil {
		log.Println("Error fetching all penjualan:", err.Error(), "penjualan_handler.go:GetAllPenjualan", "Error at line 104")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

func (h *PenjualanHandler) GetPenjualanByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	data, err := h.repo.GetPenjualanByID(uint(id))
	if err != nil {
		log.Println("Error fetching penjualan by ID:", err.Error(), "penjualan_handler.go:GetPenjualanByID", "Error at line 121")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"header":  data,
		"details": data.Details,
	})
}
