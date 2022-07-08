package controllers

import "companies/middleware"

func (server *Server) initializeRoutes() {
	// Login Route
	server.Router.HandleFunc("/login", middleware.JSON(server.Login)).Methods("POST")
	server.Router.HandleFunc("/register", middleware.JSON(server.RegisterUser)).Methods("POST")

	// Company routes
	server.Router.HandleFunc("/company", middleware.Auth(middleware.LocationIP(server.CreateCompany))).Methods("POST")
	server.Router.HandleFunc("/company", middleware.JSON(server.GetCompanies)).Methods("GET")
	server.Router.HandleFunc("/company/{name}", middleware.JSON(server.GetCompany)).Methods("GET")
	server.Router.HandleFunc("/company/{name}", middleware.JSON(server.Update)).Methods("PUT")
	server.Router.HandleFunc("/company/{name}", middleware.Auth(middleware.LocationIP(server.DeleteCompany))).Methods("DELETE")
}
