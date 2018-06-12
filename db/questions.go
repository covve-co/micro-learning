package db

import (
	"gitlab.com/covveco/special-needs/model"
)

func (db *DB) GetQuestionsByCourseID(id int) ([]model.Question, error) {
	qq := []model.Question{}

	err := db.db.Select(&qq, "SELECT id, course_id, title, content, image FROM questions WHERE course_id = $1", id)
	if err != nil {
		return nil, err
	}
	return qq, nil
}

func (db *DB) GetOptionsByQuestionID(id int) ([]model.Option, error) {
	oo := []model.Option{}

	err := db.db.Select(&oo, "SELECT question_id, position, content, correct FROM options WHERE question_id = $1", id)
	if err != nil {
		return nil, err
	}
	return oo, nil
}
