package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/yourusername/draft-forge/internal/agents"
	apiHandlers "github.com/yourusername/draft-forge/internal/api"
	"github.com/yourusername/draft-forge/internal/auth"
	"github.com/yourusername/draft-forge/internal/db"
	dbagent "github.com/yourusername/draft-forge/internal/db/agent"
	dbproject "github.com/yourusername/draft-forge/internal/db/project"
	"github.com/yourusername/draft-forge/internal/projects"
	"github.com/yourusername/draft-forge/internal/scaffold"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(databaseURL) == "" {
		log.Fatal("DATABASE_URL is required")
	}

	// Connect to database
	dbConn, err := db.Connect(databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer dbConn.Close()

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "draftforge-api",
		})
	})

	// API routes (to be implemented)
	api := app.Group("/api/v1")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "DraftForge API v1",
		})
	})

	// Auth setup
	jwtSecret := os.Getenv("JWT_SECRET")
	refreshSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	if jwtSecret == "" || refreshSecret == "" {
		log.Fatal("JWT_SECRET and REFRESH_TOKEN_SECRET are required")
	}
	tokenManager := auth.NewTokenManager(jwtSecret, refreshSecret, 15*time.Minute, 30*24*time.Hour)

	githubClientID := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	githubRedirectURI := os.Getenv("GITHUB_REDIRECT_URI")
	if githubClientID == "" || githubClientSecret == "" {
		log.Fatal("GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET are required")
	}
	stateSecret := os.Getenv("STATE_SECRET")
	if stateSecret == "" {
		stateSecret = jwtSecret
	}

	ghClient := auth.NewOAuthClient(nil, githubClientID, githubClientSecret, githubRedirectURI)
	userStore := auth.NewSQLUserStore(dbConn)
	authService := auth.NewService(userStore, ghClient, tokenManager, githubClientID, githubRedirectURI, stateSecret)
	authHandler := apiHandlers.NewAuthHandler(authService, tokenManager, userStore)
	authHandler.Register(api)

	artifactDir := os.Getenv("AGENT_ARTIFACT_DIR")
	if artifactDir == "" {
		artifactDir = ".draftforge/agent-runs"
	}

	sqlxDB := sqlx.NewDb(dbConn, "postgres")

	agentStore := dbagent.NewStore(sqlxDB)
	agentService := agents.NewService(agentStore, artifactDir)
	agentHandler := apiHandlers.NewAgentHandler(agentService)
	protected := api.Group("", apiHandlers.AuthMiddleware(tokenManager, userStore))

	projectStore := dbproject.NewStore(sqlxDB)
	scaffoldRoot := os.Getenv("SCAFFOLD_ROOT")
	if scaffoldRoot == "" {
		scaffoldRoot = "scaffolds"
	}
	localScaffolder := scaffold.NewLocalScaffolder(scaffoldRoot)
	githubScaffolder := scaffold.NewGitHubScaffolder(nil)
	projectScaffolder := &scaffold.CompositeScaffolder{
		Remote: githubScaffolder,
		Local:  localScaffolder,
	}
	projectService := projects.NewService(projectStore, projectScaffolder)
	projectHandler := apiHandlers.NewProjectHandler(projectService)
	projectHandler.Register(protected)

	agentHandler.Register(protected)

	// Start server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting DraftForge API on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": message,
	})
}
