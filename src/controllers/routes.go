package controllers

import "net/http"

func (server *Server) initialiseRoutes() {

	fileServer := http.FileServer(http.Dir("../static"))

	server.Router.Handle("/", fileServer)
	server.Router.HandleFunc("/products", server.handleProducts)
	server.Router.HandleFunc("/products/{id}", server.handleProductsID)
}
