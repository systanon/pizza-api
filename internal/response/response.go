package response

import "github.com/gin-gonic/gin"

type successResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Success bool      `json:"success"`
	Error   errorBody `json:"error"`
}

func RespondSuccess(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, successResponse{Success: true, Data: data, Message: message})
}

func RespondError(c *gin.Context, status int, code string, message string) {
	c.JSON(status, errorResponse{
		Success: false,
		Error:   errorBody{Code: code, Message: message},
	})
}
