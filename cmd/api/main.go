package main

import (
	"fmt"
	"log"

	"healthcare-app/config"
	"healthcare-app/internal/handlers"
	"healthcare-app/internal/repositories"
	"healthcare-app/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "healthcare-app/docs" // Generated Swagger docs
)

// @title           Healthcare Management API
// @version         1.0
// @description     API for healthcare management with receptionist and doctor portals
// @host            localhost:8080
// @BasePath        /api/v1

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	patientRepo := repositories.NewPatientRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	userService := services.NewUserService(userRepo)
	patientService := services.NewPatientService(patientRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	patientHandler := handlers.NewPatientHandler(patientService)

	// Set up the router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		v1.POST("/login", authHandler.Login)

		// User routes
		v1.POST("/users", authHandler.RequireAuth(userHandler.CreateUser))
		
		// Patient routes - Receptionist access
		receptionistRoutes := v1.Group("/patients")
		receptionistRoutes.Use(authHandler.RequireAuth(authHandler.RequireReceptionist))
		{
			receptionistRoutes.POST("", patientHandler.CreatePatient)
			receptionistRoutes.GET("", patientHandler.GetAllPatients)
			receptionistRoutes.GET("/:id", patientHandler.GetPatient)
			receptionistRoutes.PUT("/:id", patientHandler.UpdatePatient)
			receptionistRoutes.DELETE("/:id", patientHandler.DeletePatient)
		}

		// Patient routes - Doctor access
		doctorRoutes := v1.Group("/doctor/patients")
		doctorRoutes.Use(authHandler.RequireAuth(authHandler.RequireDoctor))
		{
			doctorRoutes.GET("", patientHandler.GetAllPatients)
			doctorRoutes.GET("/:id", patientHandler.GetPatient)
			doctorRoutes.PUT("/:id/medical", patientHandler.UpdatePatientMedicalInfo)
		}
	}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 