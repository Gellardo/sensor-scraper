package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize the web server
	r := gin.Default()

	// Define a route to display the line chart
	r.GET("/", handleDataRequest)

	// Run the web server on a specific port
	r.Run(":8080")
}

func handleDataRequest(c *gin.Context) {
	// Open the SQLite database file
	db, err := sql.Open("sqlite3", "your-database-file.db")
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}
	defer db.Close()

	// Query the database to fetch the tabular data
	rows, err := db.Query("SELECT timestamp, value FROM your_table_name")
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

	// Render an HTML page with the retrieved data
	c.HTML(http.StatusOK, "chart.html", gin.H{
		"data": data,
	})
}
