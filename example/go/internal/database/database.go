package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "orders.db")
	if err != nil {
		return nil, err
	}

	// 创建表
	err = createTables(db)
	if err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}

func createTables(db *sql.DB) error {
	// Create orders table
	createOrdersTable := `
	CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		order_number TEXT UNIQUE NOT NULL,
		customer_name TEXT NOT NULL,
		customer_email TEXT NOT NULL,
		product_id TEXT NOT NULL,
		product_token_id TEXT NOT NULL,
		quantity INTEGER NOT NULL,
		total_amount TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		reddio_payment_id TEXT UNIQUE,
		reddio_pay_link TEXT,
		reddio_status TEXT,
		transaction_hash TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		paid_at DATETIME
	);`

	// Execute table creation SQL
	if _, err := db.Exec(createOrdersTable); err != nil {
		return err
	}

	return nil
}