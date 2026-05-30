package handlers

import (
	"net/http"
	"strconv"

	"github.com/therizza/freakdream/backend/internal/middleware"
	"github.com/therizza/freakdream/backend/internal/repository"
)

type UserHandler struct {
	users  *repository.UserRepository
	photos *repository.PhotoRepository
}

func NewUserHandler(users *repository.UserRepository, photos *repository.PhotoRepository) *UserHandler {
	return &UserHandler{users: users, photos: photos}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	user, err := h.users.GetByID(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	photo, _ := h.photos.GetProfilePhoto(r.Context(), userID)
	foto := ""
	if photo != nil {
		foto = photo.Foto
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"id":         user.ID,
		"nome":       user.Nome,
		"sobre":      user.Sobre,
		"email":      user.Email,
		"foto":       foto,
		"created_at": user.CreatedAt,
	})
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.users.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, "user not found")
		return
	}

	photo, _ := h.photos.GetProfilePhoto(r.Context(), id)
	foto := ""
	if photo != nil {
		foto = photo.Foto
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"id":         user.ID,
		"nome":       user.Nome,
		"sobre":      user.Sobre,
		"foto":       foto,
		"created_at": user.CreatedAt,
	})
}

func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		respondJSON(w, http.StatusOK, []interface{}{})
		return
	}

	users, err := h.users.Search(r.Context(), q)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "search failed")
		return
	}

	result := make([]map[string]interface{}, 0, len(users))
	for _, u := range users {
		photo, _ := h.photos.GetProfilePhoto(r.Context(), u.ID)
		foto := ""
		if photo != nil {
			foto = photo.Foto
		}
		result = append(result, map[string]interface{}{
			"id":    u.ID,
			"nome":  u.Nome,
			"sobre": u.Sobre,
			"foto":  foto,
		})
	}

	respondJSON(w, http.StatusOK, result)
}
