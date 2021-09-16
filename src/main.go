package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Product struct {
	Id       int
	Name     string
	Price    float64
	Discount float64
	Category string
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

func main() {

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

	var games []Game
	games = append(games, Game{Genre: "rpg", ReleaseDate: "19/07/2004", Product: Product{Id: 1, Name: "trails", Category: "game", Price: 29.99, Discount: 1.00}})
	games = append(games, Game{Genre: "rpg", ReleaseDate: "19/07/2004", Product: Product{Id: 1, Name: "trails", Category: "game", Price: 29.99, Discount: 1.00}})
	games = append(games, Game{Genre: "rpg", ReleaseDate: "19/07/2004", Product: Product{Id: 1, Name: "trails", Category: "game", Price: 29.99, Discount: 1.00}})
	games = append(games, Game{Genre: "rpg", ReleaseDate: "19/07/2004", Product: Product{Id: 1, Name: "trails", Category: "game", Price: 29.99, Discount: 1.00}})
	games = append(games, Game{Genre: "rpg", ReleaseDate: "19/07/2004", Product: Product{Id: 1, Name: "trails", Category: "game", Price: 29.99, Discount: 1.00}})
	games = append(games, Game{Genre: "rpg", ReleaseDate: "19/07/2004", Product: Product{Id: 1, Name: "trails", Category: "game", Price: 29.99, Discount: 1.00}})

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
