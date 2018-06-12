package handler

import (
	"html/template"
	"net/http"

	"gitlab.com/covveco/special-needs/model"
)

func renderMenu(w http.ResponseWriter, r *http.Request, title string, content template.HTML) {
	user, _ := r.Context().Value("user").(*model.User)

	renderView(w, "menu", struct {
		Title    string
		Username string
		Content  template.HTML
	}{
		Title:    title,
		Username: user.Name,
		Content:  content,
	})
}
