package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/IlyaChgn/voblako/internal/models"
	fileinterfaces "github.com/IlyaChgn/voblako/internal/pkg/file"
	"github.com/IlyaChgn/voblako/internal/pkg/server/delivery/responses"
	"github.com/gorilla/mux"
)

type FileHandler struct {
	usecases   fileinterfaces.FileUsecases
	ctxUserKey string
}

func NewFileHandler(usecases fileinterfaces.FileUsecases, ctxUserKey string) *FileHandler {
	return &FileHandler{
		usecases:   usecases,
		ctxUserKey: ctxUserKey,
	}
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrBadForm)
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		log.Println(err)
		responses.SendErrResponse(w, responses.StatusInternalServerError, responses.ErrInternalServer)
		return
	}

	var contentType string
	contentType = header.Header.Get("Content-Type")
	if contentType == "" {
		buffer := make([]byte, 512)
		if _, err := file.Read(buffer); err == nil {
			contentType = http.DetectContentType(buffer)
			file.Seek(0, 0)
		}
	}

	filename := header.Filename
	if filename == "" {
		filename = "Новый файл"
	}

	user := ctx.Value(h.ctxUserKey).(*models.User)
	userID := user.ID

	metadata, err := h.usecases.UploadFile(ctx, userID, &models.GeneralFileData{
		Filename:    filename,
		ContentType: contentType,
		File:        buf.Bytes(),
		Size:        header.Size,
	})
	if err != nil {
		if errors.Is(err, models.InvalidInputError) {
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrWrongFilename)
			return
		}

		log.Println(err)
		responses.SendErrResponse(w, responses.StatusInternalServerError, responses.ErrInternalServer)
		return
	}

	responses.SendOkResponse(w, metadata)
}

func (h *FileHandler) GetFilesList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var options models.FilesListOptions
	err := json.NewDecoder(r.Body).Decode(&options)
	if err != nil {
		responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrBadJSON)
		return
	}

	user := ctx.Value(h.ctxUserKey).(*models.User)
	userID := user.ID

	files, err := h.usecases.GetFilesList(ctx, userID, options)
	if err != nil {
		if errors.Is(err, models.InvalidInputError) {
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrInvalidURLParams)
			return
		}

		log.Println(err)
		responses.SendErrResponse(w, responses.StatusInternalServerError, responses.ErrInternalServer)
		return
	}

	responses.SendOkResponse(w, files)
}

func (h *FileHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	user := ctx.Value(h.ctxUserKey).(*models.User)
	userID := user.ID

	file, err := h.usecases.GetFile(ctx, userID, id)
	if err != nil {
		switch {
		case errors.Is(err, models.InvalidInputError):
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrInvalidID)
		case errors.Is(err, models.PermissionDeniedError):
			responses.SendErrResponse(w, responses.StatusForbidden, responses.ErrForbidden)
		default:
			log.Println(err)
			responses.SendErrResponse(w, responses.StatusInternalServerError, responses.ErrInternalServer)
		}

		return
	}

	w.Header().Set("Content-Type", file.ContentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", file.Size))
	w.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=%s", file.Filename))

	if _, err := io.Copy(w, bytes.NewReader(file.File)); err != nil {
		log.Println(err)
		http.Error(w, responses.ErrInternalServer, responses.StatusInternalServerError)
		return
	}
}

func (h *FileHandler) GetMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["id"]

	user := ctx.Value(h.ctxUserKey).(*models.User)
	userID := user.ID

	metadata, err := h.usecases.GetMetadata(ctx, userID, id)
	if err != nil {
		switch {
		case errors.Is(err, models.InvalidInputError):
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrInvalidID)
		case errors.Is(err, models.PermissionDeniedError):
			responses.SendErrResponse(w, responses.StatusForbidden, responses.ErrForbidden)
		default:
			log.Println(err)
			responses.SendErrResponse(w, responses.StatusInternalServerError, responses.ErrInternalServer)
		}

		return
	}

	responses.SendOkResponse(w, metadata)
}

func (h *FileHandler) UpdateFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id := vars["id"]

	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrBadForm)
		return
	}
	defer file.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, file)
	if err != nil {
		log.Println(err)
		responses.SendErrResponse(w, responses.StatusInternalServerError, responses.ErrInternalServer)
		return
	}

	user := ctx.Value(h.ctxUserKey).(*models.User)
	userID := user.ID

	err = h.usecases.UpdateFile(ctx, userID, id, buf.Bytes(), header.Size)
	if err != nil {
		switch {
		case errors.Is(err, models.InvalidInputError):
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrInvalidID)
		case errors.Is(err, models.PermissionDeniedError):
			responses.SendErrResponse(w, responses.StatusForbidden, responses.ErrForbidden)
		default:
			log.Println(err)
			responses.SendErrResponse(w, responses.StatusInternalServerError, responses.ErrInternalServer)
		}

		return
	}

	responses.SendOkResponse(w, nil)
}

func (h *FileHandler) UpdateFilename(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id := vars["id"]

	var reqData *models.UpdateFilenameRequest
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrBadJSON)
		return
	}

	user := ctx.Value(h.ctxUserKey).(*models.User)
	userID := user.ID

	err = h.usecases.UpdateFilename(ctx, userID, id, reqData.Filename)
	if err != nil {
		switch {
		case errors.Is(err, models.InvalidFilenameError):
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrWrongFilename)
		case errors.Is(err, models.InvalidInputError):
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrInvalidID)
		case errors.Is(err, models.PermissionDeniedError):
			responses.SendErrResponse(w, responses.StatusForbidden, responses.ErrForbidden)
		default:
			log.Println(err)
			responses.SendErrResponse(w, responses.StatusInternalServerError, responses.ErrInternalServer)
		}

		return
	}

	responses.SendOkResponse(w, nil)
}

func (h *FileHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok || id == "" {
		responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrInvalidID)
		return
	}

	user := ctx.Value(h.ctxUserKey).(*models.User)
	userID := user.ID

	err := h.usecases.DeleteFile(ctx, userID, id)
	if err != nil {
		switch {
		case errors.Is(err, models.InvalidInputError):
			responses.SendErrResponse(w, responses.StatusBadRequest, responses.ErrInvalidID)
		case errors.Is(err, models.PermissionDeniedError):
			responses.SendErrResponse(w, responses.StatusForbidden, responses.ErrForbidden)
		default:
			log.Println(err)
			responses.SendErrResponse(w, responses.StatusInternalServerError, responses.ErrInternalServer)
		}

		return
	}

	responses.SendOkResponse(w, nil)
}
