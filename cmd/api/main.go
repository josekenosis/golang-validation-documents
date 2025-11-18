package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/neoway/golang-validation-documents/internal/config"
	"github.com/neoway/golang-validation-documents/internal/database"
	"github.com/neoway/golang-validation-documents/internal/handler"
	"github.com/neoway/golang-validation-documents/internal/middleware"
	"github.com/neoway/golang-validation-documents/internal/repository"
	"github.com/neoway/golang-validation-documents/internal/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	if err := config.Load(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	if err := database.Connect(); err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	app.Use(middleware.RequestCounter())
	app.Use(middleware.SanitizeInput())

	authHandler := handler.NewAuthHandler()
	statusHandler := handler.NewStatusHandler()

	app.Post("/login", authHandler.Login)
	app.Get("/status", statusHandler.GetStatus)

	documentRepo := repository.NewDocumentRepository(database.DB)
	documentService := service.NewDocumentService(documentRepo)
	documentHandler := handler.NewDocumentHandler(documentService)

	api := app.Group("/api/v1")
	api.Use(middleware.AuthMiddleware())

	api.Post("/documents", documentHandler.Create)
	api.Get("/documents", documentHandler.List)
	api.Get("/documents/search", documentHandler.GetByNumber)
	api.Get("/documents/:id", documentHandler.GetByID)
	api.Put("/documents/:id", documentHandler.Update)
	api.Delete("/documents/:id", documentHandler.Delete)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Info().Msg("Shutting down server...")
		_ = app.Shutdown()
	}()

	port := ":" + config.AppConfig.Port
	log.Info().Str("port", port).Msg("Server starting")
	if err := app.Listen(port); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

