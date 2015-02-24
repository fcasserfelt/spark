package membership

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
)

type RegistrationManager struct {
	repo UserRepo
}

type Application struct {
	Email    string
	Password string
	Confirm  string
}

func NewApplication(email string, password string, confirm string) *Application {
	return &Application{
		Email:    email,
		Password: password,
		Confirm:  confirm,
	}
}

func NewRegistrationManager(repo UserRepo) *RegistrationManager {
	return &RegistrationManager{repo}
}

func (r RegistrationManager) Apply(application *Application) (*User, string, error) {
	if len(application.Email) == 0 {
		return nil, "", fmt.Errorf("Email required")
	}

	if len(application.Password) == 0 {
		return nil, "", fmt.Errorf("Password required")
	}

	if application.Password != application.Confirm {
		return nil, "", fmt.Errorf("Passowrd and confirm does not match")
	}

	if r.repo.GetByEmail(application.Email) != nil {
		return nil, "", fmt.Errorf("Email already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(application.Password), 10)
	if err != nil {
		return nil, "", err
	}

	var user = NewUser(application.Email)
	user.Password = string(hashedPassword)
	id := r.repo.Add(user)
	user.Id = id

	tokenString, err := GenerateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("Could not generate token %v", err)
	}

	return user, tokenString, nil
}
