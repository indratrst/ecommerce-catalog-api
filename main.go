package main

import (
	"log"

	"ecommerce-catalog-api/config"
	"ecommerce-catalog-api/handler"
	"ecommerce-catalog-api/middleware"
	"ecommerce-catalog-api/repository"
	"ecommerce-catalog-api/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Database Connection
	db := config.ConnectDB()

	// 2. Dependency Injection
	// --- Auth ---
	userRepo := repository.NewUserRepository(db)
	authServ := service.NewAuthService(userRepo)
	authHand := handler.NewAuthHandler(authServ)

	// --- Product ---
	productRepo := repository.NewProductRepository(db)
	productServ := service.NewProductService(productRepo)
	productHand := handler.NewProductHandler(productServ)

	// --- Category ---
	categoryRepo := repository.NewCategoryRepository(db)
	categoryServ := service.NewCategoryService(categoryRepo)
	categoryHand := handler.NewCategoryHandler(categoryServ)

	// 3. Initialize Fiber App
	app := fiber.New()

	// 4. Global Middlewares (WAJIB DITARUH SEBELUM ROUTE)
	app.Use(logger.New()) // Tambahan: Log HTTP Request ke terminal
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // Origin Next.js
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// 5. Routes Setup
	api := app.Group("/api/v1")

	// --- Public Routes ---
	// Auth
	api.Post("/auth/register", authHand.Register)
	api.Post("/auth/login", authHand.Login)

	// Products
	api.Get("/products", productHand.GetProducts)
	api.Get("/products/:id", productHand.GetProductByID)

	// Categories (Get Categories biasanya public agar bisa dibaca di katalog)
	api.Get("/categories", categoryHand.GetCategories)

	// --- Protected Routes (Perlu Token JWT) ---
	protected := api.Group("/", middleware.Protected())

	// Protected Products
	protected.Post("/products", productHand.CreateProduct)
	protected.Put("/products/:id", productHand.UpdateProduct)
	protected.Delete("/products/:id", productHand.DeleteProduct)

	// Protected Categories
	protected.Post("/categories", categoryHand.CreateCategory)

	// 6. Start Server
	log.Println("Server running on port :8080")
	log.Fatal(app.Listen(":8080"))
}
