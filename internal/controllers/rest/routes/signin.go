package routes

import (
	"errors"
	"log"
	"net/http"

	"social-network/internal/controllers/rest/response"
	"social-network/internal/errapp"
	"social-network/internal/service/auth"

	"github.com/gin-gonic/gin"
)

/*
[post] /signin - sign in with login and password (return )
*/

func SignIn(r *gin.Engine, service auth.Service) {
	v1 := r.Group("/v1")
	v1.POST("/signin", signIn(service))
}

func signIn(service auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataSignIn auth.SignInData
		if err := c.ShouldBindJSON(&dataSignIn); err != nil {
			log.Println(err)
			response.ErrorMessageJSON(c, http.StatusBadRequest, err.Error())
			return
		}

		data, err := service.SignIn(dataSignIn)
		if err != nil {
			log.Println(err)
			if errors.Is(err, errapp.AccessDataNotFound) {
				response.ErrorMessageJSON(c, http.StatusConflict, errapp.AccessDataNotFound.Error())
				return
			}

			response.ErrorMessageJSON(c, http.StatusUnauthorized, response.InternalError)
			return
		}

		response.SuccessMessageJSON(c, http.StatusOK, data)
	}
}
