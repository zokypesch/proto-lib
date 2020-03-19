package core

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"regexp"

	proto "github.com/golang/protobuf/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	logrus "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
)

type responseBody interface {
	XXX_ResponseBody() interface{}
}

type errorBody struct {
	Err     string `json:"message,omitempty"`
	Success bool   `json:"success"`
	Code    string `json:"errorCode"`
}

type successResponse struct {
	Err     string      `json:"message,omitempty"`
	Success bool        `json:"success"`
	Code    string      `json:"errorCode"`
	Data    interface{} `json:"data"`
}

// LocalForward for handling localforward append message
func LocalForward(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, resp proto.Message, opts ...func(context.Context, http.ResponseWriter, proto.Message) error) {
	w.Header().Set("Content-Type", "application/json")

	var buf []byte
	var err error

	if rb, ok := resp.(responseBody); ok {
		buf, err = marshaler.Marshal(rb.XXX_ResponseBody())
	} else {
		buf, err = marshaler.Marshal(resp)
	}

	if err != nil {
		grpclog.Infof("Marshal error: %v", err)
		runtime.HTTPError(ctx, mux, marshaler, w, req, err)
		return
	}
	var resMap map[string]interface{}

	err = json.Unmarshal(buf, &resMap)
	if err != nil {
		runtime.HTTPError(ctx, mux, marshaler, w, req, err)
		return
	}

	res := &successResponse{
		Success: true,
		Err:     "no error",
		Code:    "0",
		Data:    resMap,
	}

	resByte, errResByte := json.Marshal(res)

	if errResByte != nil {
		runtime.HTTPError(ctx, mux, marshaler, w, req, errResByte)
		return
	}
	if _, err = w.Write(resByte); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
}

var re = regexp.MustCompile("\\[(.*?)\\]")

// CustomHTTPError for hadling custom message error
func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const fallback = `{"message": "failed to marshal error message", "success": false}`

	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(grpc.Code(err)))
	code := "9999"
	errString := grpc.ErrorDesc(err)
	findErrorCode := re.FindStringSubmatch(errString)

	if len(findErrorCode) > 0 {
		code = findErrorCode[1]
	}
	jErr := json.NewEncoder(w).Encode(errorBody{
		Err:     grpc.ErrorDesc(err),
		Success: false,
		Code:    code,
	})

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}

// customLogger for custome logger
func customLogger(code codes.Code) logrus.Level {
	if code == codes.OK {
		return logrus.InfoLevel
	}

	return logrus.WarnLevel
}

// CreateLogger for creating logger
func CreateLogger() ([]grpc_logrus.Option, *logrus.Entry) {
	logger := &logrus.Logger{}
	logger.SetFormatter(&logrus.TextFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	customFunc := customLogger

	logrusEntry := logrus.NewEntry(logger)

	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(customFunc),
	}

	grpc_logrus.ReplaceGrpcLogger(logrusEntry)

	return opts, logrusEntry
}

// RegisterGRPC for registration GRPC
func RegisterGRPC(whiteList []string, APIPassword string) *grpc.Server {
	opts, logrusEntry := CreateLogger()
	auth := NewAuthInterceptor(APIPassword)
	interceptor := auth.GetUnaryCustom(whiteList)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
				grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
				interceptor,
			),
		),
	)
	return server
}

// RegisterGRPCWithInterceptor for registration GRPC
func RegisterGRPCWithInterceptor(interceptor ...grpc.UnaryServerInterceptor) *grpc.Server {
	opts, logrusEntry := CreateLogger()

	intercep := append(interceptor,
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
	)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				intercep...,
			),
		),
	)
	return server
}
