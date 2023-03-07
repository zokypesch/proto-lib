package core_test

import (
	"encoding/json"
	"fmt"
	"github.com/zokypesch/proto-lib/core"
	"github.com/zokypesch/proto-lib/grpc/pb/protolib"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"testing"
)

type DummyError struct {
}

func (d DummyError) Code() string {
	return "ERR999"
}

func (d DummyError) Error() string {
	return fmt.Sprintf("[%s] something went wrong", d.Code())
}

func (d DummyError) GRPCStatus() *status.Status {
	return status.New(codes.Internal, d.Error())
}

func (d DummyError) GetData() interface{} {
	return nil
}

func TestGRPCStatus(t *testing.T) {
	type args struct {
		param error
		data  interface{}
	}

	var testData = []struct {
		name     string
		args     args
		expected string
	}{
		{
			name: "case1 interface",
			args: args{
				param: DummyError{},
				data: map[string]interface{}{
					"a": "b",
					"c": []int{1, 2, 3},
				},
			},
			expected: `[{"Data":"{\"a\":\"b\",\"c\":[1,2,3]}"}]`,
		},
		{
			name: "case2 pb struct",
			args: args{
				param: DummyError{},
				data: &protolib.Dummy{
					Field1: "field1",
					Field2: 2,
					Field3: []*protolib.DummyInner{
						{Name: "febri"},
						{Name: "wong"},
					},
				},
			},
			expected: `[{"Field1":"field1","Field2":2,"Field3":[{"name":"febri"},{"name":"wong"}]}]`,
		},
	}
	for _, tt := range testData {
		t.Run(tt.name, func(t *testing.T) {
			code := core.NewCustomErrDataCode(tt.args.param, tt.args.data)

			grpcStatus := code.GRPCStatus()

			log.Println(code, grpcStatus)

			details := grpcStatus.Details()
			b, err := json.Marshal(details)
			if err != nil {
				t.Error(err)
				return
			}
			bStr := string(b)
			if tt.expected != bStr {
				t.Errorf("expect %s got %s", tt.expected, bStr)
				return
			}
		})
	}
}
