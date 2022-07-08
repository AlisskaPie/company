package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	DB     *sql.DB
	Router *mux.Router
}

func (server *Server) Initialize(DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		DbHost, DbPort, DbUser, DbName, DbPassword)

	server.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run() {
	fmt.Println("Listening to port 3000")
	log.Fatal(http.ListenAndServe(":3000", server.Router))
}
