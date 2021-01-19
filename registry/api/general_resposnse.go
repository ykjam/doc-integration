package api

import (
	"errors"
)

type GeneralResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseErrorCodeAndMessage struct {
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_msg"`
}

const (
	ErrorCodeOK                  int = 200
	ErrorCodeBadRequest          int = 400
	ErrorCodeForbidden           int = 403
	ErrorCodeNotFound            int = 404
	ErrorCodeInternalServerError int = 500
)

var ErrOK = errors.New("OK")
var ErrBadRequest = errors.New("bad request")
var ErrNotFound = errors.New("not found")
var ErrInternalServerError = errors.New("internal server error")

const (
	ErrorMessageOK                  = "ok"
	ErrorMessageBadRequest          = "bad_request"
	ErrorMessageNotFound            = "not_found"
	ErrorMessageInternalServerError = "internal_server_error"
)
