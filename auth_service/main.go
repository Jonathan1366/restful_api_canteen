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
	"github.com/jackc/pgx/v5"
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

		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		if len(jwtSecret) == 0 {
			log.Fatal("JWT_SECRET environment variable is not set")
		}

		baseHandler:= &handlers.BaseHandler{
			DB: dbPool,	
			RedisClient: utils.RedisClient,
			DefaultQueryExecMode: pgx.QueryExecModeExec,
			Presigner: s3Bucket,
			JWTSecret: jwtSecret,
		}

		userHandler:= handlers.NewUserHandlers(baseHandler)
		sellerHandlers:= handlers.NewSellerHandler(baseHandler)
		googleHandler:= handlers.NewGoogleHandlers(baseHandler)
	
		routes.SetupRoutes(app, sellerHandlers, userHandler, googleHandler)

		port := os.Getenv("PORT")
		if port==""{
			port="3000"
		}
		if err := app.Listen("0.0.0.0:" + port); err != nil {
			log.Fatalf("Failed to start server: %v\n", err)
	}
}




