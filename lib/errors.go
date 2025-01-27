package lib

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const (
	ErrorCodeUnknown              = 11
	ErrorCodeCannotLaunchChrome   = 12
	ErrorCodeNameNotResolved      = 13
	ErrorCodeNameResolutionFailed = 14
	ErrorCodeTimeout              = 15
	ErrorCodeAddressUnreachable   = 18
	ErrorCodeEmptyResponse        = 19
	ErrorCodeCannotSaveFile       = 20
	ErrorCodeOther                = 21
)

var (
	errorMessageCodeMapping = map[string]int{
		"net::ERR_NAME_NOT_RESOLVED":                    13,
		"net::ERR_NAME_RESOLUTION_FAILED":               14,
		"Navigation Timeout Exceeded: 30000ms exceeded": 15,
		"Navigation Timeout Exceeded: 60000ms exceeded": 15,
		"net::ERR_CERT_DATE_INVALID":                    16,
		"PAGE_RETURNED_HTTP_302":                        17,
		"net::ERR_ADDRESS_UNREACHABLE":                  18,
		"net::ERR_EMPTY_RESPONSE":                       19,
	}
)

type Error struct {
	message string
	code    int
	err     error
}

func NewError(message string, code int, err error) *Error {
	return &Error{
		message: message,
		code:    code,
		err:     err,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.message, e.err)
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Code() int {
	return e.code
}

var _ error = &Error{}

func GetCodeFromError(err error) int {
	if errors.Is(err, context.DeadlineExceeded) {
		return ErrorCodeTimeout
	}
	message := err.Error()
	for pattern, code := range errorMessageCodeMapping {
		if strings.HasSuffix(message, pattern) {
			return code
		}
	}
	return 0
}
