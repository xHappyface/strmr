package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID         int64  `db:"id"`
	UserID     string `db:"user_id"`
	UserType   string `db:"user_type"`
	InsertTime int64  `db:"insert_time"`
}

func (database *Database) GetUserByID(id int64) (*User, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetUserByID: " + err.Error()
		return nil, errors.New(msg)
	}
	u, err := database.getUserByID(tx, id)
	if err != nil {
		msg := "cannot get user in GetUserByID: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetUserByID: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetUserByID: " + err.Error()
		return nil, errors.New(msg)
	}
	return u, nil
}

func (database *Database) getUserByID(tx *sqlx.Tx, id int64) (*User, error) {
	cols := `id, user_id, user_type, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM user WHERE id = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getUserByID: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(id)
	var u User
	err = row.StructScan(&u)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal user from getUserByID: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &u, nil
}

func (database *Database) GetUserByTypeAndID(user_type, user_id string) (*User, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetUserByTypeAndID: " + err.Error()
		return nil, errors.New(msg)
	}
	u, err := database.getUserByTypeAndID(tx, user_type, user_id)
	if err != nil {
		msg := "cannot get user in GetUserByTypeAndID: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetUserByTypeAndID: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetUserByTypeAndID: " + err.Error()
		return nil, errors.New(msg)
	}
	return u, nil
}

func (database *Database) getUserByTypeAndID(tx *sqlx.Tx, user_type, user_id string) (*User, error) {
	cols := `id, user_id, user_type, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM user WHERE user_type = $1 AND user_id = $2`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getUserByTypeAndID: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(user_type, user_id)
	var u User
	err = row.StructScan(&u)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal user from getUserByTypeAndID: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &u, nil
}

func (database *Database) InsertUser(user_type string, user_id string) error {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for InsertUser: " + err.Error()
		return errors.New(msg)
	}
	u, err := database.getUserByTypeAndID(tx, user_type, user_id)
	if err != nil {
		msg := "cannot get user in InsertUser: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback from select in InsertUser: " + msg + ": " + roll_err.Error()
			return errors.New(fatal)
		}
	}
	if u == nil {
		err := database.insertUser(tx, user_type, user_id)
		if err != nil {
			msg := "cannot insert user in InsertUser: " + err.Error()
			roll_err := tx.Rollback()
			if roll_err != nil {
				fatal := "cannot rollback from insert in InsertUser: " + msg + ": " + roll_err.Error()
				return errors.New(fatal)
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in InsertUser: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) insertUser(tx *sqlx.Tx, user_type string, user_id string) error {
	cols := `user_type, user_id`
	query := fmt.Sprintf(`INSERT INTO user (%s) VALUES(LOWER($1), LOWER($2)) ON CONFLICT (%s) DO NOTHING`, cols, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in insertUser: " + err.Error()
		return errors.New(msg)
	}
	defer stmt.Close()
	_, err = stmt.Exec(user_type, user_id)
	if err != nil {
		msg := "cannot execute query in insertUser: " + err.Error()
		return errors.New(msg)
	}
	return nil
}
