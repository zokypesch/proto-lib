package core

import (
	"google.golang.org/grpc/status"
)

type ErrorInterface interface {
	GetData() interface{}
	Error() string
	GRPCStatus() *status.Status
}
