package db

import "github.com/covveco/micro-learning/model"

func (db *DB) GetCompletionsByUserID(id int) ([]model.Completion, error) {
	cc := []model.Completion{}

	err := db.db.Select(&cc, `SELECT user_id, course_id FROM completions WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return cc, nil
}

func (db *DB) GetCompletionsByUserIDCourseID(uID, cID int) ([]model.Completion, error) {
	cc := []model.Completion{}

	err := db.db.Select(&cc, `SELECT user_id, course_id FROM completions WHERE user_id = $1 AND course_id = $2`, uID, cID)
	if err != nil {
		return nil, err
	}

	return cc, nil
}

func (db *DB) CreateCompletion(c *model.Completion) error {
	_, err := db.db.Exec(`INSERT INTO completions (user_id, course_id) VALUES ($1, $2)`, c.UserID, c.CourseID)
	return err
}
