package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/therizza/freakdream/backend/internal/middleware"
	"github.com/therizza/freakdream/backend/internal/models"
	"github.com/therizza/freakdream/backend/internal/repository"
	"github.com/therizza/freakdream/backend/internal/utils"
)

type PostHandler struct {
	posts     *repository.PostRepository
	photos    *repository.PhotoRepository
	users     *repository.UserRepository
	uploadDir string
}

func NewPostHandler(
	posts *repository.PostRepository,
	photos *repository.PhotoRepository,
	users *repository.UserRepository,
	uploadDir string,
) *PostHandler {
	return &PostHandler{posts: posts, photos: photos, users: users, uploadDir: uploadDir}
}

func (h *PostHandler) Feed(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	limit := 20
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if o, err := strconv.Atoi(v); err == nil && o >= 0 {
			offset = o
		}
	}

	posts, err := h.posts.Feed(r.Context(), userID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to load feed")
		return
	}

	h.enrichPosts(r, posts)
	respondJSON(w, http.StatusOK, posts)
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	var req struct {
		Texto string `json:"texto"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Texto == "" {
		respondError(w, http.StatusBadRequest, "texto is required")
		return
	}

	if len(req.Texto) > 2000 {
		respondError(w, http.StatusBadRequest, "texto must be at most 2000 characters")
		return
	}

	id, err := h.posts.Create(r.Context(), userID, models.PostTypeText, req.Texto)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create post")
		return
	}

	post, _ := h.posts.GetByID(r.Context(), id)
	respondJSON(w, http.StatusCreated, post)
}

func (h *PostHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		respondError(w, http.StatusBadRequest, "failed to parse form")
		return
	}

	file, header, err := r.FormFile("imagem")
	if err != nil {
		respondError(w, http.StatusBadRequest, "imagem is required")
		return
	}
	defer file.Close()

	filename, err := utils.SaveUpload(h.uploadDir, file, header)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	texto := r.FormValue("texto")
	fotoID, err := h.photos.Create(r.Context(), filename, texto)
	if err != nil {
		_ = utils.DeleteUpload(h.uploadDir, filename)
		respondError(w, http.StatusInternalServerError, "failed to save photo")
		return
	}

	usePerfil := r.FormValue("perfil") == "1"
	if usePerfil {
		_ = h.photos.SetProfilePhoto(r.Context(), userID, fotoID)
	}

	postID, err := h.posts.Create(r.Context(), userID, models.PostTypePhoto, strconv.FormatInt(fotoID, 10))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create photo post")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"post_id": postID,
		"foto_id": fotoID,
		"foto":    filename,
	})
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	postID, err := parseIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	post, err := h.posts.GetByID(r.Context(), postID)
	if err != nil {
		respondError(w, http.StatusNotFound, "post not found")
		return
	}

	// If photo post, delete the photo file too
	if post.Tipo == models.PostTypePhoto {
		fotoID, _ := strconv.ParseInt(post.Texto, 10, 64)
		photo, photoErr := h.photos.GetByID(r.Context(), fotoID)
		if photoErr == nil && photo != nil {
			_ = utils.DeleteUpload(h.uploadDir, photo.Foto)
			_ = h.photos.DeleteProfilePhoto(r.Context(), userID, fotoID)
			_ = h.photos.Delete(r.Context(), fotoID)
		}
	}

	if err := h.posts.Delete(r.Context(), postID, userID); err != nil {
		respondError(w, http.StatusForbidden, "cannot delete this post")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func (h *PostHandler) Share(w http.ResponseWriter, r *http.Request) {
	userIDStr := middleware.GetUserID(r.Context())
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	var req struct {
		PostID int64 `json:"post_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PostID <= 0 {
		respondError(w, http.StatusBadRequest, "post_id is required")
		return
	}

	if _, err := h.posts.GetByID(r.Context(), req.PostID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			respondError(w, http.StatusNotFound, "source post not found")
		} else {
			respondError(w, http.StatusInternalServerError, "failed to load source post")
		}
		return
	}

	id, err := h.posts.Create(r.Context(), userID, models.PostTypeShare, strconv.FormatInt(req.PostID, 10))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to share post")
		return
	}

	post, _ := h.posts.GetByID(r.Context(), id)
	respondJSON(w, http.StatusCreated, post)
}

func (h *PostHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r, "id")
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid user id")
		return
	}

	limit := 20
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if l, err := strconv.Atoi(v); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if o, err := strconv.Atoi(v); err == nil && o >= 0 {
			offset = o
		}
	}

	posts, err := h.posts.GetByUser(r.Context(), id, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to load posts")
		return
	}

	h.enrichPosts(r, posts)
	respondJSON(w, http.StatusOK, posts)
}

func (h *PostHandler) enrichPosts(r *http.Request, posts []models.Post) {
	for i := range posts {
		user, err := h.users.GetByID(r.Context(), posts[i].UserID)
		if err != nil {
			continue
		}
		photo, _ := h.photos.GetProfilePhoto(r.Context(), posts[i].UserID)
		foto := ""
		if photo != nil {
			foto = photo.Foto
		}
		posts[i].User = &models.UserPublic{
			ID:    user.ID,
			Nome:  user.Nome,
			Sobre: user.Sobre,
			Foto:  foto,
		}
	}
}
