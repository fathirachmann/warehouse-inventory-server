package handlers

import (
	"os"
	"regexp"
	"time"

	"warehouse-inventory-server/models"
	"warehouse-inventory-server/repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Route Handlers
func (h *UserHandler) RegisterRoute(r fiber.Router) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
}

type UserHandler struct {
	repo *repositories.UserRepository
}

func NewUserHandler(repo *repositories.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// Methods
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req *models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Data input tidak valid"})
	}

	// Register validations
	switch {
	case req.Username == "":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "username tidak boleh kosong"})
	case req.Email == "":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email tidak boleh kosong"})
	case req.Password == "":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password tidak boleh kosong"})
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format email tidak valid"})
	}

	user, err := h.repo.FindByUsername(req.Username)
	if err == nil && user != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username sudah digunakan"})
	}

	user, err = h.repo.FindByEmail(req.Email)
	if err == nil && user != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email sudah digunakan"})
	}

	pwdRegex := regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*]{8,}$`)
	if !pwdRegex.MatchString(req.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Password minimal 8 karakter, mengandung huruf dan angka dan spesial karakter"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Hash password gagal"})
	}

	// Create user
	userInput := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashed),
		FullName: req.FullName,
		Role:     "user",
	}

	if err := h.repo.Create(&userInput); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	registerResponse := models.RegisterResponse{
		Username: user.Username,
		Email:    user.Email,
		FullName: user.FullName,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "Register sukses",
		"user":   registerResponse,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req *models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": "Data input tidak valid"})
	}

	// Login validations
	switch {
	case req.Email == "":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email tidak boleh kosong"})
	case req.Password == "":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password tidak boleh kosong"})
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format email tidak valid"})
	}

	user, err := h.repo.FindByEmail(req.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau password salah"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau password salah"})
	}

	// JWT generation
	secret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "generate token gagal"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": signedToken,
		"user": fiber.Map{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	})
}
