package xcode

import (
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
