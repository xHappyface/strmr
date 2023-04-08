package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Stream struct {
	ID         int64 `db:"id"`
	StartTime  int64 `db:"start_time"`
	EndTime    int64 `db:"end_time"`
	InsertTime int64 `db:"insert_time"`
}

func (database *Database) GetStreamByID(id int64) (*Stream, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetStreamByID: " + err.Error()
		return nil, errors.New(msg)
	}
	s, err := database.getStreamByID(tx, id)
	if err != nil {
		msg := "cannot get task in GetStreamByID: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetStreamByID: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetStreamByID: " + err.Error()
		return nil, errors.New(msg)
	}
	return s, nil
}

func (database *Database) getStreamByID(tx *sqlx.Tx, id int64) (*Stream, error) {
	cols := `id, start_time, end_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM stream WHERE id = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getStreamByID: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(id)
	var s Stream
	err = row.StructScan(&s)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal user from getStreamByID: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &s, nil
}

func (database *Database) GetLatestStream() (*Stream, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetLatestStream: " + err.Error()
		return nil, errors.New(msg)
	}
	s, err := database.getLatestStream(tx)
	if err != nil {
		msg := "cannot get task in GetLatestStream: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetLatestStream: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetLatestStream: " + err.Error()
		return nil, errors.New(msg)
	}
	return s, nil
}

func (database *Database) getLatestStream(tx *sqlx.Tx) (*Stream, error) {
	cols := `id, start_time, end_time, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM stream WHERE end_time IS NULL ORDER BY insert_time DESC LIMIT 1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getLatestStream: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx()
	var s Stream
	err = row.StructScan(&s)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal user from getLatestStream: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &s, nil
}

func (database *Database) EndActiveStreams() error {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for EndActiveStreams: " + err.Error()
		return errors.New(msg)
	}
	err = database.endActiveStreams(tx)
	if err != nil {
		msg := "cannot get task in EndActiveStreams: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in EndActiveStreams: " + msg + ": " + roll_err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in EndActiveStreams: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) endActiveStreams(tx *sqlx.Tx) error {
	query := `UPDATE stream SET end_time = CAST(strftime('%s', 'now') AS INTEGER) WHERE end_time IS NULL`
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in endActiveStreams: " + err.Error()
		return errors.New(msg)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		msg := "cannot execute query in endActiveStreams: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) InsertStream() error {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for InsertStream: " + err.Error()
		return errors.New(msg)
	}
	err = database.insertStream(tx)
	if err != nil {
		msg := "cannot insert stream in InsertStream: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback from insert in InsertStream: " + msg + ": " + roll_err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in InsertStream: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) insertStream(tx *sqlx.Tx) error {
	cols := `end_time`
	query := fmt.Sprintf(`INSERT INTO stream (%s) VALUES(NULL)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in insertStream: " + err.Error()
		return errors.New(msg)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		msg := "cannot execute query in insertStream: " + err.Error()
		return errors.New(msg)
	}
	return nil
}
