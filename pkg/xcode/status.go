package xcode

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"strconv"
)

var _ XCode = (*Status)(nil)

type Status struct {
	sts *spb.Status
}

func Error(code Code) *Status {
	return &Status{&spb.Status{Code: code.Code(), Message: code.Message()}}
}

func Errorf(code Code, format string, args ...any) *Status {
	code.msg = fmt.Sprintf(format, args...)
	return Error(code)
}

func (s *Status) Error() string {
	return s.sts.GetMessage()
}

func (s *Status) Code() int32 {
	return s.sts.Code
}

func (s *Status) Message() string {
	if s.sts.Message == "" {
		return strconv.Itoa(int(s.sts.Code))
	}
	return s.sts.Message
}

func (s *Status) Details() []any {
	if s == nil || s.sts == nil {
		return nil
	}
	details := make([]any, 0, len(s.sts.Details))
	for _, d := range s.sts.Details {
		detail := &ptypes.DynamicAny{}
		if err := d.UnmarshalTo(detail); err != nil {
			details = append(details, err)
			continue
		}
		details = append(details, detail.Message)
	}
	return details
}

func (s *Status) WithDetails(msgs ...proto.Message) (*Status, error) {
	for _, msg := range msgs {
		anyMsg, err := anypb.New(msg)
		if err != nil {
			return s, err
		}
		s.sts.Details = append(s.sts.Details, anyMsg)
	}
	return s, nil
}

func (s *Status) Proto() *spb.Status {
	return s.sts
}

func FromCode(code Code) *Status {
	return &Status{sts: &spb.Status{Code: code.code, Message: code.Message()}}
}

func FromProto(pbMsg proto.Message) XCode {
	if msg, ok := pbMsg.(*spb.Status); ok {
		if len(msg.Message) == 0 || msg.Message == strconv.FormatInt(int64(msg.Code), 10) {
			return Code{code: msg.Code}
		}
		return &Status{sts: msg}
	}
	return Errorf(ServerErr, "invalid proto message get %v", pbMsg)
}

func toXCode(grpcStatus *status.Status) Code {
	grpcCode := grpcStatus.Code()
	switch grpcCode {
	case codes.OK:
		return OK
	case codes.InvalidArgument:
		return NotFound
	case codes.PermissionDenied:
		return AccessDenied
	case codes.Unauthenticated:
		return Unauthorized
	case codes.Unimplemented:
		return MethodNotAllowed
	case codes.DeadlineExceeded:
		return Deadline
	case codes.Unavailable:
		return ServiceUnavailable
	case codes.Unknown:
		return String(grpcStatus.Message())
	}
	return ServerErr
}
func CodeFromError(err error) XCode {
	err = errors.Cause(err)
	if code, ok := err.(XCode); ok {
		return code
	}

	switch err {
	case context.Canceled:
		return Canceled
	case context.DeadlineExceeded:
		return Deadline
	}
	return ServerErr
}

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

func GrpcStatusToXCode(gstatus *status.Status) XCode {
	details := gstatus.Details()
	for i := len(details) - 1; i >= 0; i-- {
		detail := details[i]
		if pb, ok := detail.(proto.Message); ok {
			return FromProto(pb)
		}
	}
	return toXCode(gstatus)
}
