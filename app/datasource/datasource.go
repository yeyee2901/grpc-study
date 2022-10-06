package datasource

import (
	"fmt"

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

    fmt.Printf("ds.DB: %v\n", ds.DB)
    fmt.Printf("newBook: %v\n", newBook)
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

    id, err := result.LastInsertId()

    fmt.Println("Inserted ID: ", id)
    return nil
}
