package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Task struct {
	ID         int64  `db:"id"`
	TaskText   string `db:"task_text"`
	InsertTime int64  `db:"insert_time"`
}

func (database *Database) GetTaskByID(id int64) (*Task, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetTaskByID: " + err.Error()
		return nil, errors.New(msg)
	}
	u, err := database.getTaskByID(tx, id)
	if err != nil {
		msg := "cannot get task in GetTaskByID: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetTaskByID: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetTaskByID: " + err.Error()
		return nil, errors.New(msg)
	}
	return u, nil
}

func (database *Database) getTaskByID(tx *sqlx.Tx, id int64) (*Task, error) {
	cols := `id, task_text, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM task WHERE id = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getTaskByID: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(id)
	var t Task
	err = row.StructScan(&t)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal user from getTaskByID: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &t, nil
}

func (database *Database) InsertTask(task_text string) error {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for InsertTask: " + err.Error()
		return errors.New(msg)
	}
	err = database.insertTask(tx, task_text)
	if err != nil {
		msg := "cannot insert task in InsertTask: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback from insert in InsertTask: " + msg + ": " + roll_err.Error()
			return errors.New(fatal)
		}
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in InsertUser: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) insertTask(tx *sqlx.Tx, task_text string) error {
	cols := `task_text`
	query := fmt.Sprintf(`INSERT INTO task (%s) VALUES($1)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in insertTask: " + err.Error()
		return errors.New(msg)
	}
	defer stmt.Close()
	_, err = stmt.Exec(task_text)
	if err != nil {
		msg := "cannot execute query in insertTask: " + err.Error()
		return errors.New(msg)
	}
	return nil
}
