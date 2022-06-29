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
			response.ErrorMessageJSON(c, http.StatusUnauthorized, "empty")
			return
		}
		allow, err := service.CheckToken(authHeader)
		if err != nil {
			response.ErrorMessageJSON(c, http.StatusUnauthorized, err.Error())
			return
		}
		if !allow {
			response.ErrorMessageJSON(c, http.StatusUnauthorized, "not allow")
			return
		}

		return
	}
}
