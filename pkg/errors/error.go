package errors

import (
	"errors"
)

var (
	ErrUnknownType       = errors.New("storage type is unknown or  not set")
	ErrMaxAttemptsExceed = errors.New("render page attempts exceeded")
	ErrPageIsNil         = errors.New("page data is nil")
	ErrUnknownTrigger    = errors.New("don't know what to wait")
	ErrLinkHasBaseURL    = errors.New("link contains base url")
	ErrNoBaseURL         = errors.New("base_url is not set in config")
	ErrWrongLookupType   = errors.New("lookup type is wrong or not set")
)
