package response

const (
	CodeSuccess            = 1  // Success
	ErrCodeParamInvalid    = -4 // Email is invalid
	ErrInvalidToken        = -3 // Token invalid
	ErrorCodeUserHasExited = -2 // User is existed
	ErrInvalidOTP          = -5 /// Invalid otp
	ErrSendEmailOTP        = -6

	//Login response

	ErrCodeOTPNotExists = -7
	ErrSetRedis         = -11
)

// message
var msg = map[int]string{
	CodeSuccess:            "Success",
	ErrCodeParamInvalid:    "Email is invalid",
	ErrInvalidToken:        "Token is invalid",
	ErrorCodeUserHasExited: "User has existed",
	ErrInvalidOTP:          "OTP error",
	ErrSendEmailOTP:        "Failed to send email",
	ErrSetRedis:            "Failed to set value to redis",
	ErrCodeOTPNotExists:    "OTP loign is not exists",
}
