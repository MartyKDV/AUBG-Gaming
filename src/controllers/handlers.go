package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"main/src/models"
	"net/http"
	"os"
	"regexp"
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

func (server *Server) handleFilter(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("./views/products.html")
	checkError(err)
	log.Println("In handle func")
	r.ParseForm()
	filters := r.Form
	log.Println(filters)

	sql := "SELECT * FROM products "
	if len(filters["filter-category"]) > 0 || len(filters["filter-game"]) > 0 || len(filters["filter-hardware"]) > 0 {
		sql += "WHERE "

		extra := false
		if len(filters["filter-category"]) > 0 {
			if extra {
				sql += "AND "
			}
			extra = true
			sql += "category IN ("
			for i, v := range filters["filter-category"] {
				log.Println(v)
				if i == (len(filters["filter-category"]) - 1) {
					sql += "'" + v + "') "
				} else {
					sql += "'" + v + "',"
				}
			}
		}

		if len(filters["filter-game"]) > 0 {
			if extra {
				sql += "AND "
			}
			extra = true
			sql += "genre IN ("
			for i, v := range filters["filter-game"] {
				log.Println(v)
				if i == (len(filters["filter-game"]) - 1) {
					sql += "'" + v + "') "
				} else {
					sql += "'" + v + "',"
				}
			}
		}

		if len(filters["filter-hardware"]) > 0 {
			if extra {
				sql += "AND "
			}
			extra = true
			sql += "hardware_type IN ("
			for i, v := range filters["filter-hardware"] {
				log.Println(v)
				if i == (len(filters["filter-hardware"]) - 1) {
					sql += "'" + v + "') "
				} else {
					sql += "'" + v + "',"
				}
			}
		}
	}
	var products []models.Item
	results, err := server.Db.Query(sql)
	checkError(err)

	for results.Next() {
		var p models.Item
		err = results.Scan(&p.Id, &p.Name, &p.Price, &p.Discount, &p.Genre, &p.ReleaseDate, &p.Features, &p.HardwareType, &p.ServiceType, &p.Category)
		checkError(err)
		log.Println(p)
		products = append(products, p)
	}

	err = templ.Execute(w, products)
	checkError(err)
}

func (server *Server) handleOrder(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("./views/status.html")
	checkError(err)

	// Here should be the transaction code which performs the transaction itself
	// but for this is out of the scope of this project

	intial := os.Getenv("INITIAL_CITY")
	goal := r.FormValue("city")

	path := Search(intial, goal)

	response := "Order is Successful<br>" + path

	err = templ.Execute(w, response)
	checkError(err)
}
func (server *Server) handleCheckout(w http.ResponseWriter, r *http.Request) {

	templ, err := template.ParseFiles("./views/checkout.html")
	checkError(err)

	cookie, err := r.Cookie("cartCookie")
	checkError(err)

	user := cookie.Value
	log.Println("Got: " + user)
	cart := server.GetCart(user)

	var total, price float64 = 0.00, 0.00
	type cartInfo struct {
		CartItems []models.CartItemDetails
		Total     float64
	}
	var cartItems cartInfo

	for _, i := range cart.CartItems {

		idString := strconv.Itoa(i.ItemID)
		result := server.Db.QueryRow("SELECT id, name, price FROM products WHERE id = ?", idString)
		var cartItem models.CartItemDetails
		err := result.Scan(&cartItem.ID, &cartItem.Name, &price)
		checkError(err)
		price = price * float64(i.Quantity)
		total += price
		cartItem.Quantity = i.Quantity
		cartItems.CartItems = append(cartItems.CartItems, cartItem)
	}
	cartItems.Total = total

	err = templ.Execute(w, cartItems)
	checkError(err)
}

