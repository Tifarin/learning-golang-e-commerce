package config

import "gotoko/handlers"

func (server *Server) initializeRoutes() {
	server.Router.HandleFunc("/", handlers.Home).Methods("GET")
}