package datasource

import (
	"github.com/jmoiron/sqlx"
)

type DataSource struct {
	DB *sqlx.DB
}

func NewDataSource(db *sqlx.DB) *DataSource {
	return &DataSource{
		DB: db,
	}
}

func (ds *DataSource) SaveBook(newBook interface{}) (err error) {
	query := `
    INSERT INTO books
        (title, isbn, tahun)
    VALUES
        (:title, :isbn, :tahun)
    `

	tx := ds.DB.MustBegin()

	result, err := tx.NamedExec(query, newBook)
	if err != nil {
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return
	}

	_, err = result.LastInsertId()

	return nil
}

func (ds *DataSource) GetUserById(result interface{}, id int64) (err error) {
	var args []interface{}
	query := `
        SELECT
            name, email, DATE_TRUNC('second', created_at) AS created_at
        FROM
            users
        WHERE
            id = $1
    `

	args = append(args, id)
	err = ds.DB.Get(result, query, args...)

	return
}
