package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	instance *sql.DB
}

type Reminder struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	fmt.Println("\033[32m- Successful connection to database\033[0m")

	return &DB{instance: db}, nil
}

func (d *DB) checkInitialConditions() {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS reminders (
        id INT AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(100) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `

	_, err := d.instance.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
		return
	}

	fmt.Println("\033[32m- Initial conditions checked\033[0m")
}

func (d *DB) GetInstance() *sql.DB {
	return d.instance
}

func (d *DB) Close() error {
	return d.instance.Close()
}
