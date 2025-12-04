package handlers

import (
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

// RegisterRoute mendaftarkan seluruh endpoint pembelian sesuai main.go
func (h *PembelianHandler) RegisterRoute(r fiber.Router) {
	r.Post("/", h.CreatePembelian)
	r.Get("/", h.GetAllPembelian)
	r.Get("/:id", h.GetPembelianByID)
}

// CreatePembelian meng-handle pembuatan pembelian baru beserta update stok dan history
func (h *PembelianHandler) CreatePembelian(c *fiber.Ctx) error {
	var req models.BeliHeaderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Data input tidak valid"})
	}

	switch {
	case req.NoFaktur == "":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no_faktur tidak boleh kosong"})
	case req.Supplier == "":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "supplier tidak boleh kosong"})
	case len(req.Details) == 0:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "details tidak boleh kosong"})
	}

	var userID uint
	if claims, ok := c.Locals("user").(jwt.MapClaims); ok {
		if sub, ok := claims["sub"]; ok {
			switch v := sub.(type) {
			case float64:
				userID = uint(v)
			case int:
				userID = uint(v)
			}
		}
	}

	header := models.BeliHeader{
		NoFaktur:  req.NoFaktur,
		Supplier:  req.Supplier,
		UserID:    userID,
		Status:    "selesai",
		CreatedAt: time.Now(),
	}

	var details []models.BeliDetail
	total := 0.0
	for _, d := range req.Details {
		if d.Qty <= 0 || d.Harga <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "qty dan harga tidak boleh kurang dari sama dengan 0"})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
			"error":  err.Error(),
		})
	}

	// Update stok & history untuk setiap detail (stok masuk)
	for _, d := range details {
		stok, err := h.stokRepo.GetOrCreateByBarangID(d.BarangID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "Server error",
				"error":  err.Error(),
			})
		}
		stokSebelum := stok.StokAkhir
		stokSesudah := stokSebelum + d.Qty
		stok.StokAkhir = stokSesudah
		if err := h.stokRepo.UpdateStok(stok); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "Server error",
				"error":  err.Error(),
			})
		}

		history := models.HistoryStok{
			BarangID:       d.BarangID,
			UserID:         userID,
			JenisTransaksi: "masuk",
			Jumlah:         d.Qty,
			StokSebelumnya: stokSebelum,
			StokSesudah:    stokSesudah,
			Keterangan:     "Pembelian " + header.NoFaktur,
		}
		if err := h.stokRepo.CreateHistory(&history); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "Server error",
				"error":  err.Error(),
			})
		}
	}

	// Reload header + details untuk response
	created, err := h.repo.GetPembelianByID(header.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
			"error":  err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

// GetAllPembelian adalah method untuk pengambilan semua data pembelian
func (h *PembelianHandler) GetAllPembelian(c *fiber.Ctx) error {
	data, err := h.repo.GetAllPembelian()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
			"error":  err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(data)
}

// GetPembelianByID adalah method untuk pengambilan data pembelian berdasarkan ID
func (h *PembelianHandler) GetPembelianByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Data input tidak valid"})
	}
	data, err := h.repo.GetPembelianByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "Server error",
			"error":  err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(data)
}
