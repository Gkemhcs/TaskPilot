package db

import (
	"database/sql"
	"fmt"

	"github.com/Gkemhcs/taskpilot/internal/config"
	userdb "github.com/Gkemhcs/taskpilot/internal/user/gen"

	"github.com/sirupsen/logrus"
)

func InitDB(logger *logrus.Logger,config *config.Config )(*userdb.Queries){
	dbURL:=fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",config.DBUser,config.DBPassword,config.DBHost,config.DBPort,config.DBName)

	


	

	 conn, err := sql.Open("postgres", dbURL)
    if err != nil {
        logger.Fatal("Cannot open DB: ", err)
    }

	if err := conn.Ping(); err != nil {
        logger.Fatal("Cannot ping DB: ", err)
    }
	db := userdb.New(conn)
	return db
}


