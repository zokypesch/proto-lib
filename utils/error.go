package utils

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HTTPBadRequest for bad request error
func HTTPBadRequest(code string, message string) error {
	return status.Errorf(codes.InvalidArgument, fmt.Sprintf("[%s] %s", code, message))
}

// HTTPUnAuthorize for bad request error
func HTTPUnAuthorize(code string, message string) error {
	return status.Errorf(codes.PermissionDenied, fmt.Sprintf("[%s] %s", code, message))
}

// HTTPNotFound for bad request error
func HTTPNotFound(code string, message string) error {
	return status.Errorf(codes.NotFound, fmt.Sprintf("[%s] %s", code, message))
}
