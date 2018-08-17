package db

import (
	"errors"

	"github.com/covveco/micro-learning/model"
)

func (db *DB) GetUserByID(id int) (*model.User, error) {
	u := model.User{}

	err := db.db.Get(&u, `SELECT id, staff_no, name, nric, password, registered FROM USERS WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (db *DB) GetUserByToken(token string) (*model.User, error) {
	u := model.User{}

	err := db.db.Get(&u, `SELECT id, staff_no, name, nric, password, registered, token FROM USERS WHERE token = $1`, token)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (db *DB) GetUserByStaffNo(staffNo string) (*model.User, error) {
	u := model.User{}

	err := db.db.Get(&u, `SELECT id, staff_no, name, nric, password, registered, token FROM USERS WHERE staff_no = $1`, staffNo)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (db *DB) SetUserPassword(u *model.User, p string) error {
	_, err := db.db.Exec(`UPDATE USERS SET password = $1, registered = TRUE WHERE staff_no = $2`, p, u.StaffNo)
	if err != nil {
		return err
	}
	u.Password = model.String(p)

	return nil
}

func (db *DB) CreateUser(u *model.User) error {
	_, err := db.db.Exec(`INSERT INTO users (staff_no, name, nric, token) VALUES ($1, $2, $3, $4)`, u.StaffNo, u.Name, u.NRIC, u.Token)
	return err
}

func (db *DB) DeleteUserByID(id int) error {
	res, err := db.db.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	r, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if r == 0 {
		return errors.New(`no rows affected`)
	}
	return nil
}
