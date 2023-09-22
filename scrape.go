package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
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
	url := "http://localhost:8080/ping"
	sensorid := 1
	jsonPath := "pong.ping"
	resp, err := http.Get(url)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
		return
	}

	value, err := extractJsonPath(body, jsonPath)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
	}

	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
	}
	defer db.Close()

	// Insert the value into the database
	_, err = db.Exec("INSERT INTO "+tableName+" (timestamp, value, sensorid) VALUES (?, ?, ?)",
		time.Now().Unix(), value, sensorid)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error: %v", err))
	}
	c.String(http.StatusOK, "")
}
