package core

import (
	"context"
	"time"

	"github.com/zokypesch/proto-lib/utils"
	"google.golang.org/grpc/metadata"
)

// CreateContextForOutgoing for create context for outgoing request
func CreateContextForOutgoing(ctxIncome context.Context, req map[string]string) context.Context {
	md := metadata.MD{}
	for k, v := range req {
		md[k] = []string{v}
	}
	auth, ok := utils.GetTokenFromCtx(ctxIncome)

	if ok {
		md["authorization"] = []string{auth}
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
