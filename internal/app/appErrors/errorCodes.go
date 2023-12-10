package apperrors

import (
	"main/internal/services/logger"
	"net/http"
)

type AppError struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	Details        string `json:"details"`
	HttpStatusCode int    `json:"-"`
}

type AppErrorCodeMap map[string]AppError

var AppErrors = AppErrorCodeMap{
	"E001": AppError{

		Code:           "E001",
		Message:        "Internal server error",
		Details:        "Something went wrong on our side. Please contact us",
		HttpStatusCode: http.StatusInternalServerError,
	},
	"E002": AppError{
		Code:           "E002",
		Message:        "Service unavailable",
		Details:        "Service is currently unavailable. Please try again later",
		HttpStatusCode: http.StatusServiceUnavailable,
	},
	"E003": AppError{
		Code:           "E003",
		Message:        "Invalid request body",
		Details:        "The request body is invalid",
		HttpStatusCode: http.StatusBadRequest,
	},
	"E004": AppError{
		Code:           "E004",
		Message:        "Not found",
		Details:        "The requested resource was not found",
		HttpStatusCode: http.StatusNotFound,
	},
	"E005": AppError{
		Code:           "E005",
		Message:        "Unauthorized",
		Details:        "Missing or invalid authorization data",
		HttpStatusCode: http.StatusUnauthorized,
	},
	"E006": AppError{
		Code:           "E006",
		Message:        "Forbidden",
		Details:        "Credentials not active",
		HttpStatusCode: http.StatusForbidden,
	},
}

type ErrorResponse struct {
	Errors         []AppError `json:"errors"`
	HttpStatusCode int        `json:"-"`
}

func New(code string) ErrorResponse {
	appErr, ok := AppErrors[code]

	if !ok {
		appErr = AppErrors["E001"]
		logger.AppMiscLog.Errorw("Error code not found", "code", code)
	}

	return ErrorResponse{
		Errors: []AppError{
			appErr,
		},
		HttpStatusCode: appErr.HttpStatusCode,
	}
}
