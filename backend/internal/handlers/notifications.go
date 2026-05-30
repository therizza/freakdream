package handlers

import (
	"net/http"
	"strconv"

	"github.com/therizza/freakdream/backend/internal/middleware"
	"github.com/therizza/freakdream/backend/internal/repository"
)

type NotificationHandler struct {
	notifs *repository.NotificationRepository
}

func NewNotificationHandler(notifs *repository.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{notifs: notifs}
}

func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	notifs, err := h.notifs.List(r.Context(), userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list notifications")
		return
	}

	respondJSON(w, http.StatusOK, notifs)
}

func (h *NotificationHandler) Remove(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	notifID, err := parseIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid notification id")
		return
	}

	if err := h.notifs.Remove(r.Context(), notifID, userID); err != nil {
		respondError(w, http.StatusNotFound, "notification not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "notification removed"})
}

func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	notifID, err := parseIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid notification id")
		return
	}

	if err := h.notifs.MarkRead(r.Context(), notifID, userID); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to mark as read")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "marked as read"})
}
