package main

import (
	"log"
	"os"
	db_drivers "ubm-canteen/db/db_drivers"
	"ubm-canteen/handlers"
	"ubm-canteen/internal/usecase"
	"ubm-canteen/routes"
	"ubm-canteen/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

	func main() {

		viper.AutomaticEnv()
		utils.InitRedis()
		
		err:= usecase.InitTextractClient("ap-southeast-1")
		if err != nil {
			log.Fatalf("Failed to initialize textract client: %v", err)
		} 
		s3Bucket := utils.Presigner(utils.Presigner{})
		
		//inisialisasi aplikasi fiber dengan konfigurasi khusus
		app := fiber.New(fiber.Config{
			Prefork: false,
			BodyLimit: 10 * 1024 * 1024,
			AppName: "restapi_canteen v1.0.0",
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				c.Status(fiber.StatusInternalServerError)
				return c.SendString("Error: " + err.Error())
			},
		})

		dbPool := db_drivers.InitDbPool(app)
		defer dbPool.Close()

		baseHandler:= &handlers.BaseHandler{
			DB: dbPool,	
			Presigner: s3Bucket,
		}
		userHandler:= handlers.NewUserHandlers(baseHandler)
		sellerHandlers:= handlers.NewSellerHandler(baseHandler)
		// googleHandler:= handlers.NewGoogleHandlers(baseHandler)
	
		routes.SetupRoutes(app, sellerHandlers, userHandler)

		port := os.Getenv("PORT")
		if port==""{
			port="3000"
		}
		if err := app.Listen("0.0.0.0:" + port); err != nil {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}




