package main

import (
	"log"
	"my-go-api/config"
	"my-go-api/handler"
	"my-go-api/middleware"
	"my-go-api/repository"
	"my-go-api/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// 1. Database
	db := config.ConnectDB()

	// 2. Dependency Injection
	productRepo := repository.NewProductRepository(db)
	productServ := service.NewProductService(productRepo)
	productHand := handler.NewProductHandler(productServ)

	// Dependency Injection Auth
	userRepo := repository.NewUserRepository(db)
	authServ := service.NewAuthService(userRepo)
	authHand := handler.NewAuthHandler(authServ)

	// 3. Express/Fiber Routes
	app := fiber.New()
	api := app.Group("/api/v1")

	// Public Routes (Auth)
	api.Post("/auth/register", authHand.Register)
	api.Post("/auth/login", authHand.Login)

	// Public Routes (Read Product)
	api.Get("/products", productHand.GetProducts)
	api.Get("/products/:id", productHand.GetProductByID)

	// Protected Routes (Butuh Header "Authorization: Bearer <token>")
	protected := api.Group("/", middleware.Protected())
	protected.Post("/products", productHand.CreateProduct)
	protected.Put("/products/:id", productHand.UpdateProduct)
	protected.Delete("/products/:id", productHand.DeleteProduct)

	// Dependency Injection Category
	categoryRepo := repository.NewCategoryRepository(db)
	categoryServ := service.NewCategoryService(categoryRepo)
	categoryHand := handler.NewCategoryHandler(categoryServ)

	// Routes Category (Protected)
	protected.Post("/categories", categoryHand.CreateCategory)
	protected.Get("/categories", categoryHand.GetCategories)
	// 4. Start Server
	log.Fatal(app.Listen(":8080"))
}
