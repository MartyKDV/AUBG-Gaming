package controllers

import "net/http"

func (server *Server) initialiseRoutes() {

	fileServer := http.FileServer(http.Dir("./views"))

	server.Router.Handle("/", fileServer)
	server.Router.Handle("/products", server.isLogged(server.handleProducts)).Methods("GET")
	server.Router.Handle("/product/{id}", server.isLogged(server.handleProductID)).Methods("GET")
	server.Router.HandleFunc("/register", server.handleRegister)
	server.Router.HandleFunc("/login", server.handleLogin)
	server.Router.Handle("/cart", server.isLogged(server.handleCart)).Methods("GET")
	server.Router.Handle("/cart/{id}", server.isLogged(server.handleCart)).Methods("POST")
}
