package database

import (
	qrProcesser "agm/internal/qrProcesser"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

// saveToDatabase saves JSON data to a PostgreSQL database.
func SaveToDatabase(jsonData []byte) {
	connStr := "user=yourusername dbname=yourdbname sslmode=disable"

	// Open a connection to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Error closing database connection: %v", err)
		}
	}(db) // Close the database connection when done

	// Unmarshal JSON data into QRData struct
	var qrData qrProcesser.QRData
	err = json.Unmarshal(jsonData, &qrData)
	if err != nil {
		log.Fatalf("Error decoding JSON: %v", err)
	}

	// Insert QR data content into the database
	_, err = db.Exec("INSERT INTO qr_codes (content) VALUES ($1)", qrData.Content)
	if err != nil {
		log.Fatalf("Error inserting data into the database: %v", err)
	}
}

// FetchJSONFromDB retrieves JSON data from a PostgreSQL database.
func FetchJSONFromDB(content string) ([]byte, error) {
	db, err := sql.Open("postgres", "user=yourusername dbname=yourdbname sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			return
		}
	}()

	var qrData qrProcesser.QRData
	err = db.QueryRow("SELECT content FROM qr_codes WHERE content = $1", content).Scan(&qrData.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from database: %v", err)
	}

	jsonData, err := json.Marshal(qrData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode JSON: %v", err)
	}

	return jsonData, nil
}
