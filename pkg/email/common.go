package email

var (
	service SendEmailService
)

type SendEmailService interface {
	SendEmail(email string, subject string, body string) error
}

func GetService() SendEmailService {
	return service
}

func InitService(s SendEmailService) {
	service = s
}
