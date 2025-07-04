package main

import (
	"fmt"

	"github.com/Gkemhcs/taskpilot/cmd/server"
	"github.com/Gkemhcs/taskpilot/internal/config"
	"github.com/Gkemhcs/taskpilot/internal/db"

	"github.com/Gkemhcs/taskpilot/internal/utils"
	_ "github.com/lib/pq" // PostgreSQL driver for database/sql
)

// main is the entry point of the application. It initializes the logger, loads configuration, sets up the database connection, and starts the HTTP server.
func main() {
	// Initialize a structured logger for the application
	logger := utils.NewLogger()

	// Load configuration from environment variables or .env file
	config := config.LoadConfig()

	// Initialize the database connection using the loaded config and logger
	// Returns a userdb.Queries instance for database operations
	dbConn := db.InitDB(logger, config)

	// Start the HTTP server with the provided config, logger, and database connection
	err := server.NewServer(config, logger, dbConn)

	if err != nil {
		// Panic if the server fails to start
		panic(err)
	} else {
		fmt.Println("server running fine")
	}
}
