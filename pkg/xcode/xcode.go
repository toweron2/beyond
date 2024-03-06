package xcode

import "strconv"

type XCode interface {
    Error() string
    Code() int
    Message() string
    Details() []any
}

type Code struct {
    code int
    msg  string
}

func New(code int, msg string) Code {
    return Code{code, msg}
}

func add(code int, msg string) Code {
    return Code{code, msg}
}

func (c Code) Error() string {
    if len(c.msg) > 0 {
        return c.msg
    }
    return strconv.Itoa(c.code)
}

func (c Code) Code() int {
    return c.code
}

func (c Code) Message() string {
    return c.Error()
}

func (c Code) Details() []any {
    return nil
}
