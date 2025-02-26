package server

import "errors"

var (
	ErrInvalidRequest       = errors.New("incorrect request")
	ErrUnknownFunction      = errors.New("unknown function")
	ErrDeviceNotFound       = errors.New("device not found")
	ErrInvalidAddress       = errors.New("invalid address")
	ErrMsgServerFail        = "message server fail"
	ErrMsgServerNotStarted  = "mpdbus message server not started"
	ErrMsgServerFailConnect = "server fail connect"
	ErrMsgServerFailRead    = "server fail read"
	ErrMsgFailLoadConfig    = "fail load config"
	ErrMsgCovertConfig      = "fail convert config"
)
