package core

import (
	"context"
	"fmt"
	"net/http"

	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

var pattern_health_check_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0}, []string{"health"}, ""))

// RunHTTP for running http
func RunHTTP(init func() error, registerHandler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error), GRPCAddress string, GRPCPort string, HTTPPort string) error {
	if err := init(); err != nil {
		return err
	}
	// breaking in version 2
	// runtime.HTTPError = CustomHTTPError
	runtime.WithErrorHandler(CustomHTTPError)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := registerHandler(ctx, mux, fmt.Sprintf("%s:%s", GRPCAddress, GRPCPort), opts)

	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%s", HTTPPort), mux)
}

// RunHTTPWithCustomMatcher for running http
func RunHTTPWithCustomMatcher(init func() error, registerHandler func(ctx context.Context, mux *runtime.ServeMux,
	endpoint string, opts []grpc.DialOption) (err error), customMatcher func(key string) (string, bool),
	GRPCAddress string, GRPCPort string, HTTPPort string) error {
	if err := init(); err != nil {
		return err
	}
	// breaking in version 2
	// runtime.HTTPError = CustomHTTPError

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(customMatcher),
		runtime.WithErrorHandler(CustomHTTPError),
	)

	mux.Handle("GET", pattern_health_check_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		w.WriteHeader(200)
	})

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := registerHandler(ctx, mux, fmt.Sprintf("%s:%s", GRPCAddress, GRPCPort), opts)

	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%s", HTTPPort), mux)
}
