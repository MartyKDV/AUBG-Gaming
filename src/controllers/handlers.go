package controllers

import (
	"main/src/models"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

func (server *Server) handleProducts(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("./static/products.html")
	checkError(err)

	var products []models.Item
	results, err := server.Db.Query("SELECT * FROM products")
	checkError(err)

	for results.Next() {
		var p models.Item
		err = results.Scan(&p.Id, &p.Name, &p.Price, &p.Discount, &p.Genre, &p.ReleaseDate, &p.Features, &p.Category)
		checkError(err)
		products = append(products, p)
	}

	err = templ.Execute(w, products)
	checkError(err)
}

func (server *Server) handleProductID(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("./static/product_details.html")
	checkError(err)

	vars := mux.Vars(r)
	id := vars["id"]

	var item models.Item

	result := server.Db.QueryRow("SELECT * FROM products WHERE id = " + id)
	err = result.Scan(&item.Id, &item.Name, &item.Price, &item.Discount, &item.Genre, &item.ReleaseDate, &item.Features, &item.Category)
	checkError(err)

	err = templ.Execute(w, item)
	checkError(err)
}
