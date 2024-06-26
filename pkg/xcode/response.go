package xcode

import (
	"context"
	"github.com/pkg/errors"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"net/http"
)

func ErrHandler(err error) (int, any) {
	code := CodeFromError(err)

	return http.StatusOK, spb.Status{
		Code:    code.Code(),
		Message: code.Message(),
	}
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
