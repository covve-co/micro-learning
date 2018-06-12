package model

type Question struct {
	ID       int     `db:"id"`
	CourseID int     `db:"course_id"`
	Title    string  `db:"title"`
	Content  string  `db:"content"`
	Image    *string `db:"image"`

	// For the view code
	Options []Option
	Correct bool
}

type Option struct {
	QuestionID int    `db:"question_id"`
	Position   int    `db:"position"`
	Content    string `db:"content"`
	Correct    bool   `db:"correct"`
}
