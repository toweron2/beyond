package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type (
	TokenOptions struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
		RefreshAfter  int64
		Fields        map[string]any
	}
	Token struct {
		AccessToken   string `json:"access_token"`
		AccessExpire  int64  `json:"access_expire"`
		RefreshAfter  int64  `json:"refresh_ffter"`
		RefreshToken  string `json:"refresh_token"`
		RefreshExpire int64  `json:"refresh_expire"`
	}
)

func BuildTokens(opt TokenOptions) (Token, error) {
	var err error
	now := time.Now().Add(-time.Minute).Unix()
	token := Token{
		AccessExpire:  now + opt.AccessExpire,
		RefreshAfter:  now + opt.RefreshAfter,
		RefreshExpire: now + opt.RefreshExpire,
	}

	token.AccessToken, err = getToken(now, opt.AccessSecret, opt.Fields, opt.AccessExpire)
	if err != nil {
		return token, err
	}

	token.RefreshToken, err = getToken(now, opt.RefreshSecret, opt.Fields, opt.RefreshExpire)
	if err != nil {
		return token, err
	}

	return token, nil
}

func getToken(iat int64, secretKey string, payloads map[string]any, second int64) (string, error) {
	claims := jwt.MapClaims{
		"exp": iat + second,
		"iat": iat,
	}
	for k, v := range payloads {
		claims[k] = v
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secretKey))
}
