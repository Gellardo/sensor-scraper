package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
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
            sensorid INTEGER,
            timestamp INTEGER,
            value REAL
        );
    `
)

func main() {
	r := gin.Default()
	r.Use(gin.Logger())

	// Initialize the database
	if err := initializeDatabase(); err != nil {
		log.Printf("Error initializing database: %v\n", err)
		return
	}

	setupScraper()

	SetupTemplatesAndStatic(r)
	r.GET("/", handleHomeRequest)
	r.GET("/graph", handleDataRequest)
	r.GET("/ping", func(c *gin.Context) { c.String(http.StatusOK, "{\"pong\":{\"ping\":1.0}}") })
	r.GET("/scrape", triggerScrape)

	r.Run("localhost:8080") // TODO make configurable
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
		log.Println("Database file and table created.")

		// Insert dummy values for yesterday and today if they don't exist
		yesterday := time.Now().Add(-24 * time.Hour).Unix()
		today := time.Now().Unix()

		err = insertDummyValues(db, 1, yesterday, 1)
		if err != nil {
			return err
		}
		err = insertDummyValues(db, 1, today, 4)
		if err != nil {
			return err
		}
		log.Println("Inserted dummy values")
	}
	return nil
}

func insertDummyValues(db *sql.DB, sensorid int64, timestamp int64, value float64) error {
	insertSQL := "INSERT INTO " + tableName + " (sensorid, timestamp, value) VALUES (?, ?, ?)"
	_, err := db.Exec(insertSQL, sensorid, timestamp, value)
	return err
}

func handleHomeRequest(c *gin.Context) {
	config, err := loadSensorConfig("config.toml")
	if err != nil {
		log.Printf("Unable to load config: %v\n", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}

	// Render the homepage with links to graph pages
	c.HTML(http.StatusOK, "home.html", gin.H{"sensors": config.Sensors})
}
func handleDataRequest(c *gin.Context) {
	// Open the SQLite database file
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}
	defer db.Close()

	sensorid := c.DefaultQuery("sensorid", "1")
	daysBack, _ := strconv.Atoi(c.DefaultQuery("range", "14"))
	daysBackStr := fmt.Sprintf("-%d days", daysBack)
	// Query the database to fetch the tabular data
	rows, err := db.Query("SELECT timestamp, value FROM "+tableName+" where sensorid = ? and timestamp >= strftime('%s', 'now', '"+daysBackStr+"') ORDER BY timestamp", sensorid)
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
		//timeStr := time.Unix(timestamp, 0).Format(time.RFC3339)

		// Store the data in a map
		dataPoint := map[string]interface{}{
			"timestamp": timestamp,
			"value":     value,
		}

		data = append(data, dataPoint)
	}
	log.Printf("Fetched %d data points\n", len(data))

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}

	// Render an HTML page with the retrieved JSON data
	c.HTML(http.StatusOK, "chart.html", gin.H{
		"jsonData": template.HTML(jsonData), // Use template.HTML to render as raw HTML
		"sensorid": sensorid,
	})
}
