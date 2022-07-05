package routes

import (
	"context"
	"errors"
	"log"
	"net/http"

	"social-network/internal/controllers/rest/response"
	"social-network/internal/errapp"
	"social-network/internal/service/auth"

	"github.com/gin-gonic/gin"
)

/*
[post] /signin - sign in with login and password (return token and user id)
*/

func SignIn(r *gin.Engine, service auth.Service) {
	v1 := r.Group("/v1")
	v1.POST("/signin", signIn(service))
}

func signIn(service auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(c)
		defer cancel()

		var dataSignIn auth.SignInData
		if err := c.ShouldBindJSON(&dataSignIn); err != nil {
			log.Println(err)
			response.ErrorMessageJSON(c, http.StatusBadRequest, err.Error())
			return
		}

		data, err := service.SignIn(ctx, dataSignIn)
		if err != nil {
			log.Printf("service.SignIn error: %v", err)

			if errors.Is(err, errapp.AccessDataNotFound) ||
				errors.Is(err, errapp.PasswordCheckError) {
				response.ErrorMessageJSON(c, http.StatusUnauthorized, response.AccessDataError)
				return
			}

			response.ErrorMessageJSON(c, http.StatusUnauthorized, response.InternalError)
			return
		}

		response.SuccessMessageJSON(c, http.StatusOK, data)
	}
}
