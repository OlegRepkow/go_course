package main

import (
	"log"
	"module_6/cmd/server/handlers"
	"module_6/cmd/server/middlewares"
	_ "module_6/docs"
	"module_6/internal/clients"
	"module_6/internal/config"
	"module_6/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	swagger "github.com/swaggo/fiber-swagger"
)

// @title Chat API
// @version 1.0
// @description WebSocket chat server with authentication
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	log.Println("[INFO] Starting server...")
	cfg := config.Load()
	log.Printf("[INFO] Config loaded: port=%s, db=%s\n", cfg.Port, cfg.DBName)

	mongoClient, err := clients.NewMongo(cfg.MongoURI)
	if err != nil {
		log.Fatalf("[FATAL] Failed to connect to MongoDB: %v\n", err)
	}
	log.Println("[INFO] Connected to MongoDB")
	db := mongoClient.Database(cfg.DBName)

	authService := services.NewAuthService(db, cfg.JWTSecret)
	chatService := services.NewChatService(db)

	authHandler := handlers.NewAuthHandler(authService)
	chatHandler := handlers.NewChatHandler(chatService)

	app := fiber.New()

	app.Get("/swagger/*", swagger.WrapHandler)

	app.Post("/auth/sign-up", authHandler.SignUp)
	app.Post("/auth/sign-in", authHandler.SignIn)

	protected := app.Group("", middlewares.Auth(cfg.JWTSecret))
	protected.Get("/channel/history", chatHandler.History)
	protected.Post("/channel/send", chatHandler.Send)

	app.Use("/channel/listen", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/channel/listen", websocket.New(chatHandler.Listen))

	log.Printf("[INFO] Server listening on port %s\n", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
