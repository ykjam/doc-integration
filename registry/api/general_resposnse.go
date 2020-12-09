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
	ErrorCodeTfaRequired         int = 250
	ErrorCodeBadRequest          int = 400
	ErrorCodeUnauthorized        int = 401
	ErrorCodeForbidden           int = 403
	ErrorCodeNotFound            int = 404
	ErrorCodeConflict            int = 409
	ErrorCodeExpired             int = 408
	ErrorCodeFileSizeTooLarge    int = 413
	ErrorCodeTooManyRequests     int = 429
	ErrorCodeInternalServerError int = 500
)

var ErrOK = errors.New("OK")
var ErrBadRequest = errors.New("Bad request")
var ErrUnauthorized = errors.New("Unauthorized")
var ErrForbidden = errors.New("Forbidden")
var ErrNotFound = errors.New("Not found")
var ErrConflict = errors.New("Conflict")
var ErrExpired = errors.New("Expired")
var ErrFileSizeTooLarge = errors.New("File size too large")
var ErrInternalServerError = errors.New("Internal server error")

const (
	ErrorMessageOK                  = "ok"
	ErrorMessageTfaRequired         = "tfa_required"
	ErrorMessageBadRequest          = "bad_request"
	ErrorMessageUnauthorized        = "unauthorized"
	ErrorMessageForbidden           = "forbidden"
	ErrorMessageNotFound            = "not_found"
	ErrorMessageConflict            = "conflict"
	ErrorMessageExpired             = "expired"
	ErrorMessageFileSizeTooLarge    = "file_size_too_large"
	ErrorMessageTooManyRequests     = "otp_retry_limit_exceeded"
	ErrorMessageInternalServerError = "internal_server_error"
)
