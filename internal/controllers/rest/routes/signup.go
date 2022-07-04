package routes

import (
	"context"
	"errors"
	"log"
	"net/http"

	"social-network/internal/controllers/rest/response"
	"social-network/internal/errapp"
	"social-network/internal/service/signup"

	"github.com/gin-gonic/gin"
)

/*
[post] /signup - create new user
*/

func SignUp(r *gin.Engine, service signup.Service) {
	v1 := r.Group("/v1")
	v1.POST("/signup", signUp(service))
}

func signUp(service signup.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(c)
		defer cancel()

		var data signup.Data
		if err := c.ShouldBindJSON(&data); err != nil {
			response.ErrorMessageJSON(c, http.StatusBadRequest, err.Error())
			return
		}

		err := service.CreateUser(ctx, data)
		if err != nil {
			log.Printf("service.CreateUser error: %v", err)

			if errors.Is(err, errapp.LoginExist) {
				response.ErrorMessageJSON(c, http.StatusConflict, errapp.LoginExist.Error())
				return
			}

			response.ErrorMessageJSON(c, http.StatusInternalServerError, response.InternalError)
			return
		}

		response.SuccessMessageJSON(c, http.StatusCreated, nil)
		return
	}
}
