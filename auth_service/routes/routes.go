package routes

import (
	"ubm-canteen/handlers"
	"ubm-canteen/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, sellerHandlers *handlers.SellerHandler, userHandlers *handlers.UserHandler, googleHandlers  *handlers.GoogleHandler) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":"404 not found",
		})
	});
	
	v1:= app.Group("/api/v1")
	
	auth:= v1.Group("/auth")

	//SELLER
	seller:= auth.Group("/seller", middleware.AuthMiddleware("seller"))
	seller.Post("/register", sellerHandlers.RegisterSeller)
	seller.Post("/login", sellerHandlers.LoginSeller)
	seller.Post("/otp/send", sellerHandlers.SendOTP)
	seller.Post("/otp/verify", sellerHandlers.VerifyOTP) 
	seller.Post("/logout", sellerHandlers.LogoutSeller)
	seller.Put("/store/location", sellerHandlers.StoreLocSeller)

	//S3 BUCKET
	// seller.Post("/presignurl", Seller.GeneratePresignedUploadURL)
	
	//USER
	user:= auth.Group("/user")
	user.Post("/register", userHandlers.RegisterUser)
	user.Post("/login", userHandlers.LoginUser)
	user.Post("/logout", userHandlers.LogoutUser)

	//OAUTH2
	google:= auth.Group("/google")
	google.Post("/login", googleHandlers.GoogleSignIn)
}