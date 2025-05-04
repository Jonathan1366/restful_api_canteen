package handlers

import (
	entity "ubm-canteen/models"
	"ubm-canteen/utils"

	"github.com/gofiber/fiber/v2"
)

func (h *BaseHandler) SendOTP(c *fiber.Ctx) error {
	input := new(entity.OTPRequest)
	ctx:=c.Context()
	
	//parse body request
	if err:=c.BodyParser(input); err!=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":"error",
			"message":"Gagal memproses permintaan OTP",
		})
	}

	//validasi input
	if input.IdSeller=="" || input.PhoneNum=="" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":"error",
			"message":"Id Seller dan PhoneNum diperlukan",
		})	
	}
	
	conn, err:= h.DB.Acquire(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":"Error",
			"message":"gagal mengakses database",
		})
	}
	defer conn.Release()
	
	err = utils.SendOTPwithTwilio(input.PhoneNum)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":"error",
			"message":"gagal mengirim OTP via Twilio",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":"success",
		"message":"TOTP berhasil dikirim ke nomor telp anda",
	})
}