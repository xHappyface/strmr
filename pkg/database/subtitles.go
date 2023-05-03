package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Subtitle struct {
	ID         int64   `db:"id"`
	Subtitle   string  `db:"subtitle"`
	Duration   float64 `db:"duration"`
	InsertTime int64   `db:"insert_time"`
}

func (database *Database) GetSubtitlesByTimeRange(start int64, end int64) ([]Subtitle, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetSubtitlesByTimeRange: " + err.Error()
		return nil, errors.New(msg)
	}
	u, err := database.getSubtitlesByTimeRange(tx, start, end)
	if err != nil {
		msg := "cannot get subtitles in GetSubtitlesByTimeRange: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetSubtitlesByTimeRange: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetSubtitlesByTimeRange: " + err.Error()
		return nil, errors.New(msg)
	}
	return u, nil
}

func (database *Database) getSubtitlesByTimeRange(tx *sqlx.Tx, start int64, end int64) ([]Subtitle, error) {
	cols := `id, subtitle, duration, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM subtitles WHERE insert_time >= $1 AND insert_time <= $2`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getSubtitlesByTimeRange: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	rows, err := stmt.Queryx(start, end)
	if err != nil {
		msg := "cannot query subtitles from getSubtitlesByTimeRange: " + err.Error()
		return nil, errors.New(msg)
	}
	subtitles := []Subtitle{}
	for rows.Next() {
		var s Subtitle
		err = rows.StructScan(&s)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return nil, nil
			default:
				msg := "cannot unmarshal subtitle from getSubtitlesByTimeRange: " + err.Error()
				return nil, errors.New(msg)
			}
		}
		subtitles = append(subtitles, s)
	}
	return subtitles, nil
}

func (database *Database) InsertSubtitle(text string, duration float64) error {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for InsertSubtitle: " + err.Error()
		return errors.New(msg)
	}
	err = database.insertSubtitle(tx, text, duration)
	if err != nil {
		msg := "cannot insert subtitle in InsertSubtitle: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback from insert in InsertSubtitle: " + msg + ": " + roll_err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in InsertSubtitle: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) insertSubtitle(tx *sqlx.Tx, text string, duration float64) error {
	cols := `subtitle, duration`
	query := fmt.Sprintf(`INSERT INTO subtitles (%s) VALUES($1, $2)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in insertSubtitle: " + err.Error()
		return errors.New(msg)
	}
	defer stmt.Close()
	_, err = stmt.Exec(text, duration)
	if err != nil {
		msg := "cannot execute query in insertSubtitle: " + err.Error()
		return errors.New(msg)
	}
	return nil
}
