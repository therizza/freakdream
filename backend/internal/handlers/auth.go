package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/therizza/freakdream/backend/internal/models"
	"github.com/therizza/freakdream/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	users     *repository.UserRepository
	jwtSecret string
	jwtExpiry time.Duration
	limiter   http.Handler
}

func NewAuthHandler(users *repository.UserRepository, jwtSecret string, jwtExpiry time.Duration) *AuthHandler {
	return &AuthHandler{users: users, jwtSecret: jwtSecret, jwtExpiry: jwtExpiry}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Senha == "" {
		respondError(w, http.StatusBadRequest, "email and senha are required")
		return
	}

	user, err := h.users.GetByEmail(r.Context(), req.Email)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid email or senha")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Senha), []byte(req.Senha)); err != nil {
		respondError(w, http.StatusUnauthorized, "invalid email or senha")
		return
	}

	ip := r.RemoteAddr
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		ip = fwd
	}
	_ = h.users.RecordLogin(r.Context(), user.ID, ip)

	tokenStr, err := h.generateToken(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"token": tokenStr,
		"user": models.UserPublic{
			ID:    user.ID,
			Nome:  user.Nome,
			Sobre: user.Sobre,
		},
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Senha == "" || req.Nome == "" {
		respondError(w, http.StatusBadRequest, "nome, email, and senha are required")
		return
	}

	if len(req.Senha) < 8 {
		respondError(w, http.StatusBadRequest, "senha must be at least 8 characters")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Senha), bcrypt.DefaultCost)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	id, err := h.users.Create(r.Context(), req.Nome, req.Sobre, req.Email, string(hash))
	if err != nil {
		respondError(w, http.StatusConflict, "email already registered")
		return
	}

	tokenStr, err := h.generateToken(id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"token": tokenStr,
		"user": models.UserPublic{
			ID:    id,
			Nome:  req.Nome,
			Sobre: req.Sobre,
		},
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Stateless JWT: logout is handled client-side by discarding the token.
	respondJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h *AuthHandler) generateToken(userID int64) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", userID),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(h.jwtExpiry)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}

// helpers shared by all handlers

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}

func parseIDParam(r *http.Request, paramName string) (int64, error) {
	v := r.PathValue(paramName)
	return strconv.ParseInt(v, 10, 64)
}
