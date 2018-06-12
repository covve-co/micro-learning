package model

type User struct {
	ID         int     `db:"id" json:"id"`
	StaffNo    string  `db:"staff_no" json:"staff_no"`
	Name       string  `db:"name" json:"name"`
	NRIC       string  `db:"nric" json:"nric"`
	Password   *string `db:"password"`
	Registered bool    `db:"registered"`
	Token      string  `db:"token" json:"token"`
}
