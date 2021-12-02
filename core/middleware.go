package core

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"regexp"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	runtime "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	logrus "github.com/sirupsen/logrus"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/reflect/protoreflect"
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
func LocalForward(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, resp protoreflect.ProtoMessage, opts ...func(context.Context, http.ResponseWriter, protoreflect.ProtoMessage) error) {
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
func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
	const fallback = `{"message": "failed to marshal error message", "success": false}`

	// no longer needed in new version v2
	// w.Header().Set("Content-type", marshaler.ContentType())
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, POST, GET, PUT, DELETE, PATCH")

	httpStatusFromCode := runtime.HTTPStatusFromCode(grpc.Code(err))
	w.WriteHeader(httpStatusFromCode)
	code := "9999"
	errString := grpc.ErrorDesc(err)
	findErrorCode := re.FindStringSubmatch(errString)

	errWithoutCode := errString
	if len(findErrorCode) > 0 {
		code = findErrorCode[1]
		errWithoutCode = strings.TrimSpace(strings.Replace(errString, fmt.Sprintf("[%s]", code), "", -1))
	}
	respBody := errorBody{
		Err:     errWithoutCode,
		Success: false,
		Code:    code,
	}
	jErr := json.NewEncoder(w).Encode(respBody)

	Logs.ReportCaller = false
	toLog := AddCtxToLog(ctx, Logs)
	toLog = toLog.WithFields(logrus.Fields{
		"request": func() string {
			body, _ := httputil.DumpRequest(req, true)
			return string(body)
		}(),
		"response": fmt.Sprintf("%s", map[string]interface{}{
			"statusCode": httpStatusFromCode,
			"content":    respBody,
			"errCode":    code,
		}),
	})
	toLog.Error(err.Error())

	if jErr != nil {
		w.Write([]byte(fallback))
	}
}

// customLogger for custome logger
func customLogger(code codes.Code) logrus.Level {
	if code == codes.OK {
		Logs.SetOutput(os.Stdout)
		return logrus.InfoLevel
	}

	Logs.SetOutput(os.Stderr)
	return logrus.ErrorLevel
}

// DefaultMessageProducer writes the default message
func CustomMessageProducer(ctx context.Context, format string, level logrus.Level, code codes.Code, err error, fields logrus.Fields) {
	if err != nil {
		fields[logrus.ErrorKey] = err
	}
	entry := ctxlogrus.Extract(ctx).WithContext(ctx).WithFields(fields)
	switch level {
	case logrus.DebugLevel:
		entry.Debugf(format)
	case logrus.InfoLevel:
		entry.Infof(format)
	case logrus.WarnLevel:
		entry.Warningf(format)
	case logrus.ErrorLevel:
		entry.Errorf(format)
	case logrus.FatalLevel:
		entry.Fatalf(format)
	case logrus.PanicLevel:
		entry.Panicf(format)
	}
}

// CreateLogger for creating logger
func CreateLogger() ([]grpc_logrus.Option, *logrus.Entry) {
	logger := Logs
	logger.ReportCaller = false
	customFunc := customLogger
	customMessage := CustomMessageProducer

	logrusEntry := logrus.NewEntry(logger)

	opts := []grpc_logrus.Option{
		grpc_logrus.WithLevels(customFunc),
		grpc_logrus.WithMessageProducer(customMessage),
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

func GetPrometheusHandler(opts promhttp.HandlerOpts) http.Handler {
	return promhttp.HandlerFor(reg, opts)
}

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
	RegisterPrometheusWithoutServer(server)

	httpServer := &http.Server{Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), Addr: fmt.Sprintf("0.0.0.0:%d", port)}

	// Start your http server for prometheus.
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal("Unable to start a http server.")
		}
	}()

}

func RegisterPrometheusWithoutServer(server *grpc.Server) {
	initProm()

	// Initialize all metrics.
	grpcMetrics.InitializeMetrics(server)
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
		grpcMetrics.UnaryServerInterceptor(),
		grpc_logrus.PayloadUnaryServerInterceptor(
			logrusEntry,
			func(ctx context.Context, fullMethodName string, servingObject interface{}) bool {
				return true
			},
		),
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
