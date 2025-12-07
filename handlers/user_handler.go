package handlers

import (
	"log"
	"os"
	"regexp"
	"time"

	"warehouse-inventory-server/middleware"
	"warehouse-inventory-server/models"
	"warehouse-inventory-server/repositories"
	"warehouse-inventory-server/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Route Handlers
func (h *UserHandler) RegisterRoute(r fiber.Router) {
	r.Post("/register", middleware.Authentication(), middleware.GuardAdmin(), h.Register) // Simple Authorization: Only admin can register new staff
	r.Post("/login", h.Login)
}

type UserHandler struct {
	repo *repositories.UserRepository
}

func NewUserHandler(repo *repositories.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// Register godoc
// @Summary Register new user (Admin only)
// @Description Register a new staff user. Only admin can perform this action.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.RegisterRequest true "Register Request"
// @Success 201 {object} models.RegisterResponse "Created"
// @Failure 400 {object} middleware.SpecificErrorResponse "Bad Request"
// @Failure 422 {object} middleware.ErrorResponse "Unprocessable Entity"
// @Failure 500 {object} middleware.ErrorResponse "Internal Server Error"
// @Security BearerAuth
// @Router /api/auth/register [post]
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req *models.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("Error parsing registration body:", err.Error(), "user_handler.go:Register", "Error at line 48")
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Data input tidak valid")
	}

	// Register validations
	errMap := make(map[string]string)

	// Username validation
	if req.Username == "" {
		errMap["username"] = "username tidak boleh kosong"
	} else if len(req.Username) < 4 {
		errMap["username"] = "username minimal 4 karakter"
	} else {
		user, err := h.repo.FindByUsername(req.Username)
		if err == nil && user != nil {
			errMap["username"] = "Username sudah digunakan"
		}
	}

	// FullName validation
	if req.FullName == "" {
		errMap["full_name"] = "full name tidak boleh kosong"
	}

	// Email validation
	if req.Email == "" {
		errMap["email"] = "email tidak boleh kosong"
	} else {
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(req.Email) {
			errMap["email"] = "Format email tidak valid"
		} else {
			user, err := h.repo.FindByEmail(req.Email)
			if err == nil && user != nil {
				errMap["email"] = "Email sudah digunakan"
			}
		}
	}

	// Password validation
	if req.Password == "" {
		errMap["password"] = "password tidak boleh kosong"
	} else if !utils.ValidatePassword(req.Password) {
		errMap["password"] = "Password minimal 8 karakter, mengandung huruf besar, huruf kecil, angka, dan simbol"
	}

	if len(errMap) > 0 {
		return &middleware.ValidationError{
			Message: "validation error",
			Errors:  errMap,
		}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		log.Println("Error hashing password during registration:", err.Error(), "user_handler.go:Register", "Error at line 102")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	// Create user
	userInput := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashed),
		FullName: req.FullName,
		Role:     "staff",
	}

	if err := h.repo.Create(&userInput); err != nil {
		log.Println("Error creating user during registration:", err.Error(), "user_handler.go:Register", "Error at line 117")
		return fiber.NewError(fiber.StatusBadRequest, "database error")
	}

	registerResponse := models.RegisterResponse{
		Username: userInput.Username,
		Email:    userInput.Email,
		FullName: userInput.FullName,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "Register sukses",
		"user":   registerResponse,
	})
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "Login Request"
// @Success 200 {object} models.LoginResponse "OK"
// @Failure 400 {object} middleware.SpecificErrorResponse "Bad Request"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 422 {object} middleware.ErrorResponse "Unprocessable Entity"
// @Failure 500 {object} middleware.ErrorResponse "Internal Server Error"
// @Router /api/auth/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req *models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("Error parsing login body:", err.Error(), "user_handler.go:Login", "Error at line 149")
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Data input tidak valid")
	}

	// Login validations
	errMap := make(map[string]string)

	if req.Email == "" {
		errMap["email"] = "email tidak boleh kosong"
	}

	if req.Password == "" {
		errMap["password"] = "password tidak boleh kosong"
	}

	if len(errMap) > 0 {
		return &middleware.ValidationError{
			Message: "validation error",
			Errors:  errMap,
		}
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		errMap["email"] = "Format email tidak valid"
	}

	if len(errMap) > 0 {
		return &middleware.ValidationError{
			Message: "validation error",
			Errors:  errMap,
		}
	}

	// Authenticate user
	user, err := h.repo.FindByEmail(req.Email)
	if err != nil || user == nil {
		// Jangan lanjut ke pengecekan password jika user tidak ditemukan
		// Agar tidak terjadi panic nil pointer dan pesan tetap generic
		return &middleware.ValidationError{
			Message: "validation error",
			Errors: map[string]string{
				"email":    "Email atau password salah",
				"password": "Email atau password salah",
			},
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return &middleware.ValidationError{
			Message: "validation error",
			Errors: map[string]string{
				"email":    "Email atau password salah",
				"password": "Email atau password salah",
			},
		}
	}

	if len(errMap) > 0 {
		return &middleware.ValidationError{
			Message: "validation error",
			Errors:  errMap,
		}
	}

	// JWT generation
	secret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Println("Error signing JWT token during login:", err.Error(), "user_handler.go:Login", "Error at line 190")
		return fiber.NewError(fiber.StatusInternalServerError, "Server error")
	}

	loginResponse := models.LoginResponse{
		Token: signedToken,
	}

	return c.Status(fiber.StatusOK).JSON(loginResponse)
}
