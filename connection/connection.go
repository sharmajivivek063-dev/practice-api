package connection

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Global database connection
var DB *sql.DB

// Connect to MySQL database
func ConnectDB() {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}

	// Check database connection
	if err := db.Ping(); err != nil {
		db.Close()
		panic(err)
	}

	DB = db

	fmt.Println("MySQL Connected Successfully!")
}

// Close database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
