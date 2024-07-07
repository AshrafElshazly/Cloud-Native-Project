package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

type Order struct {
	ID        int `json:"id"`
	UserID    int `json:"user_id"`
	ProductID int `json:"product_id"`
}

type OrderDetail struct {
	OrderID      int     `json:"order_id"`
	UserName     string  `json:"user_name"`
	UserEmail    string  `json:"user_email"`
	ProductName  string  `json:"product_name"`
	ProductPrice float64 `json:"product_price"`
	ProductCount int     `json:"product_count"`
}

var db *sql.DB

func initDB() {
	var err error
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=require"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	createTable := `CREATE TABLE IF NOT EXISTS orders (
        id SERIAL PRIMARY KEY,
        user_id INT NOT NULL,
        product_id INT NOT NULL
    );`
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var order Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", order.UserID).Scan(&userExists)
	if err != nil || !userExists {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var productExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)", order.ProductID).Scan(&productExists)
	if err != nil || !productExists {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	err = db.QueryRow("INSERT INTO orders (user_id, product_id) VALUES ($1, $2) RETURNING id", order.UserID, order.ProductID).Scan(&order.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	params := mux.Vars(r)
	query := `
    SELECT 
        orders.id AS order_id, 
        users.name AS user_name, 
        users.email AS user_email, 
        products.name AS product_name, 
        products.price AS product_price, 
        products.count AS product_count
    FROM orders
    JOIN users ON orders.user_id = users.id
    JOIN products ON orders.product_id = products.id
    WHERE orders.id = $1;
    `
	row := db.QueryRow(query, params["id"])

	var orderDetail OrderDetail
	err := row.Scan(&orderDetail.OrderID, &orderDetail.UserName, &orderDetail.UserEmail, &orderDetail.ProductName, &orderDetail.ProductPrice, &orderDetail.ProductCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(orderDetail)
}

func getAllOrders(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	query := `
    SELECT 
        orders.id AS order_id, 
        users.name AS user_name, 
        users.email AS user_email, 
        products.name AS product_name, 
        products.price AS product_price, 
        products.count AS product_count
    FROM orders
    JOIN users ON orders.user_id = users.id
    JOIN products ON orders.product_id = products.id;
    `
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var orderDetails []OrderDetail
	for rows.Next() {
		var orderDetail OrderDetail
		err := rows.Scan(&orderDetail.OrderID, &orderDetail.UserName, &orderDetail.UserEmail, &orderDetail.ProductName, &orderDetail.ProductPrice, &orderDetail.ProductCount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		orderDetails = append(orderDetails, orderDetail)
	}

	json.NewEncoder(w).Encode(orderDetails)
}

func main() {
	initDB()
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/order", createOrder).Methods("POST")
	router.HandleFunc("/order/{id}", getOrder).Methods("GET")
	router.HandleFunc("/orders", getAllOrders).Methods("GET")
	log.Fatal(http.ListenAndServe(":8002", router))
}
