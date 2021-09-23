package controllers

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Server struct {
	Db     *sql.DB
	Router *mux.Router
}

func (server *Server) Initialise(dns string) {

	var err error
	server.Db, err = sql.Open("mysql", dns)
	checkError(err)

	server.Router = mux.NewRouter()

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
