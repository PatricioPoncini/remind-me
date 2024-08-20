package db

import (
	"database/sql"
	"fmt"
	"log"
	"remind_me/src/utils"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	instance *sql.DB
}

type Reminder struct {
	ID         int        `json:"id"`
	Title      string     `json:"title"`
	Expiration *time.Time `json:"expiration"`
	ChatID     int64      `json:"chat_id"`
	CreatedAt  time.Time  `json:"created_at"`
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	utils.SuccessLog("Successful connection to database")

	return &DB{instance: db}, nil
}

func (d *DB) CheckInitialConditions() {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS reminders (
        id INT AUTO_INCREMENT PRIMARY KEY,
        title VARCHAR(100) NOT NULL,
		expiration TIMESTAMP NOT NULL,
		chat_id BIGINT UNSIGNED NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );
    `

	_, err := d.instance.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
		return
	}

	utils.SuccessLog("Initial conditions checked")
}

func (d *DB) InsertNewReminder(text string, timestamp time.Time, chatId int64) (*Reminder, error) {
	query := "INSERT INTO reminders (title, expiration, chat_id) VALUES (?, ?, ?);"

	result, err := d.instance.Exec(query, text, timestamp, chatId)
	if err != nil {
		return nil, fmt.Errorf("error trying to insert new reminder: %w", err)
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert ID: %w", err)
	}

	return d.GetReminderById(int(lastInsertID))
}

func (d *DB) GetReminderById(id int) (*Reminder, error) {
	query := "SELECT id, title, chat_id FROM reminders WHERE id = ?;"
	var reminder Reminder

	err := d.instance.QueryRow(query, id).Scan(&reminder.ID, &reminder.Title, &reminder.ChatID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no reminder found with id %d", id)
		}
		return nil, fmt.Errorf("error trying to get reminder by id: %v", err)
	}

	return &reminder, nil
}
