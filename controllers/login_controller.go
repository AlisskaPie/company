package controllers

import (
	"companies/api/utils"
	"companies/authorization"
	"companies/models"
	"companies/responses"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (server *Server) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
	}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	user.Password = utils.HashPassword(user.Password)
	_, err = server.DB.Exec("insert into auth (username, password) values ($1,$2);", user.Username, user.Password)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusCreated, fmt.Sprintf("username: %s", user.Username))
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	token, err := server.SignIn(user.Username, user.Password)
	if err != nil {
		formattedError := utils.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(username, password string) (string, error) {
	var err error
	user := models.User{}

	sqlStmt := `select * from auth where username=$1`
	err = server.DB.QueryRow(sqlStmt, username).Scan(&user.Username, &user.Password)
	if err != nil {
		return "", err
	}
	ok := utils.CheckPasswordHash(password, user.Password)
	if !ok {
		err = bcrypt.ErrMismatchedHashAndPassword
		return "", err
	}
	return authorization.GenerateJWT(username, password)
}
