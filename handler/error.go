package handler

import (
	"net/http"

	"github.com/covveco/micro-learning/view"
)

func (h *Handler) NotFoundError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	view.Render(w, "error", struct {
		Code    int
		Message string
		Title   string
	}{
		Code:    http.StatusNotFound,
		Message: "Not Found",
		Title:   "Page Not Found",
	})
}

func (h *Handler) InternalServerError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

	view.Render(w, "error", struct {
		Code    int
		Message string
		Title   string
	}{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
		Title:   "Interal Server Error",
	})
}
