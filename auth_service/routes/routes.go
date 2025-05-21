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
	seller.Post("/presignurl", Seller.GeneratePresignedUploadURL)
	
	//USER
	user:= auth.Group("/user")
	user.Post("/register", User.RegisterUser)
	user.Post("/login", User.LoginUser)
	user.Post("/logout", User.LogoutUser)

	//OAUTH2
	auth.Get("/google/login", google.GoogleSignIn)
	// auth.Get("google/callback", google.GoogleCallBack)

	// seller:= app.Group("/seller")
	// seller.Post("/register", Seller.RegisterSeller)
	// seller.Post("/login", Seller.LoginSeller)
	// seller.Post("/send-otp", Seller.SendOTP)
	// seller.Post("/verify-otp", Seller.VerifyOTP)
	// seller.Post("/logout", Seller.LogoutSeller)

	// //OAuth2 Google SignIn
	// auth:= app.Group("/auth")
	// auth.Post("/google-signin", google.GoogleLogin)

	
	// //s3bucket
	// seller.Post("/presignurl", Seller.GeneratePresignedUploadURL)
	// // seller.Get("/download-url", authHandler.GeneratePresignedDownloadURL)
	// user:=app.Group("/user")
	// user.Post("/register", User.RegisterUser)
	// user.Post("/login", User.LoginUser)
	// user.Post("/logout", User.LogoutUser)

	// //STORE SELLER
	// seller.Post("/store", Seller.StoreSeller)



	
}