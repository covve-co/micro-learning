package model

import "time"

type Completion struct {
	UserID    int       `db:"user_id" json:"user_id"`
	CourseID  int       `db:"course_id" json:"course_id"`
	Timestamp time.Time `db:"timestamp" json:"timestamp"`
}
