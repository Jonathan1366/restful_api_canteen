package handlers

import (
	entity "ubm-canteen/models"
	"ubm-canteen/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *BaseHandler) VerifyOTP(c *fiber.Ctx) error {
	input:=new(entity.VerifyOTP)
	
// 	//parse body request
	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":"error",
			"message":"gagal memproses permintaan verifikasi",
		})
	}

// 	//validasi input
	if input.IdSeller == "" || input.OTP == "" || input.PhoneNum=="" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":"error",
			"message":"Id Seller dan OTP diperlukan",
		})
	}

	isValid, err := utils.VerifyOTPWithTwilio(input.PhoneNum, input.OTP)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "error",
			"message": "gagal memverifikasi otp", 
		})
	}	

	if !isValid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "error",
			"message": "OTP tidak valid atau sudah kadaluarsa",
		})
	}
	
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"message": "Verifikasi OTP berhasil", 
})
	
// 	//cek otp di database
// 	conn, err:=h.DB.Acquire(ctx)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"status":"error",
// 			"message":"gagal mengakses database",
// 		})
// 	}
// 	defer conn.Release()

// 	var otpData entity.OTP

// 	query := `SELECT otp, expiry_time from otps where id_seller = $1`
// 	err = conn.QueryRow(ctx, query, input.IdSeller).Scan(&otpData.OTP, &otpData.ExpiryTime)
// 	if err != nil {
// 		if err == pgx.ErrNoRows {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"status":"error",
// 				"message":"OTP tidak valid atau IdSeller tidak ditemukan",
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"status":"error",
// 			"message":"Gagal memverifikasi OTP",
// 		})
// 	}
// 	if time.Now().After(otpData.ExpiryTime) {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"status":"error",
// 			"message": "OTP sudah kadaluarsa",
// 		})
// 	}
// 	if otpData.OTP!=input.OTP {
// 		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 			"status":"error",
// 			"message":"invalid otp",
// 		})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"status": "success",
// 		"message": "Verifikasi OTP berhasil",
// 	})
// }


}
