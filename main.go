package main

import (
	// "fmt"
	"log"

	"github.com/Frezeh/Go-backend/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/api", handlers.Ping)
	app.Post("/api/signup", handlers.SignUp)
	app.Post("/api/login", handlers.Login)
	app.Patch("/api/deposit", handlers.Deposit)
	app.Patch("/api/transfer/:id", handlers.Transfer)
	app.Patch("/api/transferout/", handlers.TransferOut)
	app.Get("/api/balance", handlers.GetBalance)
	app.Get("/api/all", handlers.All)

	log.Fatal(app.Listen(":3000"))
}