package controllers

import (
	"database/sql"
	"log"
	"main/src/models"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Server struct {
	Db           *sql.DB
	Router       *mux.Router
	CartSessions map[string]models.Cart
}

func (server *Server) GetCart(user string) models.Cart {

	if _, exists := server.CartSessions[user]; exists {
		return server.CartSessions[user]
	} else {
		server.CartSessions[user] = models.Cart{}
		return server.CartSessions[user]
	}
}

func (server *Server) UpdateCart(user string, cartItem models.CartItem) {

	cart := server.CartSessions[user]
	cart.CartItems = append(cart.CartItems, cartItem)
	server.CartSessions[user] = cart
}

func (server *Server) Initialise(dns string) {

	var err error
	server.Db, err = sql.Open("mysql", dns)
	checkError(err)

	server.Router = mux.NewRouter()
	server.CartSessions["mbk180@aubg.edu"] = models.Cart{CartItems: []models.CartItem{{ItemID: 1, Quantity: 1}}}

	server.initialiseRoutes()

}

func (server *Server) Run(host string) {

	defer server.Db.Close()
	log.Println("Server Has Successfully Started at Port :8080...")
	err := http.ListenAndServe(host, server.Router)
	checkError(err)
}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
