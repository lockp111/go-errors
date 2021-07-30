package errors

import (
	"errors"
	"fmt"
)

const (
	NoCode          int32 = 0
	InternalCode    int32 = -1
	InternalMessage       = "internal error"
)

var (
	errorsMap = map[int32]*Error{
		NoCode: nil,
		InternalCode: {
			code:    InternalCode,
			message: InternalMessage,
		},
	}
)

type Error struct {
	detail  error
	message string
	code    int32
}

func WithCode(code int32) func(*Error) {
	return func(e *Error) {
		e.code = code
	}
}

func New(msg string, opts ...func(*Error)) *Error {
	e := &Error{
		message: msg,
		code:    NoCode,
	}
	for _, o := range opts {
		o(e)
	}
	return e
}

func Register(code int32, msg string) *Error {
	if _, ok := errorsMap[code]; ok {
		panic(fmt.Sprintf("duplate error code: %d", code))
	}
	e := &Error{
		message: msg,
		code:    code,
	}

	errorsMap[code] = e
	return e
}

func (e *Error) Error() string {
	errMsg := e.message
	if e.code != NoCode && e.code != InternalCode {
		errMsg = fmt.Sprintf("%s(%d)", e.message, e.code)
	}
	if e.detail != nil {
		errMsg = errMsg + ": " + e.detail.Error()
	}
	return errMsg
}

func (e *Error) Code() int32 {
	return e.code
}

func (e *Error) Detail() error {
	return e.detail
}

func (e *Error) WithError(err error) error {
	return &Error{
		detail:  err,
		message: e.message,
		code:    e.code,
	}
}

func (e *Error) WithMessage(msg string) error {
	return &Error{
		detail:  errors.New(msg),
		message: e.message,
		code:    e.code,
	}
}

func (e *Error) Is(err error) bool {
	if err == nil {
		return false
	}

	parseErr := Parse(err)
	switch parseErr.code {
	case InternalCode:
		return errors.Is(e.detail, err)
	case NoCode:
		return parseErr.message == e.message
	default:
		return parseErr.code == e.code
	}
}

func Parse(err error) *Error {
	var target *Error
	if !errors.As(err, &target) {
		target = &Error{
			code:    InternalCode,
			message: InternalMessage,
			detail:  err,
		}
	}
	return target
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func FromCode(code int32) *Error {
	return errorsMap[code]
}
