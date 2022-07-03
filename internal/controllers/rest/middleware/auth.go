package middleware

import (
	"net/http"

	"social-network/internal/controllers/rest/response"
	"social-network/internal/service/auth"

	"github.com/gin-gonic/gin"
)

func CheckAccess(service auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorMessageJSON(c, http.StatusUnauthorized, response.EmptyToken)
			return
		}

		err := service.CheckToken(authHeader)
		if err != nil {
			response.ErrorMessageJSON(c, http.StatusUnauthorized, response.InvalidToken)
			return
		}

		return
	}
}
