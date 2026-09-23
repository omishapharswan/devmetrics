package main

import (
	"log"
	"os"

	"github.com/omishapharswan/devmetrics/backend/internal/api"
	"github.com/omishapharswan/devmetrics/backend/internal/db"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DEVMETRICS_DB_PATH")
	if dbPath == "" {
		dbPath = "data/devmetrics.db"
	}

	conn, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer conn.Close()

	router := api.NewRouter(conn)

	log.Printf("devmetrics backend listening on :%s (db: %s)", port, dbPath)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
