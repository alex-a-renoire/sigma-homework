package authservice

import (
	"fmt"
	"strings"
	"time"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/alex-a-renoire/sigma-homework/service/personservice"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AuthService struct {
	ps        personservice.PersonService
	jwtSecret []byte
}

func New(ps personservice.PersonService, jwtSecret string) AuthService {
	return AuthService{
		ps:        ps,
		jwtSecret: []byte(jwtSecret),
	}
}

type SessionTokenClaims struct {
	jwt.StandardClaims
	Id   uuid.UUID
	Name string `json:"email,omitempty"`
}

type requestContextKey string

const sessionTokenTTL = 30 * time.Minute //token time to live and redis entry expiration time

func (as AuthService) GenerateSessionToken(p model.Person) (string, error) {
	claims := SessionTokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(sessionTokenTTL).Unix(),
		},
		Id:   p.Id,
		Name: p.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(as.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign string: %w", err)
	}
	return tokenString, nil
}

func (as AuthService) MyUser(tokenHeader string) (model.Person, error) {
	if tokenHeader == "" {
		return model.Person{}, fmt.Errorf("empty auth token")
	}

	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		return model.Person{}, fmt.Errorf("malformed auth token, could not split two parts")
	}

	if splitted[0] != "bearer" {
		return model.Person{}, fmt.Errorf("malformed auth token, the first part is not bearer")
	}

	tokenPart := splitted[1]

	token, err := jwt.ParseWithClaims(tokenPart, &SessionTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return as.jwtSecret, nil
	})
	if err != nil {
		return model.Person{}, fmt.Errorf("malformed auth token, could not parse with claims: %w", err)
	}

	if claims, ok := token.Claims.(*SessionTokenClaims); ok && token.Valid {
		if _, err := as.ps.GetPerson(claims.Id); err != nil {
			return model.Person{}, fmt.Errorf("no such user: %w", err)
		}

		if claims.ExpiresAt < time.Now().Unix() {
			return model.Person{}, fmt.Errorf("cannot identify session: token expired")
		}

		return model.Person{
			Id:   claims.Id,
			Name: claims.Name,
		}, nil
	} else {
		return model.Person{}, fmt.Errorf("unknown error")
	}
}
