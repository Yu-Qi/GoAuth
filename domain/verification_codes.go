package domain

// VerificationCodeService provides the service to generate and verify verification code
type VerificationCodeService interface {
	GenerateCode(uid string) (string, error)
	VerifyCode(code string) (string, error)
}
