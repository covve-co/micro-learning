package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"gitlab.com/covveco/special-needs/cert"
	"gitlab.com/covveco/special-needs/model"
	"gitlab.com/covveco/special-needs/view"
)

// Certificate is the handler for certificate.
func (h *Handler) Certificate(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user").(*model.User)

	courses, err := h.DB.GetCourses()
	var numCourses = len(courses)

	incomplete, err := h.userIncomplete(u)
	if err != nil {
		h.NotFoundError(w, r)
		return
	}

	b := bytes.NewBufferString("")
	view.Render(b, "certificate", struct {
		Completed  bool
		Completion string
	}{
		Completed:  len(incomplete) == 0,
		Completion: fmt.Sprintf("%v out of %v Courses completed.", numCourses-len(incomplete), numCourses),
	})
	renderMenu(w, r, "Cert", template.HTML(b.String()))
}

// CertificateDownload is the handler for downloading certificate.
func (h *Handler) CertificateDownload(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user").(*model.User)

	incomplete, err := h.userIncomplete(u)
	if err != nil {
		h.NotFoundError(w, r)
		return
	}
	if len(incomplete) != 0 {
		w.WriteHeader(http.StatusUnauthorized)
	}

	cert := cert.Certificate{
		Name: u.Name,
		Date: time.Now(),
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+u.Name+".pdf")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	cert.Generate(w)
}

func (h *Handler) userIncomplete(u *model.User) ([]model.Course, error) {
	cc, err := h.DB.GetCourses()
	if err != nil {
		return nil, err
	}

	var incomplete []model.Course

	for _, c := range cc {
		cc, err := h.DB.GetCompletionsByUserIDCourseID(u.ID, c.ID)
		if err != nil || len(cc) == 0 {
			incomplete = append(incomplete, c)
		}
	}

	return incomplete, nil
}
