package utils

import (
	"context"
	"fmt"

	"google.golang.org/grpc/metadata"
	meta "google.golang.org/grpc/metadata"
)

// GetTokenFromCtx for get token from context
func GetTokenFromCtx(ctx context.Context) (string, bool) {
	md, ok := meta.FromIncomingContext(ctx)

	if ok {
		v, okToken := md["authorization"]

		if !okToken {
			return "auth not found", false
		}

		if len(v) == 0 {
			return "zero length of context", false
		}

		return v[0], true
	}
	return "cannot get token", false
}

// GetXCustomKeyFromCtx for get token from context
func GetXCustomKeyFromCtx(ctx context.Context, key string) (string, bool) {
	md, ok := meta.FromIncomingContext(ctx)

	if ok {
		v, okToken := md[key]

		if !okToken {
			return fmt.Sprintf("%s not found", key), false
		}

		if len(v) == 0 {
			return "zero length of context", false
		}

		return v[0], true
	}
	return fmt.Sprintf("cannot get %s", key), false
}

// OutgoingContext for outgoing context
func OutgoingContext(ctx context.Context, params map[string]string) context.Context {
	mdReceive, okReceive := meta.FromIncomingContext(ctx)
	if !okReceive {
		return ctx
	}

	for k, v := range params {
		mdReceive.Append(k, v)
	}
	ctxPass := context.TODO()
	newCtx := meta.NewIncomingContext(ctxPass, mdReceive)
	return newCtx
}

// NewOutgoingContext for outgoing context
func NewOutgoingContext(ctx context.Context, params map[string]string) context.Context {
	mdReceive := meta.New(params)

	outGoing := metadata.NewOutgoingContext(ctx, mdReceive)
	return outGoing
}
