package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func extractJsonPath(haystack []byte, path string) (float64, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(haystack, &data); err != nil {
		return math.NaN(), err
	}

	keys := strings.Split(path, ".")
	current := data
	for _, key := range keys {
		value, ok := current[key]
		if !ok {
			return math.NaN(), errors.New(fmt.Sprintf("key %s not found in json", key))
		}

		if next, isMap := value.(map[string]interface{}); isMap {
			current = next
		} else {
			f, err := strconv.ParseFloat(fmt.Sprintf("%f", value), 64)
			if err != nil {
				return math.NaN(), err
			}
			return f, nil
		}
	}
	return math.NaN(), errors.New("Final json node is a map")

}

func triggerScrape(c *gin.Context) {
	config, err := loadSensorConfig("config.toml")
	if err != nil {
		log.Printf("Unable to load config: %v\n", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}

	errorList := []error{}
	for _, sensor := range config.Sensors {
		if err = scrapeSensor(sensor); err != nil {
			errorList = append(errorList, err)
		}
	}
	if len(errorList) > 0 {
		log.Printf("Scraping resulted in %d errors\n", len(errorList))
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", errors.Join(errorList...)))
	}
	c.String(http.StatusOK, "")
}

func scrapeSensor(sensor Sensor) error {
	url := sensor.URL
	sensorid := sensor.ID
	jsonPath := sensor.JSONPath
	resp, err := http.Get(url)
	if err != nil {
		return errors.New(fmt.Sprintf("Error: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New(fmt.Sprintf("Error: %v", err))
	}

	value, err := extractJsonPath(body, jsonPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Error: %v", err))
	}

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return errors.New(fmt.Sprintf("Error: %v", err))
	}
	defer db.Close()

	// Insert the value into the database
	_, err = db.Exec("INSERT INTO "+tableName+" (timestamp, value, sensorid) VALUES (?, ?, ?)",
		time.Now().Unix(), value, sensorid)
	if err != nil {
		return errors.New(fmt.Sprintf("Error: %v", err))
	}
	log.Printf("Successfully scraped sensor %d\n", sensorid)
	return nil
}
