package utils

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

var otpRequestCache sync.Map

	func SendOTPwithTwilio(toPhoneNumber string) error  {

		//normalize phone number
		normalizePhopneNumber, err:= NormalizePhoneNumber(toPhoneNumber)
		if err != nil {
			return fmt.Errorf("failed to normalize phopne number: %v", err)
		}

		if _, exists := otpRequestCache.Load(normalizePhopneNumber); exists {
			return fmt.Errorf("OTP request too frequent, please wait")
		}
		
		accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
		authToken := os.Getenv("TWILIO_AUTH_TOKEN")
		serviceSid := os.Getenv("TWILIO_SERVICE_SID")

		if accountSid == "" || authToken =="" || serviceSid ==""{
			return fmt.Errorf("Twilio credentials not set in env")
		}

		client:=twilio.NewRestClientWithParams(twilio.ClientParams{
			Username: accountSid,
			Password: authToken,
		})
	
		//buat parameter verifikasi
		params:= &verify.CreateVerificationParams{}
		params.SetTo(toPhoneNumber)
		params.SetChannel("sms")	

		//kirim verifikasi menggunakan layanan twilio
		resp, err := client.VerifyV2.CreateVerification(serviceSid, params)
		if err != nil {
			fmt.Println("Error sending verification:", err.Error())
        return err
		} 
			if resp.Sid!=nil{
				fmt.Println("Verification SID:", *resp.Sid)
				} else{
				fmt.Println("Verification SID is nil")
				}

				otpRequestCache.Store(normalizePhopneNumber, time.Now())
				time.AfterFunc(1*time.Minute, func(){
					otpRequestCache.Delete(normalizePhopneNumber)
				})
				
		return nil
		
	}

	func VerifyOTPWithTwilio(phoneNum string, code string) (bool , error)  {

		accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
		authToken := os.Getenv("TWILIO_AUTH_TOKEN")
		serviceSid := os.Getenv("TWILIO_SERVICE_SID")

		if accountSid == "" || authToken == "" || serviceSid == ""{
			return false, fmt.Errorf("twilio credentials not set in env")
		}

		client:=twilio.NewRestClientWithParams(twilio.ClientParams{
			Username:accountSid,
			Password: authToken,
		})

		params:= &verify.CreateVerificationCheckParams{}
		params.SetTo(phoneNum)
		params.SetCode(code)

		resp, err := client.VerifyV2.CreateVerificationCheck(serviceSid, params)
		if err != nil {
			return false, err
		}
		if resp.Status!=nil && *resp.Status == "approved" {
			return true, nil
		}	else{
			return false , nil
		}
	}


	func NormalizePhoneNumber(phone string) (string, error){

		//check if phonenumber has started with +
		if len(phone) > 0 && phone[0] == '+' {
			return phone, nil
		}
		
		//if phone num stared with "0", change with the default code's country
		if len(phone) > 1 && phone[0] =='0'{
			return "+62"+phone[1:], nil
		}

		//if format doesn't valid, return error
		return "", fmt.Errorf("invalid phone number format: %s", phone)
	}