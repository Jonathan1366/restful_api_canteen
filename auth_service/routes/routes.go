package routes

import (
	"ubm-canteen/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, sellerHandlers, userHandlers, googleHandlers) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"404 not found",
		})
	})

	v1:= app.Group("/api/v1")

	auth:= v1.Group("/auth")

	//SELLER
	seller:= auth.Group("/seller")
	seller.Post("/register", sellerHandlers.RegisterSeller)
	seller.Post("/login", sellerHandlers.LoginSeller)
	seller.Post("/otp/send", sellerHandlers.SendOTP)
	seller.Post("/otp/verify", sellerHandlers.VerifyOTP) 
	seller.Post("/logout", sellerHandlers.LogoutSeller)
	seller.Post("/store", sellerHandlers.StoreSeller)

	//S3 BUCKET
	// seller.Post("/presignurl", Seller.GeneratePresignedUploadURL)
	
	//USER
	user:= auth.Group("/user")
	user.Post("/register", userHandlers.RegisterUser)
	user.Post("/login", userHandlers.LoginUser)
	user.Post("/logout", userHandlers.LogoutUser)

	//OAUTH2
	google:= auth.Group("/google")
	google.Post("/login", googlehandlers.GoogleSignIn)
}