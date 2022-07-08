package controllers

import (
	"bytes"
	"companies/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var server = Server{}

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

func refreshCompanyTable() error {
	stmt := `drop table if exists company;`
	row := server.DB.QueryRow(stmt)
	if row.Err() != nil {
		return row.Err()
	}

	log.Printf("Successfully refreshed table")
	return nil
}

func refreshAuthTable() error {
	stmt := `drop table if exists auth;`
	row := server.DB.QueryRow(stmt)
	if row.Err() != nil {
		return row.Err()
	}

	log.Printf("Successfully refreshed table")
	return nil
}

func seedCompanies() error {
	companies := []models.Company{
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

func seedOneCompany() (models.Company, error) {
	err := refreshCompanyTable()
	if err != nil {
		log.Fatal(err)
	}
	c := models.Company{
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
		return models.Company{}, errW.Err()
	}

	stmt = `INSERT INTO company (name, code, country, website, phone)
VALUES ($1, $2, $3, $4, $5)`
	errWW := server.DB.QueryRow(stmt, c.Name, c.Code, c.Country, c.Website, c.Phone)
	if errWW.Err() != nil {
		return models.Company{}, errWW.Err()
	}
	return c, nil
}

func seedOneUser() (models.User, error) {
	err := refreshAuthTable()
	if err != nil {
		log.Fatal(err)
	}
	c := models.User{
		Username: gofakeit.Username(),
		Password: gofakeit.Password(true, false, false, false, false, 5),
	}

	stmt := `
	CREATE TABLE IF NOT EXISTS auth (
	  username TEXT UNIQUE NOT NULL,
	  password TEXT);`
	errW := server.DB.QueryRow(stmt)
	if errW.Err() != nil {
		return models.User{}, errW.Err()
	}

	stmt = `INSERT INTO auth (username, password)
VALUES ($1, $2)`
	errWW := server.DB.QueryRow(stmt, c.Username, c.Password)
	if errWW.Err() != nil {
		return models.User{}, errWW.Err()
	}
	return c, nil
}

func TestCreateCompany(t *testing.T) {
	err := refreshCompanyTable()
	if err != nil {
		log.Fatal(err)
	}
	samples := struct {
		inputJSON    string
		statusCode   int
		name         string
		code         string
		country      string
		website      string
		phone        string
		errorMessage string
	}{
		inputJSON:    `{"name":"company", "code":"code", "country":"country", "website":"http://website.com","phone":"12321"}`,
		statusCode:   201,
		name:         "company",
		code:         "code",
		country:      "country",
		website:      "http://website.com",
		phone:        "12321",
		errorMessage: "",
	}

	req, err := http.NewRequest("POST", "/company", bytes.NewBufferString(samples.inputJSON))
	if err != nil {
		t.Error(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.CreateCompany)
	handler.ServeHTTP(rr, req)

	responseMap := make(map[string]interface{})
	err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
	if err != nil {
		fmt.Printf("Cannot convert to json: %v", err)
	}
	assert.Equal(t, rr.Code, samples.statusCode)
	if samples.statusCode == 201 {
		assert.Equal(t, responseMap["name"], samples.name)
		assert.Equal(t, responseMap["code"], samples.code)
		assert.Equal(t, responseMap["country"], samples.country)
		assert.Equal(t, responseMap["website"], samples.website)
		assert.Equal(t, responseMap["phone"], samples.phone)
	}
	if samples.statusCode == 422 || samples.statusCode == 500 && samples.errorMessage != "" {
		assert.Equal(t, responseMap["error"], samples.errorMessage)
	}
}

func TestGetCompanies(t *testing.T) {
	err := refreshCompanyTable()
	if err != nil {
		log.Fatal(err)
	}
	err = seedCompanies()
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/company", nil)
	if err != nil {
		t.Error(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetCompanies)
	handler.ServeHTTP(rr, req)

	var c []models.Company
	err = json.Unmarshal(rr.Body.Bytes(), &c)
	if err != nil {
		log.Fatalf("Cannot convert to json: %v\n", err)
	}
	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(c), 2)
}

func TestGetCompanyByName(t *testing.T) {
	err := refreshCompanyTable()
	if err != nil {
		log.Fatal(err)
	}
	c, err := seedOneCompany()
	if err != nil {
		log.Fatal(err)
	}
	sample := []struct {
		statusCode   int
		name         string
		code         string
		country      string
		website      string
		phone        string
		errorMessage string
	}{
		{
			statusCode: 200,
			name:       c.Name,
			code:       c.Code,
			country:    c.Country,
			website:    c.Website,
			phone:      c.Phone,
		},
		{
			statusCode: 500,
		},
	}
	for _, v := range sample {

		req, err := http.NewRequest("GET", "/company", nil)
		if err != nil {
			t.Error(err)
		}
		req = mux.SetURLVars(req, map[string]string{"name": v.name})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetCompany)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal(rr.Body.Bytes(), &responseMap)
		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, c.Name, responseMap["name"])
			assert.Equal(t, c.Code, responseMap["code"])
			assert.Equal(t, c.Country, responseMap["country"])
			assert.Equal(t, c.Website, responseMap["website"])
			assert.Equal(t, c.Phone, responseMap["phone"])
		}
	}
}
