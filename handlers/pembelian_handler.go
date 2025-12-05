package handlers

import (
	"log"
	"time"

	"warehouse-inventory-server/models"
	"warehouse-inventory-server/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type PembelianHandler struct {
	repo     *repositories.PembelianRepository
	stokRepo *repositories.StokRepository
}

func NewPembelianHandler(repo *repositories.PembelianRepository, stokRepo *repositories.StokRepository) *PembelianHandler {
	return &PembelianHandler{
		repo:     repo,
		stokRepo: stokRepo,
	}
}

// RegisterRoute mendaftarkan seluruh endpoint "/api/pembelian"
func (h *PembelianHandler) RegisterRoute(r fiber.Router) {
	r.Post("/", h.CreatePembelian)
	r.Get("/", h.GetAllPembelian)
	r.Get("/:id", h.GetPembelianByID)
}

// CreatePembelian handle pembuatan pembelian baru beserta update stok dan history
func (h *PembelianHandler) CreatePembelian(c *fiber.Ctx) error {
	var req models.BeliHeaderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error":   "body parsing error",
			"message": "Data input tidak valid",
		})
	}

	errMap := make(map[string][]string)

	switch {
	case req.Supplier == "":
		errMap["supplier"] = append(errMap["supplier"], "Nama supplier tidak boleh kosong")
	case len(req.Details) == 0:
		errMap["details"] = append(errMap["details"], "details tidak boleh kosong")
	}

	if len(errMap) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors":  "validation error",
			"message": errMap,
		})
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
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "validation error",
				"message": "qty dan harga tidak boleh kurang dari sama dengan 0",
			})
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
		log.Println("Error CreatePembelian:", err.Error(), "pembelian_handler.go:CreatePembelian", "Error at line 84")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
		})
	}

	created, err := h.repo.GetPembelianByID(header.ID)
	if err != nil {
		log.Println("Error fetching created pembelian:", err.Error(), "pembelian_handler.go:CreatePembelian", "Error at line 91")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// GetAllPembelian adalah method untuk pengambilan semua data pembelian
func (h *PembelianHandler) GetAllPembelian(c *fiber.Ctx) error {
	data, err := h.repo.GetAllPembelian()
	if err != nil {
		log.Println("Error fetching all pembelian:", err.Error(), "pembelian_handler.go:GetAllPembelian", "Error at line 104")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

// GetPembelianByID adalah method untuk pengambilan data pembelian berdasarkan ID
func (h *PembelianHandler) GetPembelianByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"message": "Data input tidak valid"})
	}
	data, err := h.repo.GetPembelianByID(uint(id))
	if err != nil {
		log.Println("Error fetching pembelian by ID:", err.Error(), "pembelian_handler.go:GetPembelianByID", "Error at line 122")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"header":  data,
		"details": data.Details,
	})
}
