package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type StreamMetadata struct {
	ID            int64  `db:"id"`
	MetadataKey   string `db:"metadata_key"`
	MetadataValue string `db:"metadata_value"`
	InsertTime    int64  `db:"insert_time"`
}

func (database *Database) GetStreamMetadataByID(id int64) (*StreamMetadata, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetStreamMetadataByID: " + err.Error()
		return nil, errors.New(msg)
	}
	u, err := database.getStreamMetadataByID(tx, id)
	if err != nil {
		msg := "cannot get stream metadata in GetStreamMetadataByID: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetStreamMetadataByID: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetStreamMetadataByID: " + err.Error()
		return nil, errors.New(msg)
	}
	return u, nil
}

func (database *Database) getStreamMetadataByID(tx *sqlx.Tx, id int64) (*StreamMetadata, error) {
	cols := `id, metadata_key, metadata_value, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM stream_metadata WHERE id = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getStreamMetadataByID: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(id)
	var sm StreamMetadata
	err = row.StructScan(&sm)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal stream metadata from getStreamMetadataByID: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &sm, nil
}

func (database *Database) InsertStreamMetadata(metadata_key string, metadata_value string) error {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for InsertStreamMetadata: " + err.Error()
		return errors.New(msg)
	}
	err = database.insertStreamMetadata(tx, metadata_key, metadata_value)
	if err != nil {
		msg := "cannot insert stream metadata in InsertStreamMetadata: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback from insert in InsertStreamMetadata: " + msg + ": " + roll_err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in InsertStreamMetadata: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) insertStreamMetadata(tx *sqlx.Tx, metadata_key string, metadata_value string) error {
	cols := `metadata_key, metadata_value`
	query := fmt.Sprintf(`INSERT INTO stream_metadata (%s) VALUES($1, $2)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in insertStreamMetadata: " + err.Error()
		return errors.New(msg)
	}
	defer stmt.Close()
	_, err = stmt.Exec(metadata_key, metadata_value)
	if err != nil {
		msg := "cannot execute query in insertStreamMetadata: " + err.Error()
		return errors.New(msg)
	}
	return nil
}
