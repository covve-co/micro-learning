package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"gitlab.com/covveco/special-needs/model"
	"gitlab.com/covveco/special-needs/view"
)

// About is the handler for about page.
func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	user, _ := r.Context().Value("user").(*model.User)
	b := bytes.NewBufferString("")

	courses, err := h.DB.GetCourses()
	var numCourses = len(courses)

	incomplete, err := h.userIncomplete(user)
	if err != nil {
		h.NotFoundError(w, r)
		return
	}

	view.Render(b, "about", struct {
		Username   string
		Completion string
	}{
		Username:   user.Name,
		Completion: fmt.Sprintf("%v out of %v Courses completed.", numCourses-len(incomplete), numCourses),
	})

	renderMenu(w, r, "About", template.HTML(b.String()))
}
