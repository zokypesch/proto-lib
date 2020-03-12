package utils

import (
	"context"

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
