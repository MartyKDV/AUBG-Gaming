package controllers

import (
	"main/src/models"
	"net/http"
	"text/template"
)

func (server *Server) handleProducts(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("./static/products.html")
	checkError(err)

	var games []models.Game
	results, err := server.Db.Query("SELECT * FROM games")

	for results.Next() {
		var g models.Game
		err = results.Scan(&g.Id, &g.Name, &g.Price, &g.Discount, &g.Genre, &g.ReleaseDate)
		checkError(err)
		games = append(games, g)
	}

	err = templ.Execute(w, games)
	checkError(err)
}

func (server *Server) handleProductsID(w http.ResponseWriter, r *http.Request) {

}
