package response

import (
	"github.com/gin-gonic/gin"
)

type successMessage struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
}

type errorMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func SuccessMessageJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, successMessage{
		Success: true,
		Data:    data,
	})
}

func ErrorMessageJSON(c *gin.Context, status int, msg string) {
	c.AbortWithStatusJSON(status, errorMessage{
		Success: false,
		Message: msg,
	})
}
