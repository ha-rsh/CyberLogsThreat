package handlers

import (
	"cybersecuritySystem/shared/auth"
	"cybersecuritySystem/shared/constants"
	"cybersecuritySystem/shared/logger"
	"cybersecuritySystem/shared/utils"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	users map[string]User
}

type User struct {
	ID       		string 		`json:"id"`
	Username 		string 		`json:"username"`
	Password 		string 		`json:"password"`
	Role     		string 		`json:"role"`
}

type LoginRequest struct {
	Username 		string 		`json:"username"`
	Password 		string 		`json:"password"`
}

type RegisterRequest struct {
	Username 		string 		`json:"username"`
	Password 		string 		`json:"password"`
	Role     		string 		`json:"role"`
}

type AuthResponse struct {
	Token    		string 		`json:"token"`
	Username 		string 		`json:"username"`
	Role     		string 		`json:"role"`
}

func NewAuthHandler() *AuthHandler {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(constants.UserPassword), bcrypt.DefaultCost)
	
	return &AuthHandler{
		users: map[string]User{
			"admin": {
				ID:       "1",
				Username: constants.UserName,
				Password: string(hashedPassword),
				Role:     "admin",
			},
		},
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} utils.APIResponse{data=AuthResponse}
// @Failure 401 {object} utils.APIResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	user, exists := h.users[req.Username]
	if !exists {
		utils.SendErrorResponse(w, "Invalid credentials", http.StatusUnauthorized, "User not found")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.SendErrorResponse(w, "Invalid credentials", http.StatusUnauthorized, "Incorrect password")
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		logger.Error("Failed to generate token: %v", err)
		utils.SendErrorResponse(w, "Failed to generate token", http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("User logged in: %s", user.Username)

	utils.SendSuccessResponse(w, AuthResponse{
		Token:    token,
		Username: user.Username,
		Role:     user.Role,
	}, nil)
}

// Register godoc
// @Summary Register new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration data"
// @Success 201 {object} utils.APIResponse{data=AuthResponse}
// @Failure 400 {object} utils.APIResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendErrorResponse(w, "Invalid request body", http.StatusBadRequest, err.Error())
		return
	}

	if _, exists := h.users[req.Username]; exists {
		utils.SendErrorResponse(w, "User already exists", http.StatusBadRequest, "Username taken")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Failed to hash password: %v", err)
		utils.SendErrorResponse(w, "Failed to create user", http.StatusInternalServerError, err.Error())
		return
	}

	if req.Role == "" {
		req.Role = "user"
	}

	userID := generateUserID()
	user := User{
		ID:       userID,
		Username: req.Username,
		Password: string(hashedPassword),
		Role:     req.Role,
	}

	h.users[req.Username] = user

	token, err := auth.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		logger.Error("Failed to generate token: %v", err)
		utils.SendErrorResponse(w, "Failed to generate token", http.StatusInternalServerError, err.Error())
		return
	}

	logger.Info("New user registered: %s", user.Username)

	w.WriteHeader(http.StatusCreated)
	utils.SendSuccessResponse(w, AuthResponse{
		Token:    token,
		Username: user.Username,
		Role:     user.Role,
	}, nil)
}


func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		utils.SendErrorResponse(w, "Missing authorization header", http.StatusUnauthorized, "")
		return
	}

	tokenString := authHeader[7:]

	newToken, err := auth.RefreshToken(tokenString)
	if err != nil {
		utils.SendErrorResponse(w, "Failed to refresh token", http.StatusUnauthorized, err.Error())
		return
	}

	username := r.Header.Get("X-Username")
	role := r.Header.Get("X-User-Role")

	utils.SendSuccessResponse(w, AuthResponse{
		Token:    newToken,
		Username: username,
		Role:     role,
	}, nil)
}

func generateUserID() string {
	return uuid.New().String()
}