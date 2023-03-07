package core

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/zokypesch/proto-lib/grpc/pb/protolib"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorInterface interface {
	Code() string
	Error() string
	GRPCStatus() *status.Status
	GetData() interface{}
}

type CustomErrDataCode struct {
	err  error
	data interface{}
}

func NewCustomErrDataCode(err error, data interface{}) *CustomErrDataCode {
	return &CustomErrDataCode{err: err, data: data}
}

func (c CustomErrDataCode) Error() string {
	return c.err.Error()
}

func (c CustomErrDataCode) GRPCStatus() *status.Status {
	if se, ok := c.err.(interface {
		GRPCStatus() *status.Status
	}); ok {
		if message, ok := c.data.(proto.Message); ok {
			gprcStatus, err := se.GRPCStatus().WithDetails(message)
			if err != nil {
				return status.New(codes.Internal, fmt.Sprintf("[ERRPTL999] %s", err.Error()))
			}
			return gprcStatus
		}

		dataBytes, err := json.Marshal(c.data)
		if err != nil {
			return status.New(codes.Internal, fmt.Sprintf("[ERRPTL999] %s", err.Error()))
		}
		dataMsg := protolib.DataMessage{Data: string(dataBytes)}
		gprcStatus, err := se.GRPCStatus().WithDetails(&dataMsg)
		if err != nil {
			return status.New(codes.Internal, fmt.Sprintf("[ERRPTL999] %s", err.Error()))
		}
		return gprcStatus
	}

	return status.New(codes.Unknown, fmt.Sprintf("[ERRPTL999] %s", c.err.Error()))
}

func (c CustomErrDataCode) Code() string {
	if se, ok := c.err.(interface {
		Code() string
	}); ok {
		return se.Code()
	}

	return ""
}
