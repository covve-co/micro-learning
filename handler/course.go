package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"gitlab.com/covveco/special-needs/model"
	"gitlab.com/covveco/special-needs/view"
)

// Course is the handler for course.
func (h *Handler) Course(w http.ResponseWriter, r *http.Request) {
	cc, err := h.DB.GetCourses()
	if err != nil {
		fmt.Println(err)
		h.NotFoundError(w, r)
		return
	}

	b := bytes.NewBufferString("")
	view.Render(b, "course", struct {
		Courses []model.Course
	}{
		Courses: cc,
	})

	renderMenu(w, r, "Courses", template.HTML(b.String()))
}

// CourseID is the handler for module with given ID.
func (h *Handler) CourseID(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.NotFoundError(w, r)
		return
	}

	c, err := h.DB.GetCourseByTemplateName(name)
	if err != nil {
		h.NotFoundError(w, r)
		return
	}

	if id < 1 || id > c.NumSections {
		h.NotFoundError(w, r)
		return
	}

	next := fmt.Sprintf("/module/%s/%d", name, id+1)
	if id+1 > c.NumSections {
		next = fmt.Sprintf("/quiz/%d", c.ID)
	}

	b := bytes.NewBufferString("")
	view.Render(b, fmt.Sprintf("courses/%s/%d", name, id), struct {
		Next string
	}{
		Next: next,
	})

	renderMenu(w, r, fmt.Sprintf("%s - %d", c.Title, id), template.HTML(b.String()))
}
