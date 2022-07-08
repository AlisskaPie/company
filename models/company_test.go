package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

type Server struct {
	DB     *sql.DB
	Router *mux.Router
}

var server = Server{}
var companyInstance = Company{}

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error
	DBURL := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_USER"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("TEST_DB_PASSWORD"),
	)
	server.DB, err = sql.Open("postgres", DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to postgres database\n")
		log.Fatal(err)
	} else {
		fmt.Printf("connected to the postgres database\n")
	}
}

func TestFindAllCompanies(t *testing.T) {
	err := refreshCompanyTable()
	if err != nil {
		log.Fatal(err)
	}

	err = seedCompanies()
	if err != nil {
		log.Fatal(err)
	}

	companies, err := companyInstance.FindAll(server.DB)
	assert.NoError(t, err)
	assert.Equal(t, len(*companies), 2)
}

func TestCreate(t *testing.T) {
	err := refreshCompanyTable()
	if err != nil {
		log.Fatal(err)
	}
	c := Company{
		Name:    gofakeit.Company(),
		Code:    gofakeit.Password(true, false, false, false, false, 5),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
		Phone:   gofakeit.Phone(),
	}
	created, err := c.Create(server.DB)
	assert.NoError(t, nil)
	assert.Equal(t, c, *created)
}

func TestFindByAttr(t *testing.T) {
	c, err := seedOneCompany()
	if err != nil {
		log.Fatalf("cannot seed table: %v", err)
	}
	found, err := companyInstance.FindByAttr(server.DB, c.Name)

	assert.NoError(t, nil)
	assert.Equal(t, c, *found)
}

func TestUpdate(t *testing.T) {
	c, err := seedOneCompany()
	if err != nil {
		log.Fatalf("Cannot seed company: %v\n", err)
	}

	update := Company{
		Name:    c.Name,
		Code:    gofakeit.Password(true, false, false, false, false, 5),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
		Phone:   gofakeit.Phone(),
	}
	updated, err := update.Update(server.DB)

	assert.NoError(t, nil)
	assert.Equal(t, *updated, update)
}

func TestDelete(t *testing.T) {
	c, err := seedOneCompany()

	if err != nil {
		log.Fatalf("Cannot seed company: %v\n", err)
	}

	isDeleted, err := c.Delete(server.DB, c.Name)
	assert.NoError(t, nil)
	assert.Equal(t, isDeleted, int64(1))
}

func refreshCompanyTable() error {
	stmt := `drop table if exists company;`
	row := server.DB.QueryRow(stmt)
	if row.Err() != nil {
		return row.Err()
	}

	log.Printf("Successfully refreshed table")
	return nil
}

func seedCompanies() error {
	companies := []Company{
		{
			Name:    gofakeit.Company(),
			Code:    gofakeit.Password(true, false, false, false, false, 5),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
			Phone:   gofakeit.Phone(),
		},
		{
			Name:    gofakeit.Company(),
			Code:    gofakeit.Password(true, false, false, false, false, 5),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
			Phone:   gofakeit.Phone(),
		},
	}

	for i := range companies {
		c := companies[i]

		stmt := `
	CREATE TABLE IF NOT EXISTS company (
	  name TEXT UNIQUE NOT NULL,
	  code TEXT,
	  country TEXT,
	  website TEXT,
	  phone TEXT);`
		err := server.DB.QueryRow(stmt)
		if err.Err() != nil {
			return err.Err()
		}

		stmt = `INSERT INTO company (name, code, country, website, phone)
VALUES ($1, $2, $3, $4, $5)`
		err = server.DB.QueryRow(stmt, c.Name, c.Code, c.Country, c.Website, c.Phone)
		if err.Err() != nil {
			return err.Err()
		}
	}
	return nil
}

func seedOneCompany() (Company, error) {
	err := refreshCompanyTable()
	if err != nil {
		log.Fatal(err)
	}
	c := Company{
		Name:    gofakeit.Company(),
		Code:    gofakeit.Password(true, false, false, false, false, 5),
		Country: gofakeit.Country(),
		Website: gofakeit.URL(),
		Phone:   gofakeit.Phone(),
	}

	stmt := `
	CREATE TABLE IF NOT EXISTS company (
	  name TEXT UNIQUE NOT NULL,
	  code TEXT,
	  country TEXT,
	  website TEXT,
	  phone TEXT);`
	errW := server.DB.QueryRow(stmt)
	if errW.Err() != nil {
		return Company{}, errW.Err()
	}

	stmt = `INSERT INTO company (name, code, country, website, phone)
VALUES ($1, $2, $3, $4, $5)`
	errWW := server.DB.QueryRow(stmt, c.Name, c.Code, c.Country, c.Website, c.Phone)
	if errWW.Err() != nil {
		return Company{}, errWW.Err()
	}
	return c, nil
}
