package storage

import (
	"TelegramShop/models"
	"database/sql"
	"log"
	"os"

	_ "github.com/glebarez/sqlite"
)

var db *sql.DB

func OpenSQLite() error {
	if err := os.MkdirAll("./data", 0755); err != nil {
		return err
	}

	var err error
	db, err = sql.Open("sqlite", "./data/database.db")
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
    user_id INTEGER UNIQUE NOT NULL,
    username TEXT,
	firstname TEXT NOT NULL,
	lastname TEXT,
    balance INTEGER DEFAULT 0,
    language_code TEXT DEFAULT 'ru',
    role TEXT DEFAULT 'user',
	state TEXT DEFAULT 'nothing',
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
	)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS products (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
    category_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    price INTEGER NOT NULL,
    image_id TEXT,
    stock INTEGER DEFAULT 0,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
	)`)

	if err != nil {
		return err
	}

	log.Println("Подключение к SQLite успешно!")
	return nil
}

func AddUser(user_id int64, username string, firstname string, lastname string, lang_code string) error {
	query := `INSERT INTO users (user_id, username, firstname, lastname, language_code) VALUES (?, ?, ?, ?, ?)`

	_, err := db.Exec(query, user_id, username, firstname, lastname, lang_code)
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func FindUser(userid int64) (*models.User, error) {
	query := `SELECT id, user_id, username, firstname, lastname, balance, language_code, role, state, updated_at, created_at FROM users WHERE user_id = ? LIMIT 1`

	var user models.User
	err := db.QueryRow(query, userid).Scan(
		&user.ID,
		&user.UserID,
		&user.Username,
		&user.Firstname,
		&user.Lastname,
		&user.Balance,
		&user.LangCode,
		&user.Role,
		&user.State,
		&user.UpdatedAt,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func RefreshUser(userid int64, username string, firstname string, lastname string, lang_code string) (*models.User, error) {
	query := `UPDATE users SET username = ?, firstname = ?, lastname = ?, language_code WHERE user_id = ?;`

	_, err := db.Exec(query, username, firstname, lastname, lang_code, userid)
	if err != nil {
		return nil, err
	}

	return FindUser(userid)
}
