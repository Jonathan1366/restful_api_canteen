package entity

import (
	"time"

	"github.com/google/uuid"
)

type Seller struct {
	IdSeller      uuid.UUID `json:"id_seller"`
	NamaSeller    string `json:"nama_seller"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	PhoneNum      string `json:"phone_num"`
	Contract      string `json:"kontrak_canteen"`
	Profile_Image string `json:"profile_pic"`
}

type RegistSeller struct {
	IdSeller      uuid.UUID `json:"id_seller"`
	NamaSeller    string `json:"nama_seller"`
	Email         string `json:"email"`
 	Password      string `json:"password"`
	PhoneNum      string `json:"phone_num"`
}

type LoginSeller struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type OTP struct{
	IdSeller string `json:"id_seller"`
	OTP string `json:"otp"`
	ExpiryTime time.Time `json:"expiry_time"`
}

type OTPRequest struct{
	IdSeller string `json:"id_seller"`
	PhoneNum string `json:"phone_num"`
}

type VerifyOTP struct{
	IdSeller string `json:"id_seller"`
	PhoneNum string `json:"phone_num"`
	OTP string `json:"otp"`
}

type User struct {
	IdUsers      uuid.UUID `json:"id_users"`
	NamaUsers    string `json:"nama_users"`
	Email         string `json:"email"`
	Password      string `json:"password"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ContractFile struct{
	IdSeller string `json:"id_seller"`
	PhoneNum string `json:"phone_num"`
	FileContract string `json:"file_contract"`
}

type VerifyContractFile struct{
	IdSeller string `json:"id_seller"`
	PhoneNum string `json:"phone_num"`	
	FileContract string `json:"file_contract"`
	Mime_type string `json:"mime_type"`
	File_content string `json:"file_content"`
	Is_Validate string `json:"is_validate"`
	Validation_metadata string `json:"validation_metadata"`
}

type Location struct{
	IdToko uuid.UUID `json:"id_toko"`
	NamaToko string `json:"nama_toko"`
	Alamat string `json:"alamat"`
}

type GoogleLoginReq struct {
  IDToken string `json:"id_token"`
  Code    string `json:"code"`
  Role    string `json:"role"`
}

type GoogleUser struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Sub   string `json:"sub"` // unique Google user ID
}

//UPDATE SELLER STORE
type StoreSeller struct{
	IdSeller string `json:"id_seller"`
	Loc_seller string `json:"loc_seller"`
	Store_seller string `json:"store_seller"`
}

