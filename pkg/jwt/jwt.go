package jwt

import (
	"errors"
	"html"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gianlucapastori/nausicaa/cfg"
	"github.com/gianlucapastori/nausicaa/internal/entities"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type UserClaims struct {
	UserId   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

func RequestTokens(user *entities.User, cfg *cfg.Config) (string, string, error) {
	accclaims := UserClaims{
		UserId:   user.Id,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Second)),
		},
	}

	refclaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(rand.Int31n(cfg.SERVER.JWT.REFRESH_EXPIRES_AT)) * time.Minute)),
	}

	acctoken := jwt.NewWithClaims(jwt.SigningMethodHS256, accclaims)
	reftoken := jwt.NewWithClaims(jwt.SigningMethodHS256, refclaims)
	accTokenStr, err := acctoken.SignedString([]byte(cfg.SERVER.JWT.SECRET_KEY_ACCESS))
	if err != nil {
		return "", "", err
	}

	refTokenStr, err := reftoken.SignedString([]byte(cfg.SERVER.JWT.SECRET_KEY_REFRESH))
	if err != nil {
		return "", "", err
	}

	return accTokenStr, refTokenStr, nil
}

func ExtractBearer(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	bearer := strings.Split(header, " ")

	if bearer[0] != "Bearer" {
		return "", errors.New("malformed bearer token")
	}
	if len(bearer) != 2 {
		return "", errors.New("malformed bearer token")
	}

	return html.EscapeString(bearer[1]), nil
}
