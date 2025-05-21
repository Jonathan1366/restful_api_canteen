package routes

import (
	"ubm-canteen/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, Seller *handlers.SellerHandler, User *handlers.UserHandler, google *handlers.GoogleHandler) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"404 not found",
		})
	})

	v1:= app.Group("/api/v1")

	auth:= v1.Group("/auth")

	//SELLER
	seller:= auth.Group("/seller")
	seller.Post("/register", Seller.RegisterSeller)
	seller.Post("/login", Seller.LoginSeller)
	seller.Post("/otp/send", Seller.SendOTP)
	seller.Post("/otp/verify", Seller.VerifyOTP) 
	seller.Post("/logout", Seller.LogoutSeller)
	seller.Post("/store", Seller.StoreSeller)

	//S3 BUCKET
	// seller.Post("/presignurl", Seller.GeneratePresignedUploadURL)
	
	//USER
	user:= auth.Group("/user")
	user.Post("/register", User.RegisterUser)
	user.Post("/login", User.LoginUser)
	user.Post("/logout", User.LogoutUser)

	//OAUTH2
	google:= auth.Group("/google")
	google.Post("/login", google.GoogleSignIn)

}