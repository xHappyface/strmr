package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type MediaRecording struct {
	ID         int64  `db:"id"`
	FileName   string `db:"file_name"`
	Extension  string `db:"extension"`
	Directory  string `db:"directory"`
	StartTime  int64  `db:"start_time"`
	EndTime    *int64 `db:"end_time"`
	Uploaded   int64  `db:"uploaded"`
	InsertTime int64  `db:"insert_time"`
}

func (database *Database) GetMediaRecordingByID(id int64) (*MediaRecording, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetMediaRecordingByID: " + err.Error()
		return nil, errors.New(msg)
	}
	mr, err := database.getMediaRecordingByID(tx, id)
	if err != nil {
		msg := "cannot get media recording in GetMediaRecordingByID: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetMediaRecordingByID: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetMediaRecordingByID: " + err.Error()
		return nil, errors.New(msg)
	}
	return mr, nil
}

func (database *Database) getMediaRecordingByID(tx *sqlx.Tx, id int64) (*MediaRecording, error) {
	cols := `id, file_name, extension, directory, start_time, end_time, uploaded, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM media_recording WHERE id = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getMediaRecordingByID: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(id)
	var mr MediaRecording
	err = row.StructScan(&mr)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal media recording from getMediaRecordingByID: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &mr, nil
}

func (database *Database) GetLatestMediaRecording() (*MediaRecording, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetLatestMediaRecording: " + err.Error()
		return nil, errors.New(msg)
	}
	mr, err := database.getLatestMediaRecording(tx)
	if err != nil {
		msg := "cannot get media recording in GetLatestMediaRecording: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetLatestMediaRecording: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetLatestMediaRecording: " + err.Error()
		return nil, errors.New(msg)
	}
	return mr, nil
}

func (database *Database) getLatestMediaRecording(tx *sqlx.Tx) (*MediaRecording, error) {
	cols := `id, file_name, extension, directory, start_time, end_time, uploaded, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM media_recording WHERE end_time IS NULL ORDER BY insert_time DESC LIMIT 1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getLatestMediaRecording: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx()
	var mr MediaRecording
	err = row.StructScan(&mr)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal media recording from getLatestMediaRecording: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &mr, nil
}

func (database *Database) EndActiveMediaRecordings() error {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for EndActiveMediaRecordings: " + err.Error()
		return errors.New(msg)
	}
	err = database.endActiveMediaRecordings(tx)
	if err != nil {
		msg := "cannot end active media recordings in EndActiveMediaRecordings: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in EndActiveMediaRecordings: " + msg + ": " + roll_err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in EndActiveMediaRecordings: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) endActiveMediaRecordings(tx *sqlx.Tx) error {
	query := `UPDATE media_recording SET end_time = CAST(strftime('%s', 'now') AS INTEGER) WHERE end_time IS NULL`
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in endActiveMediaRecordings: " + err.Error()
		return errors.New(msg)
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		msg := "cannot execute query in endActiveMediaRecordings: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) SetMediaRecordingUploadedByID(id int64, uploaded bool) error {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for SetMediaRecordingUploadedByID: " + err.Error()
		return errors.New(msg)
	}
	err = database.setMediaRecordingUploadedByID(tx, id, uploaded)
	if err != nil {
		msg := "cannot end active media recordings in SetMediaRecordingUploadedByID: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in SetMediaRecordingUploadedByID: " + msg + ": " + roll_err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in SetMediaRecordingUploadedByID: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) setMediaRecordingUploadedByID(tx *sqlx.Tx, id int64, uploaded bool) error {
	query := `UPDATE media_recording SET uploaded = $1 WHERE id = $2`
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in setMediaRecordingUploadedByID: " + err.Error()
		return errors.New(msg)
	}
	defer stmt.Close()
	u := int64(0)
	if uploaded {
		u = int64(1)
	}
	_, err = stmt.Exec(u, id)
	if err != nil {
		msg := "cannot execute query in setMediaRecordingUploadedByID: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) InsertMediaRecording(file_name string, directory string) error {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for InsertMediaRecording: " + err.Error()
		return errors.New(msg)
	}
	err = database.insertMediaRecording(tx, file_name, directory)
	if err != nil {
		msg := "cannot insert media recording in InsertMediaRecording: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback from insert in InsertMediaRecording: " + msg + ": " + roll_err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in InsertMediaRecording: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) insertMediaRecording(tx *sqlx.Tx, file_name string, directory string) error {
	cols := `file_name, directory, end_time`
	query := fmt.Sprintf(`INSERT INTO media_recording (%s) VALUES($1, $2, NULL)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in insertMediaRecording: " + err.Error()
		return errors.New(msg)
	}
	defer stmt.Close()
	_, err = stmt.Exec(file_name, directory)
	if err != nil {
		msg := "cannot execute query in insertMediaRecording: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) GetAllMediaRecordingsByUploaded(uploaded bool) ([]MediaRecording, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetAllMediaRecordingsByUploaded: " + err.Error()
		return nil, errors.New(msg)
	}
	m, err := database.getAllMediaRecordingsByUploaded(tx, uploaded)
	if err != nil {
		msg := "cannot get media recordings in GetAllMediaRecordingsByUploaded: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetAllMediaRecordingsByUploaded: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetAllMediaRecordingsByUploaded: " + err.Error()
		return nil, errors.New(msg)
	}
	return m, nil
}

func (database *Database) getAllMediaRecordingsByUploaded(tx *sqlx.Tx, uploaded bool) ([]MediaRecording, error) {
	cols := `id, file_name, extension, directory, start_time, end_time, uploaded, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM media_recording WHERE uploaded = $1 AND end_time IS NOT NULL`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getAllMediaRecordingsByUploaded: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	u := int64(0)
	if uploaded {
		u = int64(1)
	}
	rows, err := stmt.Queryx(u)
	if err != nil {
		msg := "cannot query media recordings from getAllMediaRecordingsByUploaded: " + err.Error()
		return nil, errors.New(msg)
	}
	media_recordings := []MediaRecording{}
	for rows.Next() {
		var mr MediaRecording
		err = rows.StructScan(&mr)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return nil, nil
			default:
				msg := "cannot unmarshal media recording from getAllMediaRecordingsByUploaded: " + err.Error()
				return nil, errors.New(msg)
			}
		}
		media_recordings = append(media_recordings, mr)
	}
	return media_recordings, nil
}
