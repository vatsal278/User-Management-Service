package helpers

type LoginService interface {
	LogInUser(email string, password string) bool
}

type LoginInformation struct {
	email    string
	password string
}

func StaticLoginService() LoginService {
	return &LoginInformation{
		email:    "",
		password: "",
	}
}

func (info *LoginInformation) LogInUser(email string, password string) bool {
	return info.email == email && info.password == password
}
