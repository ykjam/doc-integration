package web

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"ykjam/doc-registry-go/api"
)

func GetRemoteAddress(r *http.Request) string {
	if val := r.Header.Get("X-Forwarded-For"); val != "" {
		return strings.Split(val, ":")[0]
	} else if val := r.Header.Get("X-Real-IP"); val != "" {
		return strings.Split(val, ":")[0]
	} else {
		return strings.Split(r.RemoteAddr, ":")[0]
	}
}

type Server struct {
	c        *api.APIController
}

type httpPostWithLog func(ctx context.Context, w http.ResponseWriter, r *http.Request, clog *log.Entry)

func NewServer(apiController *api.APIController) *Server {
	return &Server{
		c:        apiController,
	}
}

func (s *Server) handleHttpPostOrGetWithLog(handleName string, w http.ResponseWriter, r *http.Request, f httpPostWithLog) {
	ctx := r.Context()
	clog := log.WithFields(log.Fields{
		"remote-addr": GetRemoteAddress(r),
		"uri":         r.RequestURI,
		"method":      r.Method,
		"handle":      handleName,
	}).WithContext(ctx)
	if r.Method == http.MethodPost || r.Method == http.MethodGet {
		f(ctx, w, r, clog)
	} else {
		clog.Error("invalid request, method not allowed")
		s.sendResponseByCode(w, api.ErrorCodeForbidden, clog)
	}
}

func (s *Server) sendResponseByCode(w http.ResponseWriter, errCode int, clog *log.Entry) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(errCode)
	var errMessage string
	switch errCode {
	case api.ErrorCodeOK:
		errMessage = api.ErrorMessageOK
	case api.ErrorCodeTfaRequired:
		errMessage = api.ErrorMessageTfaRequired
	case api.ErrorCodeBadRequest:
		errMessage = api.ErrorMessageBadRequest
	case api.ErrorCodeFileSizeTooLarge:
		errMessage = api.ErrorMessageFileSizeTooLarge
	case api.ErrorCodeUnauthorized:
		errMessage = api.ErrorMessageUnauthorized
	case api.ErrorCodeForbidden:
		errMessage = api.ErrorMessageForbidden
	case api.ErrorCodeNotFound:
		errMessage = api.ErrorMessageNotFound
	case api.ErrorCodeExpired:
		errMessage = api.ErrorMessageExpired
	case api.ErrorCodeConflict:
		errMessage = api.ErrorMessageConflict
	case api.ErrorCodeTooManyRequests:
		errMessage = api.ErrorMessageTooManyRequests
	case api.ErrorCodeInternalServerError:
		errMessage = api.ErrorMessageInternalServerError
	}
	var resp api.GeneralResponse
	if errCode == api.ErrorCodeOK {
		resp = api.GeneralResponse{
			Success: true,
		}
	} else {
		resp = api.GeneralResponse{
			Success: false,
			Data: api.ResponseErrorCodeAndMessage{
				ErrorCode:    errCode,
				ErrorMessage: errMessage,
			},
		}
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		clog.WithError(err).Error(fmt.Sprint(" data: ", resp))
		http.Error(w, api.ErrorMessageInternalServerError, api.ErrorCodeInternalServerError)
	}
}

func (s *Server) sendResponseByError(w http.ResponseWriter, err error, clog *log.Entry) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var errCode int
	var errMessage string

	switch err {
	case api.ErrInternalServerError:
		errCode = api.ErrorCodeInternalServerError
		errMessage = api.ErrorMessageInternalServerError
	case api.ErrBadRequest:
		errCode = api.ErrorCodeBadRequest
		errMessage = api.ErrorMessageBadRequest
	case api.ErrFileSizeTooLarge:
		errCode = api.ErrorCodeFileSizeTooLarge
		errMessage = api.ErrorMessageFileSizeTooLarge
	case api.ErrNotFound:
		errCode = api.ErrorCodeNotFound
		errMessage = api.ErrorMessageNotFound
	case api.ErrForbidden:
		errCode = api.ErrorCodeForbidden
		errMessage = api.ErrorMessageForbidden
	case api.ErrUnauthorized:
		errCode = api.ErrorCodeUnauthorized
		errMessage = api.ErrorMessageUnauthorized
	case api.ErrExpired:
		errCode = api.ErrorCodeExpired
		errMessage = api.ErrorMessageExpired
	case api.ErrConflict:
		errCode = api.ErrorCodeConflict
		errMessage = api.ErrorMessageConflict
	}
	w.WriteHeader(errCode)
	resp := api.GeneralResponse{
		Success: false,
		Data: api.ResponseErrorCodeAndMessage{
			ErrorCode:    errCode,
			ErrorMessage: errMessage,
		},
	}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		clog.WithError(err).Error(fmt.Sprint(" data: ", resp))
		http.Error(w, api.ErrorMessageInternalServerError, api.ErrorCodeInternalServerError)
	}
}

func (s *Server) sendResponseOKWithData(w http.ResponseWriter, data interface{}, clog *log.Entry) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(api.ErrorCodeOK)
	resp := api.GeneralResponse{
		Success: true,
		Data:    data,
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		clog.WithError(err).Error(fmt.Sprint(" data: ", resp))
		http.Error(w, api.ErrorMessageInternalServerError, api.ErrorCodeInternalServerError)
	}
}
