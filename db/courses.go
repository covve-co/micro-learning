package db

import "gitlab.com/covveco/special-needs/model"

func (db *DB) GetCourses() ([]model.Course, error) {
	cc := []model.Course{}

	err := db.db.Select(&cc, `SELECT id, title, description, template_name, num_sections FROM courses`)
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func (db *DB) GetCourseByTemplateName(name string) (*model.Course, error) {
	c := model.Course{}

	err := db.db.Get(&c, `SELECT id, title, description, template_name, num_sections FROM courses WHERE template_name = $1`, name)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (db *DB) GetCourseByID(id int) (*model.Course, error) {
	c := model.Course{}

	err := db.db.Get(&c, `SELECT id, title, description, template_name, num_sections FROM courses WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
