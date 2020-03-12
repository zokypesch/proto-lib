package core

import (
	"context"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/zokypesch/proto-lib/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthInterceptorConfig struct for information interceptor
type AuthInterceptorConfig struct {
	InternalAPIPassword string
}

var authInterceptor *AuthInterceptorConfig

// NewAuthInterceptor for new config
func NewAuthInterceptor(InternalAPIPassword string) *AuthInterceptorConfig {
	if authInterceptor == nil {
		authInterceptor = &AuthInterceptorConfig{InternalAPIPassword}
	}
	return authInterceptor
}

// AuthInterceptor for authentication interceptor
func (auth *AuthInterceptorConfig) AuthInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	token, errAuth := utils.GetTokenFromCtx(ctx)

	if !errAuth {
		// return nil, status.Errorf(codes.Unauthenticated, "you not have access to access this")
		return nil, status.Errorf(codes.Unauthenticated, token)
	}

	if token == auth.InternalAPIPassword {
		return handler(ctx, req)
	}

	tokenJWT, errJWT := GetToken(token)
	if errJWT != nil || !tokenJWT.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "token is not valid")
	}

	// Calls the handler and for return
	h, err := handler(ctx, req)

	return h, err
}

// WhitelistMethod for authentication
func (auth *AuthInterceptorConfig) WhitelistMethod(in grpc.UnaryServerInterceptor, fullMethod []string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		for _, v := range fullMethod {
			if info.FullMethod == v {
				return handler(ctx, req)
			}
		}
		return in(ctx, req, info, handler)
	}
}

// GetUnaryCustom for get unnary custom for auth
func (auth *AuthInterceptorConfig) GetUnaryCustom(fullMethod []string) grpc.UnaryServerInterceptor {
	return auth.WhitelistMethod(grpc.UnaryServerInterceptor(auth.AuthInterceptor), fullMethod)
}

// GetToken function for get token
func GetToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
}
