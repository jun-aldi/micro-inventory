package app

import (
	"micro-inventory/user-service/configs"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	zerolog "github.com/rs/zerolog/log"
)

func RunServer() {
	cfg := configs.NewConfig()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			zerolog.Printf("Error: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		},
	})

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] $ip  ${status} - ${latency}  ${method} ${path}\n",
	}))

	container := BuildContainer()
	SetupRoutes(app, container)

	port := cfg.App.AppPort
	if port == "" {
		port = os.Getenv("APP_PORT")
		if port == "" {
			log.Fatalf("server port is not set")
		}
	}

	zerolog.Printf("Server running on port %s", port)

	go func() {
		if err := app.Listen(":" + port); err != nil {
			zerolog.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	zerolog.Printf("Server is shutting down...")

	if err := app.Shutdown(); err != nil {
		zerolog.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	zerolog.Printf("Server exited properly")

}
