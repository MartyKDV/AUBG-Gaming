package controllers

import (
	"fmt"
	"log"
	"main/src/models"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

var signKey = []byte(os.Getenv("JWT_TOKEN"))

func createJWT(user string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = user
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func (server *Server) handleCart(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		{
			// Get cart view
			templ, err := template.ParseFiles("./views/cart.html")
			checkError(err)

			cookie, err := r.Cookie("loginCookie")
			checkError(err)

			user := cookie.Value
			cart := server.GetCart(user)

			var cartItems []models.CartItemDetails

			for _, i := range cart.CartItems {

				result := server.Db.QueryRow("SELECT id, name FROM products WHERE id = " + string(i.ItemID))
				var cartItem models.CartItemDetails
				err := result.Scan(&cartItem.ID, &cartItem.Name, &cartItem.Quantity)
				checkError(err)
				cartItems = append(cartItems, cartItem)
			}

			err = templ.Execute(w, cartItems)
			checkError(err)
		}

	case "POST":
		{
			// Add item to cart
			vars := mux.Vars(r)
			id := vars["id"]
			intID, err := strconv.Atoi(id)
			checkError(err)

			cookie, err := r.Cookie("loginCookie")
			checkError(err)

			user := cookie.Value
			cart := server.GetCart(user)

			cartItem := models.CartItem{ItemID: intID, Quantity: 1}
			cart.CartItems = append(cart.CartItems, cartItem)

			fmt.Fprintf(w, "Added: "+string(cartItem.ItemID))
		}
	case "UPDATE":
		{
			//EditQuantityItem()
		}
	case "DELETE":
		{
			//DeleteItem()
		}
	}
}

func (server *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		{
			templ, err := template.ParseFiles("./views/register.html")
			checkError(err)

			err = templ.Execute(w, nil)
			checkError(err)
		}
	case "POST":
		{
			user := r.FormValue("user")
			pass := []byte(r.FormValue("password"))
			var hashedPassString string
			hashedPass, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
			checkError(err)

			err = bcrypt.CompareHashAndPassword(hashedPass, pass)
			checkError(err)

			hashedPassString = string(hashedPass)
			_, err = server.Db.Exec("INSERT INTO credentials (user, password) VALUES ('" + user + "', '" + hashedPassString + "')")
			checkError(err)

			response := "Successfully Registered!"
			templ, err := template.ParseFiles("./views/status.html")
			checkError(err)
			templ.Execute(w, response)
		}
	}
}
func (server *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case "GET":
		{

			templ, err := template.ParseFiles("./views/login.html")
			checkError(err)

			err = templ.Execute(w, nil)
			checkError(err)
		}
	case "POST":
		{
			user := r.FormValue("user")
			pass := []byte(r.FormValue("password"))
			var hashedPass string
			var response string

			result := server.Db.QueryRow("SELECT password FROM credentials WHERE user = '" + user + "'")
			err := result.Scan(&hashedPass)
			if err != nil {
				response = "Invalid User"
				log.Printf(err.Error())
			} else {
				err = bcrypt.CompareHashAndPassword([]byte(hashedPass), pass)
				if err != nil {
					response = "Invalid Password!"
				} else {
					response = "Successful Login!"
				}
			}

			// Create Authentication Token
			token, err := createJWT(user)
			checkError(err)

			http.SetCookie(w, &http.Cookie{
				Name:    "loginCookie",
				Value:   token,
				Expires: time.Now().Add(time.Minute * 30),
			})

			templ, err := template.ParseFiles("./views/status.html")
			checkError(err)
			templ.Execute(w, response)
		}
	}
}
func (server *Server) handleProducts(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("./views/products.html")
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

	templ, err := template.ParseFiles("./views/product_details.html")
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
func (server *Server) isLogged(path func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if cookie, err := r.Cookie("loginCookie"); err == nil {

			token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Invalid Signing Method!")
				}
				return signKey, nil
			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {

				path(w, r)
			}
		} else {

			response := "Unauthorized! Please log in the system!"
			templ, err := template.ParseFiles("./views/status.html")
			checkError(err)
			templ.Execute(w, response)
		}
	})
}
