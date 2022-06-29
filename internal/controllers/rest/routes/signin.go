package routes

import (
	"net/http"

	"social-network/internal/controllers/rest/response"
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
			response.ErrorMessageJSON(c, http.StatusBadRequest, err.Error())
			return
		}

		data, err := service.SignIn(dataSignIn)
		if err != nil {
			response.ErrorMessageJSON(c, http.StatusUnauthorized, err.Error())
			return
		}

		response.SuccessMessageJSON(c, http.StatusOK, data)
	}
}
