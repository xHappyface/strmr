package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Metadata struct {
	ID            int64  `db:"id"`
	MetadataKey   string `db:"metadata_key"`
	MetadataValue string `db:"metadata_value"`
	InsertTime    int64  `db:"insert_time"`
}

func (database *Database) GetMetadataByID(id int64) (*Metadata, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetMetadataByID: " + err.Error()
		return nil, errors.New(msg)
	}
	sm, err := database.getMetadataByID(tx, id)
	if err != nil {
		msg := "cannot get metadata in GetMetadataByID: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetMetadataByID: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetMetadataByID: " + err.Error()
		return nil, errors.New(msg)
	}
	return sm, nil
}

func (database *Database) getMetadataByID(tx *sqlx.Tx, id int64) (*Metadata, error) {
	cols := `id, metadata_key, metadata_value, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM metadata WHERE id = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getMetadataByID: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(id)
	var m Metadata
	err = row.StructScan(&m)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal metadata from getMetadataByID: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &m, nil
}

func (database *Database) InsertMetadata(metadata_key string, metadata_value string) error {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for InsertMetadata: " + err.Error()
		return errors.New(msg)
	}
	err = database.insertMetadata(tx, metadata_key, metadata_value)
	if err != nil {
		msg := "cannot insert metadata in InsertMetadata: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback from insert in InsertMetadata: " + msg + ": " + roll_err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in InsertMetadata: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) insertMetadata(tx *sqlx.Tx, metadata_key string, metadata_value string) error {
	cols := `metadata_key, metadata_value`
	query := fmt.Sprintf(`INSERT INTO metadata (%s) VALUES($1, $2)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in insertMetadata: " + err.Error()
		return errors.New(msg)
	}
	defer stmt.Close()
	_, err = stmt.Exec(metadata_key, metadata_value)
	if err != nil {
		msg := "cannot execute query in insertMetadata: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) GetLatestMetadataByKey(metadata_key string) (*Metadata, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetLatestMetadataByKey: " + err.Error()
		return nil, errors.New(msg)
	}
	m, err := database.getLatestMetadataByKey(tx, metadata_key)
	if err != nil {
		msg := "cannot get metadata in GetLatestMetadataByKey: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetLatestMetadataByKey: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetLatestMetadataByKey: " + err.Error()
		return nil, errors.New(msg)
	}
	return m, nil
}

func (database *Database) getLatestMetadataByKey(tx *sqlx.Tx, metadata_key string) (*Metadata, error) {
	cols := `id, metadata_key, metadata_value, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM metadata WHERE metadata_key = $1 ORDER BY insert_time DESC LIMIT 1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getLatestMetadataByKey: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(metadata_key)
	var m Metadata
	err = row.StructScan(&m)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal metadata from getLatestMetadataByKey: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &m, nil
}
