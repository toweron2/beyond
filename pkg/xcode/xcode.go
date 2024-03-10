package xcode

import "strconv"

type XCode interface {
	Error() string
	Code() int32
	Message() string
	Details() []any
}

type Code struct {
	code int32
	msg  string
}

func (c Code) Error() string {
	if len(c.msg) > 0 {
		return c.msg
	}
	return strconv.Itoa(int(c.code))
}

func (c Code) Code() int32 {
	return c.code
}

func (c Code) Message() string {
	return c.Error()
}

func (c Code) Details() []any {
	return nil
}

func String(s string) Code {
	if len(s) == 0 {
		return OK
	}
	code, err := strconv.Atoi(s)
	if err != nil {
		return ServerErr
	}
	return Code{code: int32(code)}
}

func New(code int32, msg string) Code {
	return Code{code, msg}
}

func add(code int32, msg string) Code {
	return Code{code, msg}
}
