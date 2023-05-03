package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// User represents a user in the database.
type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// Order represents an order in the database.
type Order struct {
	ID              int
	Address         string
	ApartmentNumber string
	City            string
	PostalCode      string
}

// Transaction represents a transaction in the database.
type Transaction struct {
	ID        int
	UserID    int
	OrderID   int
	OrderType string
}

func main() {
	// Connect to the database.
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Define the HTTP handlers.
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Parse the form data.
			r.ParseForm()
			firstName := r.FormValue("firstname")
			lastName := r.FormValue("lastname")
			email := r.FormValue("email")
			phone := r.FormValue("phone")

			// Insert the user into the database.
			result, err := db.Exec("INSERT INTO users (firstname, lastname, email, phone) VALUES (?, ?, ?, ?)", firstName, lastName, email, phone)
			if err != nil {
				log.Fatal(err)
			}

			// Get the ID of the inserted user.
			id, err := result.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}

			// Render the template with the inserted user's data.
			user := &User{ID: int(id), FirstName: firstName, LastName: lastName, Email: email, Phone: phone}
			renderTemplate(w, "user.html", user)
		} else {
			// Render the empty user form.
			renderTemplate(w, "user.html", nil)
		}
	})

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Parse the form data.
			r.ParseForm()
			address := r.FormValue("address")
			apartmentNumber := r.FormValue("apartment_number")
			city := r.FormValue("city")
			postalCode := r.FormValue("postalcode")

			// Insert the order into the database.
			result, err := db.Exec("INSERT INTO orders (address, apartment_number, city, postalcode) VALUES (?, ?, ?, ?)", address, apartmentNumber, city, postalCode)
			if err != nil {
				log.Fatal(err)
			}

			// Get the ID of the inserted order.
			id, err := result.LastInsertId()
			if err != nil {
				log.Fatal(err)
			}

			// Render the template with the inserted order's data.
			order := &Order{ID: int(id), Address: address, ApartmentNumber: apartmentNumber, City: city, PostalCode: postalCode}
			renderTemplate(w, "order.html", order)
		} else {
			// Render the empty order form.
			renderTemplate(w, "order.html", nil)
}
})
http.HandleFunc("/transaction", func(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Parse the form data.
		r.ParseForm()
		userID, _ := strconv.Atoi(r.FormValue("user_id"))
		orderID, _ := strconv.Atoi(r.FormValue("order_id"))
		orderType := r.FormValue("order_type")

		// Insert the transaction into the database.
		result, err := db.Exec("INSERT INTO transaction (user_id, order_id, order_type) VALUES (?, ?, ?)", userID, orderID, orderType)
		if err != nil {
			log.Fatal(err)
		}

		// Get the ID of the inserted transaction.
		id, err := result.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}

		// Render the template with the inserted transaction's data.
		transaction := &Transaction{ID: int(id), UserID: userID, OrderID: orderID, OrderType: orderType}
		renderTemplate(w, "transaction.html", transaction)
	} else {
		// Render the empty transaction form.
		renderTemplate(w, "transaction.html", nil)
	}
})

// Start the HTTP server.
fmt.Println("Listening on port 8080...")
err = http.ListenAndServe(":8080", nil)
if err != nil {
	log.Fatal(err)
}
// renderTemplate renders the specified template with the provided data.
func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
tmpl, err := template.ParseFiles(name)
if err != nil {
log.Fatal(err)
}
err = tmpl.Execute(w, data)
if err != nil {
log.Fatal(err)
}
}
