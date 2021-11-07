package controllers

import "net/http"

func (server *Server) initialiseRoutes() {

	fileServer := http.FileServer(http.Dir("./views"))

	server.Router.Handle("/", fileServer)
	server.Router.Handle("/products", server.isLogged(server.handleProducts)).Methods("GET")
	server.Router.Handle("/product/{id}", server.isLogged(server.handleProductID)).Methods("GET")
	server.Router.HandleFunc("/register", server.handleRegister)
	server.Router.HandleFunc("/login", server.handleLogin)
	server.Router.Handle("/checkout", server.isLogged(server.handleCheckout)).Methods("GET")
	server.Router.Handle("/order", server.isLogged(server.handleOrder))
	server.Router.HandleFunc("/search", server.handleSearch)
	server.Router.Handle("/cart", server.isLogged(server.handleCart)).Methods("GET")
	server.Router.Handle("/cart/{id}", server.isLogged(server.handleCart)).Methods("POST")
	server.Router.Handle("/cart/quantity/{id}", server.isLogged(server.handleCartUpdate)).Methods("POST")
	server.Router.Handle("/cart/delete/{id}", server.isLogged(server.handleCartDelete)).Methods("POST")
}
