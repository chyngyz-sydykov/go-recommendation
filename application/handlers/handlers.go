package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/chyngyz-sydykov/go-recommendation/infrastructure/logger"
)

type CommonHandlerInterface interface {
	HandleError(w http.ResponseWriter, err error, statusCode int, errorResponse ErrorResponse)
	LogError(err error, statusCode int)
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type CommonHandler struct {
	logger logger.LoggerInterface
}

const INVALID_REQUEST string = "INVALID_REQUEST"
const RESOURCE_NOT_FOUND string = "RESOURCE_NOT_FOUND"
const SERVER_ERROR string = "SERVER_ERROR"
const SERVICE_NOT_AVAILABLE string = "SERVICE_NOT_AVAILABLE"

func NewCommonHandler(logger logger.LoggerInterface) *CommonHandler {

	return &CommonHandler{logger: logger}
}

func (c *CommonHandler) HandleError(w http.ResponseWriter, err error, statusCode int, errorResponse ErrorResponse) {
	c.logger.LogError(statusCode, err)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

func (c *CommonHandler) LogError(err error, statusCode int) {
	c.logger.LogError(statusCode, err)
}
