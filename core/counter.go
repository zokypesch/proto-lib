package core

import (
	"context"

	"google.golang.org/grpc"
)

// GetUnaryCounter for get unary counter
func GetUnaryCounter(svc string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		customizedCounterMetric.WithLabelValues(svc).Inc()
		return handler(ctx, req)
	}
}
