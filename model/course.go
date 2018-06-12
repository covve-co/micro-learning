package model

type Course struct {
	ID           int    `db:"id"`
	Title        string `db:"title"`
	Description  string `db:"description"`
	TemplateName string `db:"template_name"`
	NumSections  int    `db:"num_sections"`
}
