package controllers

import "net/http"

func (server *Server) initialiseRoutes() {

	fileServer := http.FileServer(http.Dir("./static"))

	server.Router.Handle("/", fileServer)
	server.Router.HandleFunc("/products", server.handleProducts).Methods("GET")
	server.Router.HandleFunc("/product/{id}", server.handleProductID).Methods("GET")
}
