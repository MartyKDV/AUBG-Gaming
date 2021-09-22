package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Product struct {
	Id              int
	Name            string
	Price           float64
	DiscountedPrice float64
	Discount        float64
}

type Game struct {
	Genre       string
	ReleaseDate string
	Product
}

type Hardware struct {
	Features []string
	Product
}

type Service struct {
	ServiceType  string
	DeliveryType string
	Product
}

func getDiscountedPrice(p Product) float64 {
	return p.Discount * p.Price
}

func main() {

	// Create router and set the functions to handle paths
	router := mux.NewRouter()
	fileServer := http.FileServer(http.Dir("../static"))

	router.Handle("/", fileServer)
	router.HandleFunc("/products", handleProducts)
	router.HandleFunc("/products/{id}", handleProductsID)

	log.Println("Server Has Successfully Started at Port :8080...")
	err := http.ListenAndServe(":8080", router)
	checkError(err)

}

func handleProducts(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("../static/products.html")
	checkError(err)

	// Connect to database
	db, err := sql.Open("mysql", "martykdv:F8>V@44qW(Wh{f*$@tcp(127.0.0.1:3306)/aubg-gaming-db")
	checkError(err)
	defer db.Close()

	var games []Game
	results, err := db.Query("SELECT * FROM games")

	for results.Next() {
		var g Game
		err = results.Scan(&g.Id, &g.Name, &g.Price, &g.Discount, &g.Genre, &g.ReleaseDate)
		checkError(err)
		games = append(games, g)
	}

	err = templ.Execute(w, games)
	checkError(err)
}

func handleProductsID(w http.ResponseWriter, r *http.Request) {

}

func checkError(err error) {

	if err != nil {
		panic(err)
	}
}
