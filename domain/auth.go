package domain

type LoginBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterBody struct {
	FullName        string `json:"full_name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type VerifyPhoneBody struct {
	Phone       int64   `json:"phone"`
	CountryCode int32   `json:"countryCode"`
	Type        OtpType `json:"type" validate:"required,oneof=sms whatsapp"`
}

// Define an enum for OTP type
type OtpType string

const (
	Sms      OtpType = "sms"
	Whatsapp OtpType = "whatsapp"
	Email    OtpType = "email"
)

// Define the request structure for verifying OTP
type VerifyOtpBody struct {
	Phone       int64  `json:"phone" validate:"required"`
	CountryCode int32  `json:"countryCode" validate:"required"`
	Code        string `json:"otp" validate:"required"`
}
