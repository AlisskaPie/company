package controllers

import (
	"companies/api/utils"
	"companies/models"
	"companies/responses"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

func (server *Server) CreateCompany(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	company := models.Company{}
	err = json.Unmarshal(body, &company)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	created, err := company.Create(server.DB)
	if err != nil {
		err = utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusCreated, created)
}
func (server *Server) GetCompanies(w http.ResponseWriter, r *http.Request) {
	c := models.Company{}
	posts, err := c.FindAll(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, posts)
}
func (server *Server) GetCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	c := models.Company{}

	received, err := c.FindByAttr(server.DB, name)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, received)
}
func (server *Server) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Start processing the request data
	update := models.Company{}
	err = json.Unmarshal(body, &update)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	update.Name = name
	updated, err := update.Update(server.DB)

	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, updated)
}
func (server *Server) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	// Check if the company exist
	c := models.Company{}
	stmt := `SELECT EXISTS (
		SELECT FROM
	company
	WHERE
	name  = $1
	);`

	_, err := server.DB.Exec(stmt, name)
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, nil)
		return
	}

	stmt = `SELECT * FROM 
		company
		WHERE name=$1`
	err = server.DB.QueryRow(stmt, name).Scan(&c.Name, &c.Code, &c.Country, &c.Website, &c.Phone)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, nil)
		return
	}

	_, err = c.Delete(server.DB, name)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, "")
}
