package handlers

import (
	"net/http"
	"strconv"

	"github.com/therizza/freakdream/backend/internal/middleware"
	"github.com/therizza/freakdream/backend/internal/repository"
)

type FriendHandler struct {
	friends *repository.FriendRepository
}

func NewFriendHandler(friends *repository.FriendRepository) *FriendHandler {
	return &FriendHandler{friends: friends}
}

func (h *FriendHandler) SendRequest(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	targetID, err := parseIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if userID == targetID {
		respondError(w, http.StatusBadRequest, "cannot add yourself")
		return
	}

	if err := h.friends.SendRequest(r.Context(), userID, targetID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to send friend request")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "friend request sent"})
}

func (h *FriendHandler) Accept(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	requesterID, err := parseIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.friends.Accept(r.Context(), userID, requesterID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to accept friend request")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "friend request accepted"})
}

func (h *FriendHandler) Remove(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	amigoID, err := parseIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.friends.Remove(r.Context(), userID, amigoID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to remove friend")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "friend removed"})
}

func (h *FriendHandler) ListRequests(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	requests, err := h.friends.ListRequests(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list friend requests")
		return
	}

	respondJSON(w, http.StatusOK, requests)
}
