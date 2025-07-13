// @title           TaskPilot API
// @version         1.0
// @description     Task and Project Management Backend built with Go, Gin, and PostgreSQL.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Koti Eswar Mani Gudi
// @contact.email  gudikotieswarmani@gmail.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer <your-token>" to authorize
package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Gkemhcs/taskpilot/docs"
	_ "github.com/Gkemhcs/taskpilot/docs"
	"github.com/Gkemhcs/taskpilot/internal/auth"
	"github.com/Gkemhcs/taskpilot/internal/config"
	"github.com/Gkemhcs/taskpilot/internal/exporter"
	exporterdb "github.com/Gkemhcs/taskpilot/internal/exporter/gen"
	"github.com/Gkemhcs/taskpilot/internal/importer"
	importerdb "github.com/Gkemhcs/taskpilot/internal/importer/gen"
	"github.com/Gkemhcs/taskpilot/internal/middleware"
	"github.com/Gkemhcs/taskpilot/internal/project"
	projectdb "github.com/Gkemhcs/taskpilot/internal/project/gen"
	"github.com/Gkemhcs/taskpilot/internal/storage"
	"github.com/Gkemhcs/taskpilot/internal/task"
	taskdb "github.com/Gkemhcs/taskpilot/internal/task/gen"
	"github.com/Gkemhcs/taskpilot/internal/user"
	userdb "github.com/Gkemhcs/taskpilot/internal/user/gen"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewServer sets up the Gin router, registers routes, and starts the HTTP server.
// It takes configuration, logger, and database connection as input.
func NewServer(config *config.Config, logger *logrus.Logger, dbConn *sql.DB) error {

	
	// Set Gin to release mode for production
	gin.SetMode(gin.ReleaseMode)
	// Create a new Gin router with default middleware (logger, recovery)
	router := gin.Default()

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", config.HOST, config.Port)
	

	redisAddr := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)

	
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		logger.Fatal("❌ Redis connection failed:", err)
	} else {
		logger.Info("✅ Redis connected:", pong)
	}

	//Initialising ratelimiter middleware with Redis client and logger
	rateLimiterMiddleware, err := middleware.RateLimiterMiddleware(redisClient, logger)
	if err != nil {
		panic(fmt.Errorf("failed to create rate limiter middleware: %w", err))
	}else{
		logger.Info("Rate limiter middleware initialized successfully")
	}

	

	// Create API v1 group with custom logger middleware
	v1 := router.Group("/api/v1",  middleware.LoggerMiddleware(logger), middleware.PrometheusMiddleware())
	v1.Use(rateLimiterMiddleware)
	// Expose Prometheus metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// swagger docs route setup
	router.GET("/api/v1/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))



	// Initialize user service with database connection
	userService := user.NewUserService(userdb.New(dbConn))

	// Initialize project service with database connection
	projectService := project.NewProjectService(projectdb.New(dbConn))

	// Initialize task service with database connection
	taskService := task.NewTaskService(taskdb.New(dbConn))

	// Initialize JWT manager for authentication
	params := auth.CreateJwtManagerParams{
		AccessTokenDuration:  config.AccessTokenDuration,
		RefreshTokenDuration: config.RefreshTokenDuration,
		AccessTokenKey:       config.JWTAccessTokenSecret,
		RefreshTokenKey:      config.JWTRefreshTokenSecret,
	}
	jwtManager := auth.NewJWTManager(params)

	// Create user handler with service, logger, and JWT manager
	userHandler := user.NewUserHandler(userService, logger, jwtManager)
	
	// Register user-related routes under /api/v1/users
	user.RegisterRoutes(v1, userHandler)

	// Create project handler with service, logger
	projectHandler := project.NewProjectHandler(logger, projectService, taskService)

	// Register project-related routes under /api/v1/projects
	project.RegisterProjectRoutes(v1, projectHandler, jwtManager)

	// Create task handler with service, logger
	taskHandler := task.NewTaskHandler(*taskService, userService, logger, projectService)

	// Register task-related routes under /api/v1/tasks
	task.RegisterTaskRoutes(v1, taskHandler, jwtManager)

	// Initiialize importhandler and service 
	ctx ,cancel:=context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	storageClient,err:=storage.StorageFactory(ctx, config.StorageType, config.StorageConfig)
	importerRepo:=importerdb.New(dbConn)
	exporterRepo:=exporterdb.New(dbConn)
	// Step 1: Connect to RabbitMQ
	conn, err := amqp091.Dial(config.RabbitMQURL)
	if err != nil {
		logger.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Step 2: Open a channel
	ch, err := conn.Channel()
	if err != nil {
		logger.Fatalf("❌ Failed to open a channel: %v", err)
	}
	defer ch.Close()


	// Step 3: Declare the queue for project imports
	err = declareQueue(ch, config.ProjectPublisher.QueueName)
	if err != nil {
		logger.Fatalf("❌ Failed to declare project queue: %v", err)
	}


	projectPublisher:=importer.NewRabbitMQPublisher(ch, config.ProjectPublisher.QueueName, config.ProjectPublisher.Exchange,config.ProjectPublisher.RoutingKey,)
	taskPublisher:=importer.NewRabbitMQPublisher(ch, config.TaskPublisher.QueueName, config.TaskPublisher.Exchange,config.TaskPublisher.RoutingKey,)
	
	importService := importer.NewImportService(storageClient, importerRepo, projectPublisher, taskPublisher, logger)
	importHandler:=importer.NewImportHandler(importService, logger)
	importer.RegisterImporterHandler(importHandler, v1, jwtManager)


	projectExportPublisher:=exporter.NewRabbitMQPublisher(ch, config.ProjectExportPublisher.QueueName, config.ProjectExportPublisher.Exchange,config.ProjectExportPublisher.RoutingKey,)
	taskExportPublisher:=exporter.NewRabbitMQPublisher(ch, config.TaskExportPublisher.QueueName, config.TaskExportPublisher.Exchange,config.TaskExportPublisher.RoutingKey,)
	
	exportService := exporter.NewExportService(exporterRepo, projectExportPublisher, taskExportPublisher, logger)
	exportHandler:=exporter.NewExportHandler(exportService, logger)
	exporter.RegisterExportHandler(exportHandler, v1, jwtManager)






	// Health check endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	// Start the HTTP server on the configured host and port
	return router.Run(fmt.Sprintf("%s:%s", config.HOST, config.Port))
}



func declareQueue(ch *amqp091.Channel, name string) error {
	_, err := ch.QueueDeclare(
		name,  // queue name
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	return err
}