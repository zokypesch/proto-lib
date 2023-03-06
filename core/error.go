package core

import (
	"google.golang.org/grpc/status"
)

type ErrorInterface interface {
	GetData() interface{}
	Code() string
	Error() string
	GRPCStatus() *status.Status
}
