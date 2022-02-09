package core

import (
	"context"
	"fmt"
	"net/http"

	runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

// RunHTTPWithCustomMatcher for running http
func RunHTTPWithCustomMatcherV2(init func() error, registerHandler func(ctx context.Context, mux *runtime.ServeMux,
	endpoint string, opts []grpc.DialOption) (err error), customMatcher func(key string) (string, bool),
	GRPCAddress string, GRPCPort string, HTTPPort string, patt runtime.Pattern, svcUriPath string) error {
	if err := init(); err != nil {
		return err
	}
	runtime.HTTPError = CustomHTTPError
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(customMatcher))

	mux.Handle("GET", patt, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		w.WriteHeader(200)
	})

	patternV1Health := runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"api", "v1", svcUriPath, "health"}, "", nil...))
	mux.Handle("GET", patternV1Health, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		w.WriteHeader(200)
	})

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
	GRPCAddress string, GRPCPort string, HTTPPort string, patt runtime.Pattern) error {
	if err := init(); err != nil {
		return err
	}
	runtime.HTTPError = CustomHTTPError
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(runtime.WithIncomingHeaderMatcher(customMatcher))

	mux.Handle("GET", patt, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		w.WriteHeader(200)
	})

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := registerHandler(ctx, mux, fmt.Sprintf("%s:%s", GRPCAddress, GRPCPort), opts)

	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%s", HTTPPort), mux)
}

// RunHTTP for running http
func RunHTTP(init func() error, registerHandler func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error), GRPCAddress string, GRPCPort string, HTTPPort string) error {
	if err := init(); err != nil {
		return err
	}
	runtime.HTTPError = CustomHTTPError
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
