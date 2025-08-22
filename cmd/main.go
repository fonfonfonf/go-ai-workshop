package main

import (
	"database/sql"
	"log"
	"os"
	"time"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver
	"github.com/gofiber/fiber/v2"
	"github.com/nattakan-n/ai-training-backend/internal/handlers"
	"github.com/nattakan-n/ai-training-backend/internal/routes"
	"github.com/nattakan-n/ai-training-backend/internal/storage/sqlite"
	"github.com/nattakan-n/ai-training-backend/internal/usecases"
)

func main() {
	app := fiber.New()

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	// SQLite connection for users
	dsn := os.Getenv("SQLITE_DSN")
	if dsn == "" {
		dsn = "file:aiworkshop.db?cache=shared&mode=rwc"
	}
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatalf("sqlite3 open error: %v", err)
	}
	userRepo, err := sqlite.NewUserSQLiteRepository(db)
	if err != nil {
		log.Fatalf("sqlite repo error: %v", err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "dev-secret"
	}
	authSvc := usecases.NewAuthService(userRepo, jwtSecret, time.Hour*24)
	authHandler := handlers.NewAuthHandler(authSvc)

	routes.Register(app, authHandler, jwtSecret)

	// Optional seed user (set SEED_EMAIL and SEED_PASSWORD)
	// if seedEmail := os.Getenv("SEED_EMAIL"); seedEmail != "" {
	// 	// log.Printf("seed user created: %s", seedEmail)
	// 	if seedPass := os.Getenv("SEED_PASSWORD"); seedPass != "" {
	// 		// log.Printf("seed pass created: %s", seedPass)
	// 		if _, ok, err := userRepo.GetByEmail("test@example.com"); err == nil && !ok {
	// 			if _, err := authSvc.Register("test@example.com", "secret123", os.Getenv("SEED_NAME")); err != nil {
	// 				log.Printf("seed user failed: %v", err)
	// 			} else {
	// 				log.Printf("seed user created: %s", seedEmail)
	// 			}
	// 		}
	// 	}
	// }

	if _, ok, err := userRepo.GetByEmail("test@example.com"); err == nil && !ok {
		if _, err := authSvc.Register("test@example.com", "secret123", os.Getenv("SEED_NAME")); err != nil {
			log.Printf("seed user failed: %v", err)
		} else {
			log.Printf("seed user created: %s", "test@example.com")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	log.Printf("Starting server on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

