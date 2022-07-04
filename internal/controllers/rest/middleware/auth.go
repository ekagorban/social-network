package middleware

import (
	"errors"
	"log"
	"net/http"

	"social-network/internal/controllers/rest/response"
	"social-network/internal/errapp"
	"social-network/internal/service/auth"

	"github.com/gin-gonic/gin"
)

func CheckAccess(service auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorMessageJSON(c, http.StatusUnauthorized, errapp.EmptyToken.Error())
			return
		}

		err := service.CheckToken(authHeader)
		if err != nil {
			log.Printf("service.CheckToken error: %v", err)

			if errors.Is(err, errapp.InvalidToken) {
				response.ErrorMessageJSON(c, http.StatusUnauthorized, errapp.InvalidToken.Error())
				return
			}

			response.ErrorMessageJSON(c, http.StatusUnauthorized, response.InternalError)
			return
		}

		return
	}
}
