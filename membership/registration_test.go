package membership

import (
	"code.google.com/p/go.crypto/bcrypt"
	"testing"
)

type stubUserRepo struct{}

func (userRepo stubUserRepo) Add(*User) (id int, err error) {
	return 1, nil
}

func (userRepo stubUserRepo) GetByEmail(email string) (*User, error) {
	if email == "existing@example.com" {
		return NewUser("existing@example.com"), nil
	}
	return nil, nil
}

func TestRegistrationWithoutEmail(t *testing.T) {
	var stub stubUserRepo
	a := NewApplication("", "secret", "secret")
	r := NewRegistrationManager(stub)

	user, _, err := r.Apply(a)

	if err == nil {
		t.Errorf("expected registration to fail got user %d", user)
	}
}

func TestRegistrationWithoutPassword(t *testing.T) {
	var stub stubUserRepo
	a := NewApplication("fredrik@example.com", "", "secret")
	r := NewRegistrationManager(stub)

	user, _, err := r.Apply(a)

	if err == nil {
		t.Errorf("expected registration to fail got user %d", user)
	}
}

func TestRegistrationNotMatchingConfirm(t *testing.T) {
	var stub stubUserRepo
	a := NewApplication("fredrik@example.com", "secret", "other")
	r := NewRegistrationManager(stub)

	user, _, err := r.Apply(a)

	if err == nil {
		t.Errorf("expected registration to fail got user %d", user)
	}
}

func TestExistingEmail(t *testing.T) {
	var stub stubUserRepo
	a := NewApplication("existing@example.com", "secret", "secret")
	r := NewRegistrationManager(stub)

	user, _, err := r.Apply(a)

	if err == nil {
		t.Errorf("expected registration to fail got user %d", user)
	}
}

func TestSucces(t *testing.T) {
	var stub stubUserRepo
	a := NewApplication("fredrik@example.com", "secret", "secret")
	r := NewRegistrationManager(stub)

	user, token, err := r.Apply(a)

	if err != nil {
		t.Errorf("Expected registration to succeed got %v", err)
	}

	if user.Id != 1 {
		t.Errorf("Expected user id 1 got %d", user.Id)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(a.Password)) != nil {
		t.Errorf("Expected hashed password got %s", user.Password)
	}

	if len(token) == 0 {
		t.Error("Expected a token got %s", token)
		return
	}

}
