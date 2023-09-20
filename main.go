package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dbName         = "db.sqlite"
	tableName      = "value_table"
	createTableSQL = `
        CREATE TABLE IF NOT EXISTS value_table (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            timestamp INTEGER,
            value REAL
        );
    `
)

func main() {
	r := gin.Default()

	// Initialize the database
	if err := initializeDatabase(); err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		return
	}

	r.GET("/", handleDataRequest)

	r.Run("localhost:8080")
}

func initializeDatabase() error {
	// Check if the database file exists
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		// If not, create the database file and the table
		db, err := sql.Open("sqlite3", dbName)
		if err != nil {
			return err
		}
		defer db.Close()

		_, err = db.Exec(createTableSQL)
		if err != nil {
			return err
		}
		fmt.Println("Database file and table created.")

		// Insert dummy values for yesterday and today if they don't exist
		yesterday := time.Now().Add(-24 * time.Hour).Unix()
		today := time.Now().Unix()

		err = insertDummyValues(db, yesterday, 1)
		if err != nil {
			return err
		}
		err = insertDummyValues(db, today, 4)
		if err != nil {
			return err
		}
	}
	return nil
}

func insertDummyValues(db *sql.DB, timestamp int64, value float64) error {
	insertSQL := "INSERT INTO " + tableName + " (timestamp, value) VALUES (?, ?)"
	_, err := db.Exec(insertSQL, timestamp, value)
	return err
}

func handleDataRequest(c *gin.Context) {
	// Open the SQLite database file
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}
	defer db.Close()

	// Query the database to fetch the tabular data
	rows, err := db.Query("SELECT timestamp, value FROM " + tableName)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}
	defer rows.Close()

	var data []map[string]interface{}

	// Iterate through the rows and convert the epoch timestamp to a human-readable format
	for rows.Next() {
		var timestamp int64
		var value float64
		err := rows.Scan(&timestamp, &value)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
			return
		}

		// Convert the timestamp to a human-readable format (e.g., RFC3339)
		timeStr := time.Unix(timestamp, 0).Format(time.RFC3339)

		// Store the data in a map
		dataPoint := map[string]interface{}{
			"timestamp": timeStr,
			"value":     value,
		}

		data = append(data, dataPoint)
	}
	fmt.Printf("Returning %d data points\n", len(data))

	// Render an HTML page with the retrieved data
	c.HTML(http.StatusOK, "chart.html", gin.H{
		"data": data,
	})
}
