package uerror

import (
	"errors"
)

var (
	ErrUnknown    = errors.New("unknown error")     // unknown error
	ErrNil        = errors.New("object is nil")     // nil
	ErrUnexpected = errors.New("un expected")       // un expected
	ErrExist      = errors.New("already exists")    // already exist
	ErrNotExist   = errors.New("does not exist")    // not exist
	ErrInvalid    = errors.New("invalid argument")  // invalid argument
	ErrPermission = errors.New("permission denied") // permission denied
	ErrClosed     = errors.New("already closed")    // already closed
	ErrTimeout    = errors.New("timeout")           // timeout
	ErrIO         = errors.New("i/o error")         // i/o error
	ErrTooLong    = errors.New("buffer too long")   // buffer too long
	ErrTooShort   = errors.New("buffer too short")  // buffer too short
)
