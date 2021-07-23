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
	errorsMap = map[int32]Error{
		NoCode: {},
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

func New(msg string) *Error {
	return &Error{
		message: msg,
		code:    NoCode,
		detail:  errors.New(msg),
	}
}

func Register(code int32, msg string) *Error {
	if _, ok := errorsMap[code]; ok {
		panic(fmt.Sprintf("duplate error code: %d", code))
	}
	return &Error{
		message: msg,
		code:    code,
		detail:  errors.New(msg),
	}
}

func (e *Error) Error() string {
	return e.detail.Error()
}

func (e *Error) Code() int32 {
	return e.code
}

func (e *Error) Detail() error {
	return e.detail
}

func (e *Error) WithError(err error) error {
	return &Error{
		detail:  fmt.Errorf("%s: %w", e.message, err),
		message: e.message,
		code:    e.code,
	}
}

func (e *Error) WithMessage(msg string) error {
	return &Error{
		detail:  fmt.Errorf("%w: %s", e, msg),
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

func Parse(err error) (target *Error) {
	if !errors.As(err, &target) {
		target = &Error{
			code:    InternalCode,
			message: InternalMessage,
			detail:  err,
		}
	}
	return
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}
