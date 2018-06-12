package handler

import (
	"context"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func (h *Handler) PrepareLog(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewV4()

		log := h.Log.WithFields(logrus.Fields{
			"request_id": id,
		})

		ctx := context.WithValue(r.Context(), "log", log)

		n.ServeHTTP(w, r.WithContext(ctx))
	})
}

type responseWriter struct {
	statusCode int
	w          http.ResponseWriter
}

func (w *responseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *responseWriter) Write(b []byte) (int, error) {
	return w.w.Write(b)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.w.WriteHeader(statusCode)
}

func (h *Handler) LogRequest(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := r.Context().Value("log").(logrus.FieldLogger)

		// Before
		log.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL,
		}).Infoln("received request")

		ww := &responseWriter{http.StatusOK, w}
		n.ServeHTTP(ww, r)

		// After
		log.WithFields(logrus.Fields{
			"status_code": ww.statusCode,
		}).Infoln("responded to request")
	})
}
