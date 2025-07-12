package main

import (
	"context"
	"log"
	"os"
	"sync"
	

	"github.com/Gkemhcs/taskpilot/internal/db"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func main() {
    ctx := context.Background()

    // Load configuration
    cfg := config.Load() // you can create a Load function from envs or .env

	logger := logrus.New()
    // Setup DB
    dbpool, err :=db.InitDB(logger,config)
    if err != nil {
        log.Fatal("failed to connect DB:", err)
    }
    queries := sqlc.New(dbpool)

    // Setup RabbitMQ
    conn, err := amqp091.Dial(cfg.RabbitMQ.URL)
    if err != nil {
        log.Fatal("failed to connect RabbitMQ:", err)
    }
    ch, err := conn.Channel()
    if err != nil {
        log.Fatal("channel err:", err)
    }

    // Setup Storage
    storageClient, err := storage.StorageFactory(ctx, cfg.Storage.Type, cfg.Storage)
    if err != nil {
        log.Fatal("storage init error:", err)
    }

    // Setup Project Service
    projectSvc := project.NewService(queries)

    // Prepare ExcelImporter
    rowHandler := func(row map[string]string) error {
        input := project.CreateProjectInput{
            Name:        row["project_name"],
            Description: row["description"],
            Color:       row["color"],
        }
        return projectSvc.Create(ctx, input)
    }

    importerInstance := importer.NewExcelImporter(
        []string{"project_name", "description", "color"},
        rowHandler,
    )

    // Setup worker
    worker := importer.NewWorker(
        "import_jobs",
        ch,
        importerInstance,
        nil, // headers handled inside
        storageClient,
        queries,
    )

    log.Println("🚀 Starting project worker...")
    if err := worker.Start(ctx); err != nil {
        log.Fatal("worker crashed:", err)
    }
}
