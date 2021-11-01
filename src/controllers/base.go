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
	Db          *sql.DB
	Router      *mux.Router
	sessionBase []Session
}

type Session struct {
	User string
	Cart models.Cart
}

func (server *Server) GetCart(user string) *models.Cart {

	for _, s := range server.sessionBase {
		if user == s.User {
			return &s.Cart
		}
	}
	session := new(Session)
	server.sessionBase = append(server.sessionBase, *session)
	return &session.Cart
}

func (server *Server) UpdateCart(user string, cartItem models.CartItem) {

	for _, s := range server.sessionBase {
		if user == s.User {
			s.Cart.CartItems = append(s.Cart.CartItems, cartItem)
		}
	}
}

func (server *Server) Initialise(dns string) {

	var err error
	server.Db, err = sql.Open("mysql", dns)
	checkError(err)

	server.Router = mux.NewRouter()
	server.sessionBase = []Session{{User: "mbk180@aubg.edu", Cart: models.Cart{CartItems: []models.CartItem{{ItemID: 1, Quantity: 1}, {ItemID: 2, Quantity: 1}}}}}

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
