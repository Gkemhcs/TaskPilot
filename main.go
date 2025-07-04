package main

import (
	"fmt"

	"github.com/Gkemhcs/taskpilot/cmd/server"
	"github.com/Gkemhcs/taskpilot/internal/config"
	"github.com/Gkemhcs/taskpilot/internal/db"

	"github.com/Gkemhcs/taskpilot/internal/utils"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {

	logger := utils.NewLogger()
	config := config.LoadConfig()
	dbConn := db.InitDB(logger, config)

	err := server.NewServer(config, logger, dbConn)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("server running fine")
	}
}
