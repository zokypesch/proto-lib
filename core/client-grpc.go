package core

import (
	"context"
	"time"

	"google.golang.org/grpc/metadata"
)

// CreateContextForOutgoing for create context for outgoing request
func CreateContextForOutgoing(req map[string]string) context.Context {
	md := metadata.MD{}
	for k, v := range req {
		md[k] = []string{v}
	}

	ctx := newCtx(5 * time.Second)

	newCtx := metadata.NewIncomingContext(ctx, md)
	incoming, _ := metadata.FromIncomingContext(newCtx)

	outGoing := metadata.NewOutgoingContext(newCtx, incoming)
	return outGoing

}

func newCtx(timeout time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.TODO(), timeout)
	return ctx
}