// Test handler for the path-finding algorithm
func (server *Server) handleSearch(w http.ResponseWriter, r *http.Request) {

	var g Graph
	g.Initialise("Blagoevgrad", "Silistra")

	path := g.aStarSearch("Blagoevgrad", "Silistra")
	answer := ""
	for i := range path {
		answer += path[i]
	}
	log.Println("Path: " + answer)
}
func (server *Server) handleCartDelete(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]
	intID, err := strconv.Atoi(id)
	checkError(err)
	cookie, err := r.Cookie("cartCookie")
	checkError(err)
	user := cookie.Value

	server.DeleteCartItem(user, intID)

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}
func (server *Server) handleCartUpdate(w http.ResponseWriter, r *http.Request) {

	quantity := r.FormValue("quantity")
	quantityInt, err := strconv.Atoi(quantity)
	checkError(err)
	vars := mux.Vars(r)
	id := vars["id"]
	intID, err := strconv.Atoi(id)
	checkError(err)
	cookie, err := r.Cookie("cartCookie")
	checkError(err)
	user := cookie.Value

	server.UpdateQuantity(user, intID, quantityInt)

	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}
func (server *Server) handleCart(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		{
			// Get cart view
			templ, err := template.ParseFiles("./views/cart.html")
			checkError(err)

			cookie, err := r.Cookie("cartCookie")
			checkError(err)

			user := cookie.Value
			log.Println("Got: " + user)
			cart := server.GetCart(user)

			var cartItems []models.CartItemDetails
			for _, i := range cart.CartItems {

				idString := strconv.Itoa(i.ItemID)
				result := server.Db.QueryRow("SELECT id, name FROM products WHERE id = ?", idString)
				var cartItem models.CartItemDetails
				err := result.Scan(&cartItem.ID, &cartItem.Name)
				checkError(err)
				cartItem.Quantity = i.Quantity
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

			cookie, err := r.Cookie("cartCookie")
			checkError(err)

			user := cookie.Value
			log.Println("Got: " + user)

			cartItem := models.CartItem{ItemID: intID, Quantity: 1}
			server.UpdateCart(user, cartItem)

			templ, err := template.ParseFiles("./views/status.html")
			checkError(err)

			err = templ.Execute(w, "Added an Item to cart!")
			checkError(err)
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
			passString := r.FormValue("password")
			pass := []byte(passString)
			var response string

			matchEmail, _ := regexp.MatchString("^[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$", user)
			if matchEmail {

				matchPass, _ := regexp.MatchString("^[a-zA-Z0-9._%+\\-]+$", passString)
				if matchPass && len(passString) >= 8 {

					result := server.Db.QueryRow("SELECT password FROM credentials WHERE user = ?", user)
					var password string
					err := result.Scan(&password)
					if err != sql.ErrNoRows {
						response = "An Account With This Email Already Exists!"
					} else {

						var hashedPassString string
						hashedPass, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
						checkError(err)

						err = bcrypt.CompareHashAndPassword(hashedPass, pass)
						checkError(err)

						hashedPassString = string(hashedPass)
						// Validate input before sql prepared statement
						_, err = server.Db.Exec("INSERT INTO credentials (user, password) VALUES (?, ?)", user, hashedPassString)
						checkError(err)

						response = "Successfully Registered!"
					}

				} else {

					response = "Invalid Password Format:\nPassword must be at least 8 characters long, consisting of lowercase and uppercase letters, numbers, and special characters"
				}
			} else {
				response = "Invalid Email Format"
			}

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

			result := server.Db.QueryRow("SELECT password FROM credentials WHERE user = ?", user)
			err := result.Scan(&hashedPass)
			if err != nil {
				response = "Invalid User"
				log.Println(err.Error())
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

			http.SetCookie(w, &http.Cookie{
				Name:    "cartCookie",
				Value:   user,
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
		err = results.Scan(&p.Id, &p.Name, &p.Price, &p.Discount, &p.Genre, &p.ReleaseDate, &p.Features, &p.HardwareType, &p.ServiceType, &p.Category)
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

	result := server.Db.QueryRow("SELECT * FROM products WHERE id = ?", id)
	err = result.Scan(&item.Id, &item.Name, &item.Price, &item.Discount, &item.Genre, &item.ReleaseDate, &item.Features, &item.HardwareType, &item.ServiceType, &item.Category)
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
