package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	db "github.com/weldonkipchirchir/rental_listing/db/sqlc"
	"github.com/weldonkipchirchir/rental_listing/middleware"
	"github.com/weldonkipchirchir/rental_listing/tasks"
)

type Server struct {
	router     *gin.Engine
	db         *sql.DB
	q          *db.Queries
	client     *asynq.Client
	redis      *redis.Client
	httpServer *http.Server
}

func NewServer() (*Server, error) {
	router := gin.Default()

	limiter := middleware.NewRateLimiter(10, 20)
	router.Use(limiter.Middleware())
	router.MaxMultipartMemory = 8 << 20

	router.Use(cors.New(
		cors.Config{
			AllowOrigins:     []string{"http://localhost:3000", "http://172.23.32.1:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}))

	connectionString := "postgresql://root:secret@localhost:5432/rental_listing?sslmode=disable"
	fmt.Println("Connection String:", connectionString) // Debugging line
	dbInstance, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := dbInstance.Ping(); err != nil {
		return nil, err
	}

	queries := db.New(dbInstance)
	server := &Server{
		db: dbInstance,
		q:  queries,
	}

	const redisAddr = "localhost:6379"

	// Initialize Asynq client
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	server.client = client

	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     redisAddr,
			Password: "",
			DB:       0,
		})
	server.redis = redisClient

	// Initialize task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeVerificationEmail, tasks.HandleVerificationEmailTask)
	mux.HandleFunc(tasks.TypeForgotPasswordEmail, tasks.HandleForgotPasswordEmailTask)

	// Run Asynq background worker

	go func() {
		log.Println("Starting Asynq server...")
		if err := asynq.NewServer(
			asynq.RedisClientOpt{Addr: redisAddr},
			asynq.Config{
				Concurrency: 10,
				Queues: map[string]int{
					"default": 1,
				},
			},
		).Run(mux); err != nil {
			log.Printf("Asynq server error: %v", err)
			panic(err)
		}
		log.Println("Asynq server started successfully")
	}()

	server.initAdminRoutes(router)
	server.initUserRoutes(router)
	server.initListingRoutes(router)
	server.initBookingRoutes(router)
	server.initFavoriteRoutes(router)
	server.initReviewRoutes(router)
	server.initNotificationRoutes(router)
	server.initVerifyRoutes(router)
	server.initPaymentRoutes(router)
	server.initStatsRoutes(router)

	server.router = router

	server.httpServer = &http.Server{
		Handler: router,
	}

	return server, nil
}

func (server *Server) Start(address string) error {
	server.httpServer.Addr = address
	return server.httpServer.ListenAndServe()
}

func (server *Server) Shutdown(ctx context.Context) error {
	return server.httpServer.Shutdown(ctx)
}
