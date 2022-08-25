package v2

import (
	"context"
	"encoding/json"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/protobuf/reflect/protoreflect"
	"net/http"
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