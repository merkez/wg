package wg

import (
	"context"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/metadata"
)

const (
	AUTH_KEY = "wg"
)

var (
	InvalidAuthKey        = errors.New("Invalid Authentication Key")
	InvalidTokenFormatErr = errors.New("Invalid token format")
	MissingKeyErr         = errors.New("No Authentication Key provided")
)

type Authenticator interface {
	AuthenticateContext(context.Context) error
}

type auth struct {
	sKey string // Signin Key
	aKey string // Auth Key
}

func NewAuthenticator(Skey, AKey string) Authenticator {
	return &auth{sKey: Skey, aKey: AKey}
}

func (a *auth) AuthenticateContext(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return MissingKeyErr
	}

	if len(md["token"]) == 0 {
		return MissingKeyErr
	}

	token := md["token"][0]
	if token == "" {
		return MissingKeyErr
	}

	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return ctx, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(a.sKey), nil
	})
	if err != nil {
		return err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return InvalidTokenFormatErr
	}

	authKey, ok := claims[AUTH_KEY].(string)
	if !ok {
		return InvalidTokenFormatErr
	}

	if authKey != a.aKey {
		return InvalidAuthKey
	}

	return nil
}
