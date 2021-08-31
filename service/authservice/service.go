package authservice

import (
	"time"

	"github.com/alex-a-renoire/sigma-homework/model"
)

type AuthService struct {
	jwtSecret []byte
}

func New(jwtSecret []byte) AuthService {
	return AuthService{
		jwtSecret: []byte(jwtSecret),
	}
}

type ConfirmationTokenClaims struct {
	//	jwt.StandardClaims
	Name string `json:"email,omitempty"`
}

type requestContextKey string

const sessionTokenTTL = 30 * time.Minute //token time to live and redis entry expiration time
const confirmationTokenTTL = 24 * time.Hour

func (as AuthService) generateConfirmationToken(p model.Person) (string, error) {
	return "", nil
}

func (as AuthService) Login(ba []byte) error {
	return nil
}
