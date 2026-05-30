package handlers

import (
	"net/http"
	"strconv"

	"github.com/therizza/freakdream/backend/internal/middleware"
	"github.com/therizza/freakdream/backend/internal/repository"
)

type FollowHandler struct {
	follows *repository.FollowRepository
}

func NewFollowHandler(follows *repository.FollowRepository) *FollowHandler {
	return &FollowHandler{follows: follows}
}

func (h *FollowHandler) Follow(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	targetID, err := parseIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if userID == targetID {
		respondError(w, http.StatusBadRequest, "cannot follow yourself")
		return
	}

	if err := h.follows.Follow(r.Context(), userID, targetID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to follow user")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "followed"})
}

func (h *FollowHandler) Unfollow(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	targetID, err := parseIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.follows.Unfollow(r.Context(), userID, targetID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to unfollow user")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "unfollowed"})
}

func (h *FollowHandler) Count(w http.ResponseWriter, r *http.Request) {
	targetID, err := parseIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	count, err := h.follows.Count(r.Context(), targetID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to count followers")
		return
	}

	respondJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (h *FollowHandler) IsFollowing(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	targetID, err := parseIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	ok, err := h.follows.IsFollowing(r.Context(), userID, targetID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to check follow status")
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{"following": ok})
}
