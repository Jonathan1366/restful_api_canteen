package handlers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	entity "ubm-canteen/models"
	"ubm-canteen/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SellerHandler struct {
	*BaseHandler
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func NewSellerHandler(base *BaseHandler) *SellerHandler {
	return &SellerHandler{base}
}

// REGISTER SELLER
func (h *SellerHandler) RegisterSeller(c *fiber.Ctx) error {
	input := new(entity.RegistSeller)
	ctx := c.Context()

	//DEALLOCATE PREPARE STATEMENT IF EXISTS
	err := utils.DeallocateAllStatement(ctx, h.DB)
	if err != nil && err != pgx.ErrNoRows {
		//ignore jika tidak ada yg terdallocate
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to deallocate statement",
		})
	}

	// PARSE BODY REQUEST
	if err := c.BodyParser(input); err != nil {
		fmt.Printf("Failed to parse body: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"code":    400,
			"message": "Invalid input: Failed to parse request body",
		})
	}

	// VALIDATE INPUT
	if input.Email == "" || input.Password == "" || input.NamaSeller == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"code":    400,
			"message": "Invalid Input: Email, sellername, and Password are required.",
			"details": fiber.Map{
				"missing_fields": []string{"email", "password", "nama_seller"},
			},
		})
	}

	// HASH PASSWORD
	hashedPass, err := utils.HashPass(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":   500,
			"status": "Failed to hash password",
		})
	}

	seller := entity.Seller{
		IdSeller:   uuid.New(),
		NamaSeller: input.NamaSeller,
		Email:      input.Email,
		Password:   hashedPass, // Assume password is hashed
		PhoneNum:   input.PhoneNum,
	}

	//SAVE SELLER TO DATABASE
	conn, err := h.DB.Acquire(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Gagal mendapatkan koneksi: %v", err),
		})
	}
	defer conn.Release() // MAKE SURE TO RELEASE CONNECTION TO POOL
	// EXECUTE PREPARE STATEMENT FOR REGISTER SELLER
	query := `INSERT INTO seller (id_seller, nama_seller, email, password, phone_num) VALUES ($1, $2, $3, $4, $5)`
	_, err = conn.Exec(ctx, query, seller.IdSeller, seller.NamaSeller, seller.Email, hashedPass, seller.PhoneNum)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  "error",
				"message": "Email already exist",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to register seller: %v", err),
		})
	}

	// go func() {
	// 	err := utils.SendOTPwithTwilio(input.PhoneNum)
	// 	if err != nil {
	// 		fmt.Printf("Failed to send OTP: %v\n", err)
	// 	} else{
	// 		fmt.Println("OTP sent successfully")
	// 	}
	// }()

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "Success",
		"message": "Seller successfully registered, OTP has been sent to your phone",
		"data": fiber.Map{
			"id_seller":   seller.IdSeller.String(),
			"email":       seller.Email,
			"nama_seller": seller.NamaSeller,
			"password":    seller.Password,
			"phone_num":   seller.PhoneNum,
		},
	})
}

// lOGIN SELLER
func (h *SellerHandler) LoginSeller(c *fiber.Ctx) error {
	seller := new(entity.Seller)
	ctx := c.Context()

	// PARSE BODY REQUEST
	if err := c.BodyParser(seller); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// VALIDATE INPUT	
	if seller.Email == "" || seller.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	// CHECK IF SELLER EXISTS AND PASSWORD MATCH
	conn, err := h.DB.Acquire(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Gagal mendapatkan koneksi: %v", err),
		})
	}
	defer conn.Release() // MAKE SURE TO RELEASE CONNECTION TO POOL

	// EXECUTE PREPARE STATEMENT FOR LOGIN
	query := `SELECT id_seller, email, password FROM seller WHERE email=$1`
	dbseller := new(entity.Seller)
	err = conn.QueryRow(ctx, query, seller.Email).Scan(&dbseller.IdSeller, &dbseller.Email, &dbseller.Password)
	if err != nil {
		if err == pgx.ErrNoRows {
			fmt.Printf("Query execution error: %v\n", err) // Log error jika ada
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to execute query",
		})
	}

	// VERIFY HASHED PASSWORD
	if !utils.CheckPassHash(seller.Password, dbseller.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid password",
		})
	}
	
	//JWT TOKEN
	tokenString, err:= utils.GenerateJWTSecret(dbseller.IdSeller.String(), dbseller.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			
		})
	}

	//REFRESH TOKEN
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_seller":  dbseller.IdSeller.String(),
		"email":      dbseller.Email,
		"ip_address": c.IP(),
		"user_agent": c.Get("User-Agent"),
		"exp":        time.Now().Add(time.Hour * 24 * 30).Unix(), //valid for 30 days
	})
	
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	//SAVE TOKEN TO REDIS
	err = utils.RedisClient.Set(ctx, "token:"+dbseller.IdSeller.String(), tokenString, time.Hour*24).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to store token in Redis",
		})
	}

	//SAVE REFRESH TOKEN TO REDIS
	err = utils.RedisClient.Set(ctx, "refresh token:"+dbseller.IdSeller.String(), refreshTokenString, time.Hour*24*30).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to store refresh token in Redis",
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "Success",
		"message": "Seller successfully logged in",
		"data": fiber.Map{
			"email":         seller.Email,
			"token":         tokenString,
			"refresh_token": refreshTokenString,
		},
	})
}

