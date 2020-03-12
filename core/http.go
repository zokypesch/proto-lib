package core

import (
	"context"
	"fmt"
	"log"
	"net/http"

	runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// RunHTTP for running http
func RunHTTP(init func() error, registerHandler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error), GRPCAddress string, GRPCPort string, HTTPPort string) error {
	if err := init(); err != nil {
		return err
	}
	runtime.HTTPError = CustomHTTPError
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Println("masuk")
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := registerHandler(ctx, mux, fmt.Sprintf("%s:%s", GRPCAddress, GRPCPort), opts)

	log.Println("masuk")
	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%s", HTTPPort), mux)
}
