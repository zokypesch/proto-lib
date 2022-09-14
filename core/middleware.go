package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	proto "github.com/golang/protobuf/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	logrus "github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmgrpc"
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
	if md, ok := runtime.ServerMetadataFromContext(ctx); ok {
		vals := md.HeaderMD.Get("x-request-id")
		if len(vals) > 0 {
			w.Header().Set("X-Request-Id", vals[0])
		}
		vals = md.HeaderMD.Get("x-service-id")
		if len(vals) > 0 {
			w.Header().Set("X-Service-Id", vals[0])
		}
	}

	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PUT, DELETE, PATCH")

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
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PUT, DELETE, PATCH")

	w.WriteHeader(runtime.HTTPStatusFromCode(grpc.Code(err)))
	code := "9999"
	errString := grpc.ErrorDesc(err)
	findErrorCode := re.FindStringSubmatch(errString)

	errWithoutCode := errString
	if len(findErrorCode) > 0 {
		code = findErrorCode[1]
		errWithoutCode = strings.TrimSpace(strings.Replace(errString, fmt.Sprintf("[%s]", code), "", -1))
	}
	jErr := json.NewEncoder(w).Encode(errorBody{
		Err:     errWithoutCode,
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

// AppendInterceptor for default register
func AppendInterceptor(interceptor ...grpc.UnaryServerInterceptor) []grpc.UnaryServerInterceptor {

	// server := grpc.NewServer(
	// 	grpc.UnaryInterceptor(
	// 		grpc_middleware.ChainUnaryServer(
	// 			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
	// 			grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
	// 			interceptor,
	// 		),
	// 	),
	// )
	opts, logrusEntry := CreateLogger()

	intercep := append(interceptor,
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
	)

	return intercep
}

// RegisterGRPC for registration GRPC
func RegisterGRPC(srvName string, whiteList []string, APIPassword string) *grpc.Server {
	auth := NewAuthInterceptor(APIPassword)
	interceptor := auth.GetUnaryCustom(whiteList)

	return RegisterGRPCWithInterceptor(srvName, interceptor)
}

var (
	// Create a metrics registry.
	reg = prometheus.NewRegistry()

	serviceName = ""
	// Create some standard server metrics.
	grpcMetrics = grpc_prometheus.NewServerMetrics()

	// Create a customized counter metric.
	customizedCounterMetric *prometheus.CounterVec
)

func registerCustomizeMetrics(svc string) {
	customizedCounterMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: svc,
		Help: "Total number of RPCs handled on the server.",
	}, []string{"name"})
}

func initProm() {
	// Register standard server metrics and customized metrics to registry.
	reg.MustRegister(grpcMetrics, customizedCounterMetric)
	customizedCounterMetric.WithLabelValues(serviceName)
}

// RegisterGRPCWithPrometh for get unnary prometheus
func RegisterGRPCWithPrometh(srvName string, interceptor ...grpc.UnaryServerInterceptor) *grpc.Server {
	registerCustomizeMetrics(srvName)

	serviceName = srvName + "_counter"

	newIntercep := append(interceptor,
		grpcMetrics.UnaryServerInterceptor(),
		GetUnaryCounter(serviceName),
	)
	intercep := AppendInterceptor(newIntercep...)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				intercep...,
			),
		),
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
	)
	return server

}

// https://www.elastic.co/guide/en/apm/agent/go/master/builtin-modules.html#builtin-modules-apmgrpc
// RegisterGRPCWithPromethAndAPM with ap agent too
func RegisterGRPCWithPromethAndAPM(srvName string, interceptor ...grpc.UnaryServerInterceptor) *grpc.Server {
	registerCustomizeMetrics(srvName)

	serviceName = srvName + "_counter"

	newIntercep := append(interceptor,
		grpcMetrics.UnaryServerInterceptor(),
		GetUnaryCounter(serviceName),
		apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery()),
	)
	intercep := AppendInterceptor(newIntercep...)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				intercep...,
			),
		),
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
	)
	return server

}

// RegisterPrometheus for registration prometheus
func RegisterPrometheus(server *grpc.Server, port int64) {
	initProm()

	// Initialize all metrics.
	grpcMetrics.InitializeMetrics(server)

	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", port)}

	// Start your http server for prometheus.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

}

// RegisterGRPCWithInterceptor for registration GRPC
func RegisterGRPCWithInterceptor(srvName string, interceptor ...grpc.UnaryServerInterceptor) *grpc.Server {
	opts, logrusEntry := CreateLogger()
	registerCustomizeMetrics(srvName)

	serviceName = srvName + "_counter"

	intercep := append(interceptor,
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
		grpcMetrics.UnaryServerInterceptor(),
		GetUnaryCounter(serviceName),
	)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				intercep...,
			),
		),
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
	)
	return server
}

// RegisterGRPCWithARM for registration GRPC With ARM
func RegisterGRPCWithARM(srvName string, interceptor ...grpc.UnaryServerInterceptor) *grpc.Server {
	opts, logrusEntry := CreateLogger()
	registerCustomizeMetrics(srvName)

	serviceName = srvName + "_counter"

	intercep := append(interceptor,
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_logrus.UnaryServerInterceptor(logrusEntry, opts...),
		apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery()),
	)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				intercep...,
			),
		),
		grpc.StreamInterceptor(grpcMetrics.StreamServerInterceptor()),
	)
	return server
}
