package models

import (
	"database/sql"

	"github.com/blockloop/scan"
)

type Company struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Country string `json:"country"`
	Website string `json:"website"`
	Phone   string `json:"phone"`
}

func (c *Company) Create(db *sql.DB) (*Company, error) {
	stmt := `
	CREATE TABLE IF NOT EXISTS company (
	  name TEXT UNIQUE NOT NULL,
	  code TEXT,
	  country TEXT,
	  website TEXT,
	  phone TEXT);`
	errW := db.QueryRow(stmt)
	if errW.Err() != nil {
		return &Company{}, errW.Err()
	}

	sqlStatement := `
	INSERT INTO company
	VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(sqlStatement, c.Name, c.Code, c.Country, c.Website, c.Phone)
	if err != nil {
		return &Company{}, err
	}

	return c, nil
}

func (c *Company) FindAll(db *sql.DB) (*[]Company, error) {
	var err error
	company := []Company{}
	rows, err := db.Query("SELECT * FROM company")
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	err = scan.Rows(&company, rows)
	if err != nil {
		return nil, err
	}
	return &company, nil
}

func (c *Company) FindByAttr(db *sql.DB, attr string) (*Company, error) {
	var err error
	rows, err := db.Query("SELECT * FROM company WHERE name=$1", attr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	err = scan.Row(c, rows)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Company) Update(db *sql.DB) (*Company, error) {
	stmt := `
	UPDATE company
	SET code = $2, country = $3, website= $4, phone=$5
	WHERE name = $1
	RETURNING *;`

	row := db.QueryRow(stmt, c.Name, c.Code, c.Country, c.Website, c.Phone)
	err := row.Scan(&c.Name, &c.Code, &c.Country, &c.Website, &c.Phone)
	switch err {
	case sql.ErrNoRows:
		// nothig to update
		return &Company{}, sql.ErrNoRows
	case nil:
		break
	default:
		return nil, err
	}

	return c, nil
}

func (c *Company) Delete(db *sql.DB, attr string) (int64, error) {
	companyToDelete := c.Name
	res, err := db.Exec("DELETE from company WHERE name=$1", companyToDelete)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
