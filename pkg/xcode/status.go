package xcode

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	spb "google.golang.org/genproto/googleapis/rpc/status"
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
