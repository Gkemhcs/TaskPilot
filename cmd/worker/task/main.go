package main

import (
	"context"
	"database/sql"
	"log"

	"strconv"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"

	"github.com/Gkemhcs/taskpilot/internal/exporter"
	exporterdb "github.com/Gkemhcs/taskpilot/internal/exporter/gen"
	"github.com/Gkemhcs/taskpilot/internal/importer"
	importerdb "github.com/Gkemhcs/taskpilot/internal/importer/gen"
	"github.com/Gkemhcs/taskpilot/internal/storage"
	"github.com/Gkemhcs/taskpilot/internal/task"
	taskdb "github.com/Gkemhcs/taskpilot/internal/task/gen"
	"github.com/Gkemhcs/taskpilot/internal/user"
	userdb "github.com/Gkemhcs/taskpilot/internal/user/gen"
)

func main() {
	cfg := LoadWorkerConfig()

	// ---------- Logger ----------
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint:     true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// ---------- PostgreSQL ----------
	dbURL := "postgres://" + cfg.DBUser + ":" + cfg.DBPassword + "@" + cfg.DBHost + ":" + cfg.DBPort + "/" + cfg.DBName + "?sslmode=disable"
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Ping the database to ensure the connection is valid
	if err := db.Ping(); err != nil {
		logger.Fatal("Cannot ping DB: ", err)
	}

	// ---------- RabbitMQ ----------
	rabbitURL := cfg.RabbitMQURL
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open channel: %v", err)
	}
	defer ch.Close()

	queueName := cfg.TaskQueue
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ Failed to declare Import queue: %v", err)
	}
	exportQueueName := cfg.TaskExportQueue
	_, err = ch.QueueDeclare(exportQueueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ Failed to declare Export queue: %v", err)
	}

	// ---------- Storage ----------
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	storageClient, err := storage.StorageFactory(ctx, cfg.StorageType, cfg.StorageConfig)
	if err != nil {
		log.Fatalf("❌ Failed to initialize storage: %v", err)
	}

	// ---------- Dependencies ----------
	taskRepo := taskdb.New(db)
	taskService := task.NewTaskService(taskRepo)

	userRepo := userdb.New(db)
	userService := user.NewUserService(userRepo)

	importRepo := importerdb.New(db)
	exportRepo := exporterdb.New(db)

	expectedHeaders := []string{"project_id", "title", "assignee_email", "description", "status", "priority", "due_date"}

	rowHandler := func(data map[string]string, userID int) error {
		ctx := context.Background()

		user, err := userService.GetUserByEmail(ctx, data["assignee_email"])
		if err != nil {

			return err
		}

		projectID, err := strconv.Atoi(data["project_id"])
		if err != nil {
			return err
		}
		dueDate, err := time.Parse(time.RFC3339, data["due_date"])
		if err != nil {
			return err
		}

		taskInput := task.CreateTaskInput{
			ProjectID:   projectID,
			Title:       data["title"],
			AssigneeID:  int(user.ID),
			Description: data["description"],
			Status:      data["status"],
			Priority:    data["priority"],
			DueDate:     dueDate,
		}

		_, err = taskService.CreateTask(ctx, taskInput)
		if err != nil {
			return err
		}

		return nil
	}

	excelImporter := importer.NewExcelImporter(expectedHeaders, rowHandler)
	sheetName := "task"
	excelExporter := exporter.NewExcelExporter(expectedHeaders,sheetName)
	
	localDir := cfg.StorageConfig.ProcessDir
	worker := NewTaskWorker(
		excelImporter,
		excelExporter,
		storageClient,
		taskService,
		importRepo,
		exportRepo,
		logger,
		expectedHeaders,
		sheetName,
		localDir,
	)
	go func() {
		logger.Println("🚀 Starting Task Import Worker...")
		if err := worker.StartConsuming(ch, queueName); err != nil {
			logger.Fatalf("❌ Import Worker failed: %v", err)
		}
	}()

	go func() {
		logger.Println("🚀 Starting Task Export Worker...")
		if err := worker.StartConsumingExport(ch, exportQueueName); err != nil {
			logger.Fatalf("❌ Export Worker failed: %v", err)
		}
	}()

	// block main thread forever or until termination
	select {}
}
