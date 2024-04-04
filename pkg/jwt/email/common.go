package email

import "github.com/Yu-Qi/GoAuth/domain"

var (
	service domain.SendEmailService
)

func GetService() domain.SendEmailService {
	return service
}

func InitService(s domain.SendEmailService) {
	service = s
}
