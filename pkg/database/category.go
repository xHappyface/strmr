package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Category struct {
	ID           int64  `db:"id"`
	CategoryName string `db:"category_name"`
	RelatedID    string `db:"related_id"`
	InsertTime   int64  `db:"insert_time"`
}

func (database *Database) GetCategoryByID(id int64) (*Category, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetCategoryByID: " + err.Error()
		return nil, errors.New(msg)
	}
	c, err := database.getCategoryByID(tx, id)
	if err != nil {
		msg := "cannot get category in GetCategoryByID: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetCategoryByID: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetCategoryByID: " + err.Error()
		return nil, errors.New(msg)
	}
	return c, nil
}

func (database *Database) getCategoryByID(tx *sqlx.Tx, id int64) (*Category, error) {
	cols := `id, category_name, related_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM category WHERE id = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getCategoryByID: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(id)
	var c Category
	err = row.StructScan(&c)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal category from getCategoryByID: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &c, nil
}

func (database *Database) InsertCategory(category_name string) error {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for InsertCategory: " + err.Error()
		return errors.New(msg)
	}
	err = database.insertCategory(tx, category_name)
	if err != nil {
		msg := "cannot insert category in InsertCategory: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback from insert in InsertCategory: " + msg + ": " + roll_err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in InsertCategory: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) insertCategory(tx *sqlx.Tx, category_name string) error {
	cols := `category_name`
	query := fmt.Sprintf(`INSERT INTO category (%s) VALUES($1)`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in insertCategory: " + err.Error()
		return errors.New(msg)
	}
	defer stmt.Close()
	_, err = stmt.Exec(category_name)
	if err != nil {
		msg := "cannot execute query in insertCategory: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) GetCategoryByName(category_name string) (*Category, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetCategoryByName: " + err.Error()
		return nil, errors.New(msg)
	}
	c, err := database.getCategoryByName(tx, category_name)
	if err != nil {
		msg := "cannot get category in GetCategoryByName: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetCategoryByName: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetCategoryByName: " + err.Error()
		return nil, errors.New(msg)
	}
	return c, nil
}

func (database *Database) getCategoryByName(tx *sqlx.Tx, category_name string) (*Category, error) {
	cols := `id, category_name, related_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM category WHERE category_name = $1`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getCategoryByName: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	row := stmt.QueryRowx(category_name)
	var c Category
	err = row.StructScan(&c)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, nil
		default:
			msg := "cannot unmarshal category from getCategoryByName: " + err.Error()
			return nil, errors.New(msg)
		}
	}
	return &c, nil
}

func (database *Database) GetAllCategories() ([]Category, error) {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for GetAllCategories: " + err.Error()
		return nil, errors.New(msg)
	}
	c, err := database.getAllCategories(tx)
	if err != nil {
		msg := "cannot get metadata in GetAllCategories: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in GetAllCategories: " + msg + ": " + roll_err.Error()
			return nil, errors.New(fatal)
		}
		return nil, errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in GetAllCategories: " + err.Error()
		return nil, errors.New(msg)
	}
	return c, nil
}

func (database *Database) getAllCategories(tx *sqlx.Tx) ([]Category, error) {
	cols := `id, category_name, related_id, insert_time`
	query := fmt.Sprintf(`SELECT %s FROM category`, cols)
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in getAllCategories: " + err.Error()
		return nil, errors.New(msg)
	}
	defer stmt.Close()
	rows, err := stmt.Queryx()
	if err != nil {
		msg := "cannot query metadata from getAllCategories: " + err.Error()
		return nil, errors.New(msg)
	}
	categories := []Category{}
	for rows.Next() {
		var c Category
		err = rows.StructScan(&c)
		if err != nil {
			switch err {
			case sql.ErrNoRows:
				return nil, nil
			default:
				msg := "cannot unmarshal metadata from getAllCategories: " + err.Error()
				return nil, errors.New(msg)
			}
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (database *Database) UpdateCategoryByName(related_id string, category_name string) error {
	tx, err := database.db.Beginx()
	if err != nil {
		msg := "cannot begin transaction for UpdateCategoryByName: " + err.Error()
		return errors.New(msg)
	}
	err = database.updateCategoryByName(tx, related_id, category_name)
	if err != nil {
		msg := "cannot update category in UpdateCategoryByName: " + err.Error()
		roll_err := tx.Rollback()
		if roll_err != nil {
			fatal := "cannot rollback in UpdateCategoryByName: " + msg + ": " + roll_err.Error()
			return errors.New(fatal)
		}
		return errors.New(msg)
	}
	err = tx.Commit()
	if err != nil {
		msg := "cannot commit transaction in UpdateCategoryByName: " + err.Error()
		return errors.New(msg)
	}
	return nil
}

func (database *Database) updateCategoryByName(tx *sqlx.Tx, related_id string, category_name string) error {
	query := `UPDATE category SET related_id = $1 WHERE category_name = $2`
	stmt, err := tx.Preparex(query)
	if err != nil {
		msg := "cannot prepare statement in updateCategoryByName: " + err.Error()
		return errors.New(msg)
	}
	defer stmt.Close()
	_, err = stmt.Exec(related_id, category_name)
	if err != nil {
		msg := "cannot execute query in updateCategoryByName: " + err.Error()
		return errors.New(msg)
	}
	return nil
}