// LOGOUT SELLER
func (h *SellerHandler) LogoutSeller(c *fiber.Ctx) error {
	//invalid jwt token (for example, by storing it in a blacklist)
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "no token provided",
		})
		// remove bearer from token string if present
	}

	token = strings.TrimPrefix(token, "Bearer ")

	//DELETE TOKEN FROM REDIS
	ctx := c.Context()
	err := utils.RedisClient.Del(ctx, "token:"+token).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to delete token from Redis",
		})
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid token or error in parsing token",
		})
	}

	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		idSeller := claims["id_seller"].(string)
		//input token to revocation list
		err := h.TokenRevocationLogic(uuid.MustParse(idSeller), token)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "failed to logout and revoke token",
			})
		}
		
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "Success",
			"message": "logged out successfully",
		})
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid token",
		})
	}
}

func (h *SellerHandler) TokenRevocationLogic(d uuid.UUID, token string) any {
	panic("unimplemented")
}

func (h *SellerHandler) StoreSeller(c *fiber.Ctx) error {

	input := new(entity.Location)
	ctx := c.Context()

	err := utils.DeallocateAllStatement(ctx, h.DB)
	if err != nil && err != pgx.ErrNoRows {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to deallocate statement",
		})
	}

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"code":    400,
			"message": "Failed to parse body",
		})
	}

	if input.NamaToko == "" || input.Alamat == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"code":    400,
			"message": "Invalid input: Nama Toko and Alamat are required",
		})
	}

	//save to database
	store := entity.Location{
		IdToko:   uuid.New(),
		NamaToko: input.NamaToko,
		Alamat:   input.Alamat,
	}

	conn, err := h.DB.Acquire(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to acquire connection: %v", err),
		})
	}
	defer conn.Release()

	//EXECUTE PREPARE STATEMENT FOR STORE LOCATION
	query := `INSERT INTO toko (id_toko, nama_toko, alamat) VALUES ($1, $2, $3)`
	_, err = conn.Exec(ctx, query, store.IdToko, store.NamaToko, store.Alamat)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  "error",
				"message": "Location already exists",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to store location: %v", err),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Location successfully stored",
		"data": fiber.Map{
			"id_toko":   store.IdToko.String(),
			"nama_toko": store.NamaToko,
			"alamat":    store.Alamat,
		},
	})
}

// PRESIGNED URL AWS S3
func (h *SellerHandler) GeneratePresignedUploadURL(c *fiber.Ctx) error {
	fileName := c.Query("filename")
	if fileName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "filename is required as query param",
		})
	}

	// Load AWS Config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to load AWS config: " + err.Error(),
		})
	}

	s3Client := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(s3Client)

	bucket := os.Getenv("S3_BUCKET_NAME")
	expireDuration := 15 * time.Minute

	// Generate presigned URL
	presignResult, err := presignClient.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fileName),
		ContentType: aws.String("application/octet-stream"),
	}, s3.WithPresignExpires(expireDuration))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to generate presigned URL: " + err.Error(),
		})
	}

	// Return presigned URL dengan metode PUT
	return c.JSON(fiber.Map{
		"upload_url": presignResult.URL,
		"method":     "PUT",
		"headers": fiber.Map{
			"Content-Type": "application/octet-stream", // Sesuai dengan ekstensi file
		},
	})
}
