package tool

import (
	"errors"
	"github.com/TskFok/OpenAi/app/global"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	Email                string `json:"email,omitempty"`
	Uid                  uint32 `json:"uid,omitempty"`
	jwt.RegisteredClaims `json:"jwt.RegisteredClaims"`
}

func JwtToken(id uint32, email string) (string, error) {
	secret := []byte(global.JwtSecret)

	newClaims := &Claims{
		Email: email,
		Uid:   id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "user_system",
			//三小时超时
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(global.JwtExpire) * time.Hour)),
			//生效时间
			NotBefore: jwt.NewNumericDate(time.Now()),
			//签发时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	token, err := t.SignedString(secret)

	if nil != err {
		return "", err
	}

	return token, nil
}

func TokenInfo(token string) (*Claims, error) {
	tokens, jwtErr := jwt.ParseWithClaims(token, &Claims{}, secret())

	if nil != jwtErr {
		return nil, jwtErr
	}

	if claims, ok := tokens.Claims.(*Claims); ok && tokens.Valid {
		return claims, nil
	}

	return nil, errors.New("unknown error")
}

func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(global.JwtSecret), nil
	}
}
