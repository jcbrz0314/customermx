package handlers

import (
	"io"
	"net/http"
	"path/filepath"

	"github.com/customermx/backend/internal/domain/eventphoto"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// uniqueFilename returns a UUID-based filename preserving the original extension.
func uniqueFilename(original string) string {
	ext := filepath.Ext(original)
	return uuid.New().String() + ext
}

const maxUploadMemory = 32 << 20 // 32 MB parse buffer (handles up to 10 × 10 MB files)

// EventPhotoHandler handles event photo endpoints
type EventPhotoHandler struct {
	photoService eventphoto.Service
}

// NewEventPhotoHandler creates a new EventPhotoHandler
func NewEventPhotoHandler(photoService eventphoto.Service) *EventPhotoHandler {
	return &EventPhotoHandler{photoService: photoService}
}

// ListPhotos returns all photos for an event
// GET /api/v1/events/{eventId}/photos
func (h *EventPhotoHandler) ListPhotos(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(chi.URLParam(r, "eventId"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	photos, err := h.photoService.List(r.Context(), eventID)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, "Failed to list photos")
		return
	}

	if photos == nil {
		photos = []*eventphoto.EventPhoto{}
	}

	RespondSuccess(w, http.StatusOK, photos)
}

// UploadPhotos handles uploading one or more photos to an event
// POST /api/v1/events/{eventId}/photos  (multipart/form-data, field "files")
func (h *EventPhotoHandler) UploadPhotos(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(chi.URLParam(r, "eventId"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	if err := r.ParseMultipartForm(maxUploadMemory); err != nil {
		RespondError(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		RespondError(w, http.StatusBadRequest, "No files provided")
		return
	}

	var uploaded []*eventphoto.EventPhoto
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			RespondError(w, http.StatusInternalServerError, "Failed to read file")
			return
		}
		defer file.Close()

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		photo, err := h.photoService.Upload(r.Context(), eventID, uniqueFilename(fileHeader.Filename), contentType, fileHeader.Size, file)
		if err != nil {
			switch err {
			case eventphoto.ErrMaxPhotosReached:
				RespondError(w, http.StatusUnprocessableEntity, err.Error())
			case eventphoto.ErrInvalidContentType:
				RespondError(w, http.StatusBadRequest, err.Error())
			case eventphoto.ErrFileTooLarge:
				RespondError(w, http.StatusRequestEntityTooLarge, err.Error())
			default:
				RespondError(w, http.StatusInternalServerError, "Failed to upload photo")
			}
			return
		}

		uploaded = append(uploaded, photo)
	}

	RespondSuccess(w, http.StatusCreated, uploaded)
}

// ServePhoto proxies a photo from S3 to the client
// GET /api/v1/events/{eventId}/photos/{photoId}
func (h *EventPhotoHandler) ServePhoto(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(chi.URLParam(r, "eventId"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	photoID, err := uuid.Parse(chi.URLParam(r, "photoId"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid photo ID")
		return
	}

	body, contentType, err := h.photoService.GetForDownload(r.Context(), eventID, photoID)
	if err != nil {
		if err == eventphoto.ErrPhotoNotFound {
			RespondError(w, http.StatusNotFound, "Photo not found")
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to retrieve photo")
		}
		return
	}
	defer body.Close()

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, body) //nolint:errcheck
}

// DeletePhoto deletes a single photo
// DELETE /api/v1/events/{eventId}/photos/{photoId}
func (h *EventPhotoHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(chi.URLParam(r, "eventId"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	photoID, err := uuid.Parse(chi.URLParam(r, "photoId"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid photo ID")
		return
	}

	if err := h.photoService.Delete(r.Context(), eventID, photoID); err != nil {
		if err == eventphoto.ErrPhotoNotFound {
			RespondError(w, http.StatusNotFound, "Photo not found")
		} else {
			RespondError(w, http.StatusInternalServerError, "Failed to delete photo")
		}
		return
	}

	RespondSuccessWithMessage(w, http.StatusOK, nil, "Photo deleted successfully")
}

// ReplacePhoto replaces the file behind an existing photo record
// PUT /api/v1/events/{eventId}/photos/{photoId}  (multipart/form-data, field "file")
func (h *EventPhotoHandler) ReplacePhoto(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(chi.URLParam(r, "eventId"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid event ID")
		return
	}

	photoID, err := uuid.Parse(chi.URLParam(r, "photoId"))
	if err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid photo ID")
		return
	}

	if err := r.ParseMultipartForm(maxUploadMemory); err != nil {
		RespondError(w, http.StatusBadRequest, "Failed to parse multipart form")
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		RespondError(w, http.StatusBadRequest, "File field 'file' is required")
		return
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	photo, err := h.photoService.Replace(r.Context(), eventID, photoID, uniqueFilename(fileHeader.Filename), contentType, fileHeader.Size, file)
	if err != nil {
		switch err {
		case eventphoto.ErrPhotoNotFound:
			RespondError(w, http.StatusNotFound, "Photo not found")
		case eventphoto.ErrInvalidContentType:
			RespondError(w, http.StatusBadRequest, err.Error())
		case eventphoto.ErrFileTooLarge:
			RespondError(w, http.StatusRequestEntityTooLarge, err.Error())
		default:
			RespondError(w, http.StatusInternalServerError, "Failed to replace photo")
		}
		return
	}

	RespondSuccess(w, http.StatusOK, photo)
}
