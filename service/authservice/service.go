package authservice

import (
	"time"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	jwtSecret []byte
}

func New(jwtSecret string) AuthService {
	return AuthService{
		jwtSecret: []byte(jwtSecret),
	}
}

type ConfirmationTokenClaims struct {
	jwt.StandardClaims
	Name string `json:"email,omitempty"`
}

type requestContextKey string

const sessionTokenTTL = 30 * time.Minute //token time to live and redis entry expiration time
const confirmationTokenTTL = 24 * time.Hour

func (as AuthService) GenerateConfirmationToken(p model.Person) (string, error) {
	claims := ConfirmationTokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(confirmationTokenTTL).Unix(),
		},
		Name: p.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(as.jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (as AuthService) Login(ba []byte) error {
	return nil
}
