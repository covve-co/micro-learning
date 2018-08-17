package handler

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/covveco/micro-learning/model"
	"github.com/covveco/micro-learning/view"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

// QuizID is the handler for quiz with the given ID.
func (h *Handler) QuizID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	courses, err := h.DB.GetCourses()
	var numCourses = len(courses)

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 || id > numCourses {
		h.NotFoundError(w, r)
		return
	}

	qq, err := h.DB.GetQuestionsByCourseID(id)
	if err != nil {
		// FIXME better handling
		h.NotFoundError(w, r)
		return
	}

	for i := range qq {
		q := qq[i]

		oo, err := h.DB.GetOptionsByQuestionID(q.ID)
		if err != nil {
			// FIXME better handling
			h.NotFoundError(w, r)
			return
		}
		q.Options = oo

		q.Correct = true

		qq[i] = q
	}

	b := bytes.NewBufferString("")
	view.Render(b, "quiz", struct {
		Completed     bool
		QuestionCount int
		Questions     []model.Question
	}{
		Completed:     false,
		QuestionCount: len(qq),
		Questions:     qq,
	})

	renderMenu(w, r, "Quiz", template.HTML(b.String()))
}

// QuizSubmit handles the quiz response
func (h *Handler) QuizSubmit(w http.ResponseWriter, r *http.Request) {
	log := r.Context().Value("log").(logrus.FieldLogger)
	idStr := chi.URLParam(r, "id")

	courses, err := h.DB.GetCourses()
	var numCourses = len(courses)

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 || id > numCourses {
		log.Warnf("no of courses invalid or out of range")
		h.NotFoundError(w, r)
		return
	}

	if err = r.ParseForm(); err != nil {
		log.Warnf("parsing form values failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	qq, err := h.DB.GetQuestionsByCourseID(id)
	if err != nil {
		log.Errorf("getting questions from db failed")
		h.NotFoundError(w, r)
		return
	}

	for i := range qq {
		qq[i].Options, err = h.DB.GetOptionsByQuestionID(qq[i].ID)
		if err != nil {
			log.Errorf("getting options for question %d from db failed", qq[i].ID)
			h.NotFoundError(w, r)
			return
		}

		answer, err := strconv.Atoi(r.FormValue(strconv.Itoa(qq[i].ID)))
		if err != nil {
			// FIXME: Find proper error
			log.Errorf("parsing form value failed")
			h.NotFoundError(w, r)
			return
		}

		for _, o := range qq[i].Options {
			if o.Position == answer {
				qq[i].Correct = o.Correct
			}
		}
	}

	var wrong bool
	for a := range qq {
		if !qq[a].Correct {
			wrong = true
			break
		}
	}

	params := struct {
		Completed     bool
		QuestionCount int
		Questions     []model.Question
		Next          string
	}{}

	if wrong {
		params.QuestionCount = len(qq)
		params.Questions = qq
		log.Warnf("user answer was wrong")
	} else {
		u := r.Context().Value("user").(*model.User)
		c := model.Completion{
			UserID:   u.ID,
			CourseID: id,
		}
		err = h.DB.CreateCompletion(&c)
		if err != nil {
			log.Errorf("creating completion failed")
			// FIXME: Check if user has finished this quiz before. Redirect to next quiz if successful
			h.NotFoundError(w, r)
			return
		}

		log.Infof("successfully created completion")

		// FIXME: Inefficient DB calls
		css, err := h.DB.GetCourses()
		if err != nil {
			h.NotFoundError(w, r)
		}
		if id == len(css) {
			redirect(w, r, "/certificate", "")
			return
		}

		cs, err := h.DB.GetCourseByID(id + 1)
		if err != nil {
			h.NotFoundError(w, r)
			return
		}

		params.Completed = true
		params.Next = fmt.Sprintf("/module/%s/1", cs.TemplateName)
	}

	fmt.Println(params)

	b := bytes.NewBufferString("")
	view.Render(b, "quiz", params)
	renderMenu(w, r, "Quiz", template.HTML(b.String()))
}
