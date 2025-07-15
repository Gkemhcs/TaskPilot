package main

import (
	"context"
	"database/sql"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"github.com/Gkemhcs/taskpilot/internal/exporter"
	exporterdb "github.com/Gkemhcs/taskpilot/internal/exporter/gen"
	"github.com/Gkemhcs/taskpilot/internal/importer"
	importerdb "github.com/Gkemhcs/taskpilot/internal/importer/gen"
	"github.com/Gkemhcs/taskpilot/internal/project"
	projectdb "github.com/Gkemhcs/taskpilot/internal/project/gen"
	"github.com/Gkemhcs/taskpilot/internal/storage"
)

// main is the entry point for the worker process.
// It sets up configuration, logging, database, RabbitMQ, storage, and starts import/export workers.
func main() {
	// Load configuration from environment and .env
	cfg := LoadWorkerConfig()

	// Set up structured logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Connect to PostgreSQL database
	dbURL := "postgres://" + cfg.DBUser + ":" + cfg.DBPassword + "@" + cfg.DBHost + ":" + cfg.DBPort + "/" + cfg.DBName + "?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Ping the database to ensure the connection is valid
	if err := db.Ping(); err != nil {
		logger.Fatal("Cannot ping DB: ", err)
	}

	// Connect to RabbitMQ
	rabbitURL := cfg.RabbitMQURL
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		logger.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Open RabbitMQ channel
	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf("❌ Failed to open channel: %v", err)
	}
	defer ch.Close()

	// Declare import and export queues
	queueName := cfg.ProjectQueue
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		logger.Fatalf("❌ Failed to declare Import queue: %v", err)
	}
	exportQueueName := cfg.ProjectExportQueue
	_, err = ch.QueueDeclare(exportQueueName, true, false, false, false, nil)
	if err != nil {
		logger.Fatalf("❌ Failed to declare Export queue: %v", err)
	}

	// Initialize storage client (local or GCP)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	storageClient, err := storage.StorageFactory(ctx, cfg.StorageType, cfg.StorageConfig)
	if err != nil {
		logger.Fatalf("❌ Failed to initialize storage: %v", err)
	}

	// Set up dependencies for project import/export
	projectRepo := projectdb.New(db)
	projectService := project.NewProjectService(projectRepo)
	importRepo := importerdb.New(db)
	expectedHeaders := []string{"name", "description", "color"}

	// Handler for each row in the imported Excel file
	rowHandler := func(data map[string]string, userID int) error {
		ctx := context.Background()
		project := project.Project{
			Name:        data["name"],
			Description: data["description"],
			Color:       data["color"],
			User:        int(userID), // set properly below
		}
		_, err := projectService.CreateProject(ctx, project)
		if err != nil {
			return err
		}
		return nil
	}

	// Create Excel importer and exporter
	excelImporter := importer.NewExcelImporter(expectedHeaders, rowHandler)
	sheetName := "projects"
	excelExporter := exporter.NewExcelExporter(expectedHeaders, sheetName)
	exportRepo := exporterdb.New(db)
	localDir := cfg.StorageConfig.ProcessDir

	// Construct the worker with all dependencies
	worker := NewProjectWorker(
		excelImporter,
		excelExporter,
		storageClient,
		*projectService,
		importRepo,
		exportRepo,
		logger,
		expectedHeaders,
		sheetName,
		localDir,
	)

	// Start import worker goroutine
	go func() {
		logger.Println("🚀 Starting Project Import Worker...")
		if err := worker.StartConsuming(ch, queueName); err != nil {
			logger.Fatalf("❌ Import Worker failed: %v", err)
		}
	}()

	// Start export worker goroutine
	go func() {
		logger.Println("🚀 Starting Project Export Worker...")
		if err := worker.StartConsumingExport(ch, exportQueueName); err != nil {
			logger.Fatalf("❌ Export Worker failed: %v", err)
		}
	}()

	// Block main thread forever or until termination
	select {}
}
