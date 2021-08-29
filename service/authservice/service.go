package authservice

import (
	"time"

	"github.com/alex-a-renoire/sigma-homework/model"
	"github.com/dgrijalva/jwt-go"
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
	jwt.StandardClaims
	Name string `json:"email,omitempty"`
}

type requestContextKey string

const sessionTokenTTL = 30 * time.Minute //token time to live and redis entry expiration time
const confirmationTokenTTL = 24 * time.Hour

func (as AuthService) generateConfirmationToken(p model.Person) (string, error) {
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
	/*	user := shop.User{}
		if err := json.Unmarshal(ba, &user); err != nil {
			s.reportError(w, http.StatusBadRequest, err)
			return
		}

		pass, group, err := s.storage.Login(user)
		if err != nil {
			if err == sql.ErrNoRows {
				s.reportError(w, http.StatusForbidden, fmt.Errorf("no user with %s email found", user.Email))
			} else {
				s.reportError(w, http.StatusInternalServerError, err)
			}
			return
		}

		if err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(user.Password)); err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				s.reportError(w, http.StatusBadRequest, storage.ErrInvalidCredentials)
			} else {
				s.reportError(w, http.StatusBadRequest, err)
			}
			return
		}

		claims := SessionTokenClaims{
			StandardClaims: jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(sessionTokenTTL).Unix(),
			},
			Email: user.Email,
			Group: group,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(s.jwtSecret)
		if err != nil {
			s.reportError(w, http.StatusInternalServerError, err)
			return
		}

		type response struct {
			Token string
		}

		r := response{Token: tokenString}

		ba, err = json.Marshal(r)
		if err != nil {
			s.reportError(w, http.StatusInternalServerError, err)
			return
		}*/
	return nil
}

func (as AuthService) Authenticate(tokenHeader string) error {
	/*

		group := "anonymous"

		if tokenHeader != "" {
			splitted := strings.Split(tokenHeader, " ")
			if len(splitted) != 2 {
				return fmt.Errorf("malformed auth token, could not split two parts")
			}

			if splitted[0] != "bearer" {
				return s.reportError(w, http.StatusForbidden, fmt.Errorf("malformed auth token, the first part is not bearer"))
			}

			tokenPart := splitted[1]

			token, err := jwt.ParseWithClaims(tokenPart, &SessionTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
				return s.jwtSecret, nil
			})
			if err != nil {
				log.Print(err)
				return fmt.Errorf("malformed auth token, could not parse with claims")
			}
			if claims, ok := token.Claims.(*SessionTokenClaims); ok && token.Valid {
				_, err := s.storage.GetUserByEmail(claims.Email)
				if err != nil {
					return fmt.Errorf("no such user")
				}
				if claims.ExpiresAt < time.Now().Unix() {
					return fmt.Errorf("token expired")
				}
				invTime, err := s.timeStorage.GetRecord(claims.Email + "_invalidationtime")
				if err != nil {
					return err
				}
				if invTime != nil && invTime.After(time.Unix(claims.IssuedAt, 0)) {
					return fmt.Errorf("token invalidated")

				}
				group = claims.Group

				ctx := context.WithValue(r.Context(), requestContextKey("email"), claims.Email)
				r = r.WithContext(ctx)
			}
		}

		path, method := r.URL.Path, r.Method
		if !s.authEnforcer.Enforce(group, path, method) {
			log.Printf("access forbidden: %s, %s, %s", group, path, method)

			s.reportError(w, http.StatusForbidden, fmt.Errorf("you are not allowed to access this page"))
			return
		}*/

	return nil
}
