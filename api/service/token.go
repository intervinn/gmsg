package service

import (
	"crypto/rsa"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	priv *rsa.PrivateKey
	pub  *rsa.PublicKey
}

func NewTokenService(pubkey string, privkey string) *TokenService {
	priv, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte(privkey))
	pub, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(pubkey))

	return &TokenService{
		priv: priv,
		pub:  pub,
	}
}

func (a *TokenService) Verify(str string) (*jwt.Token, error) {
	token, err := jwt.Parse(str, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return a.pub, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (a *TokenService) Generate(id int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(id, 10),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	tstr, err := token.SignedString(a.priv)
	if err != nil {
		return "", err
	}

	return tstr, nil
}
