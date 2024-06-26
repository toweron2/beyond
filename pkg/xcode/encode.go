package xcode

import (
	"context"
	"github.com/pkg/errors"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"strconv"
)

// FromError 将自定义的业务错误转换为grpc status
func FromError(err error) *status.Status {
	err = errors.Cause(err)
	if code, ok := err.(XCode); ok {
		grpcStatus, e := gRpcStatusFromXcode(code)
		if e == nil {
			return grpcStatus
		}
	}
	var grpcStatus *status.Status
	switch err {
	case context.Canceled:
		grpcStatus, _ = gRpcStatusFromXcode(Canceled)
	case context.DeadlineExceeded:
		grpcStatus, _ = gRpcStatusFromXcode(Deadline)
	default:
		grpcStatus, _ = status.FromError(err)
	}

	return grpcStatus
}

// gRpcStatusFromXcode 通过WithDetails方法将自定义业务错误码存放到detail中
func gRpcStatusFromXcode(code XCode) (*status.Status, error) {
	var sts *Status
	switch v := code.(type) {
	case *Status:
		sts = v
	case Code:
		sts = FromCode(v)
	default:
		sts = Error(Code{code.Code(), code.Message()})
		for _, detail := range code.Details() {
			if msg, ok := detail.(proto.Message); ok {
				_, _ = sts.WithDetails(msg)
			}
		}
	}

	stas := status.New(codes.Unknown, strconv.Itoa(int(sts.Code())))
	return stas.WithDetails(sts.Proto())
}

func FromCode(code Code) *Status {
	return &Status{sts: &spb.Status{Code: code.code, Message: code.Message()}}
}
