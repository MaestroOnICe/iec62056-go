package iec62056

import "errors"

var (
	ErrInvalidFormat      = errors.New("invalid message format")
	ErrMissingStartChar   = errors.New("missing start character")
	ErrMissingEndChar     = errors.New("missing end character")
	ErrInvalidDataBlock   = errors.New("invalid data block format")
	ErrEmptyMessage       = errors.New("empty message")
	ErrInvalidBaudRate    = errors.New("invalid baud rate")
	ErrInvalidMode        = errors.New("invalid protocol mode")
	ErrAuthenticationFail = errors.New("authentication failed")
	ErrPortClosed         = errors.New("serial port is closed")
	ErrTimeout            = errors.New("operation timeout")
)
