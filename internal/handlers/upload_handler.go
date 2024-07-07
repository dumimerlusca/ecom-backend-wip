package handlers

import (
	"ecom-backend/internal/jsonlog"
	"ecom-backend/internal/model"
	"ecom-backend/internal/service"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UploadHandler struct {
	BaseHandler
	uploadSvc *service.UploadService
}

func NewUploadHandler(logger *jsonlog.Logger, uploadSvc *service.UploadService) *UploadHandler {
	return &UploadHandler{BaseHandler: BaseHandler{logger: logger}, uploadSvc: uploadSvc}
}

func (h *UploadHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.ServerErrorResponse(w, r, err)
		return
	}

	files := r.MultipartForm.File["files"]

	type FileItem struct {
		Id string `json:"id"`
	}

	fileItems := make([]FileItem, 0)

	for _, fileHeader := range files {
		fileRecord, err := h.uploadSvc.UploadFile(r.Context(), fileHeader)

		if err != nil {
			h.ServerErrorResponse(w, r, err)
			return
		}

		fileItems = append(fileItems, FileItem{Id: fileRecord.Id})

	}

	h.WriteJson(w, http.StatusCreated, ResponseBody{Payload: Envelope{"files": fileItems}}, nil)
}

func (h *UploadHandler) ServerFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileId := ps.ByName("fileId")

	if fileId == "" {
		h.NotFoundResponse(w, r)
		return
	}

	fileRecord, err := h.uploadSvc.GetFileRecord(fileId)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			h.NotFoundResponse(w, r)
		default:
			h.ServerErrorResponse(w, r, err)
		}
		return
	}

	http.ServeFile(w, r, fmt.Sprintf("uploads/%s%s", fileRecord.Id, fileRecord.Extension))
}
