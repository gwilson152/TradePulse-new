package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/tradepulse/api/internal/database"
	"github.com/tradepulse/api/internal/middleware"
	"github.com/tradepulse/api/internal/models"
)

type TagsHandler struct {
	db *database.DB
}

func NewTagsHandler(db *database.DB) *TagsHandler {
	return &TagsHandler{db: db}
}

// ListTags handles GET /api/tags
func (h *TagsHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	tags, err := h.db.ListTags(r.Context(), userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch tags", err)
		return
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    tags,
	})
}

// GetTag handles GET /api/tags/{id}
func (h *TagsHandler) GetTag(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	tagID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid tag ID", err)
		return
	}

	tag, err := h.db.GetTag(r.Context(), tagID, userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to fetch tag", err)
		return
	}
	if tag == nil {
		sendError(w, http.StatusNotFound, "Tag not found", nil)
		return
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    tag,
	})
}

// CreateTag handles POST /api/tags
func (h *TagsHandler) CreateTag(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)

	var tag models.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Set user ID
	tag.UserID = userID

	// Validate required fields
	if tag.Name == "" {
		sendError(w, http.StatusBadRequest, "Tag name is required", nil)
		return
	}

	// Create tag
	if err := h.db.CreateTag(r.Context(), &tag); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to create tag", err)
		return
	}

	sendJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data":    tag,
	})
}

// UpdateTag handles PUT /api/tags/{id}
func (h *TagsHandler) UpdateTag(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	tagID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid tag ID", err)
		return
	}

	var tag models.Tag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate required fields
	if tag.Name == "" {
		sendError(w, http.StatusBadRequest, "Tag name is required", nil)
		return
	}

	// Update tag
	if err := h.db.UpdateTag(r.Context(), tagID, userID, &tag); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to update tag", err)
		return
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    tag,
	})
}

// DeleteTag handles DELETE /api/tags/{id}
func (h *TagsHandler) DeleteTag(w http.ResponseWriter, r *http.Request) {
	userID, _ := middleware.GetUserID(r)
	tagID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		sendError(w, http.StatusBadRequest, "Invalid tag ID", err)
		return
	}

	if err := h.db.DeleteTag(r.Context(), tagID, userID); err != nil {
		sendError(w, http.StatusInternalServerError, "Failed to delete tag", err)
		return
	}

	sendJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Tag deleted successfully",
	})
}
