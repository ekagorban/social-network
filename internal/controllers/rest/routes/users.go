package routes

import (
	"errors"
	"net/http"

	"social-network/internal/controllers/rest/middleware"
	"social-network/internal/controllers/rest/response"
	"social-network/internal/errapp"
	"social-network/internal/service/auth"
	"social-network/internal/service/user"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

/*
[get] /users - get all users (return user list)
[get] /user/:id - get user by ID (return user)
[put] /user/:id -  change user data
[get] /friends/:id - get all user friends
[put] /friend/:userID/:friendID - add other user to friend list
*/

func User(r *gin.Engine, service user.Service, authService auth.Service) {
	v1 := r.Group("/v1").Use(middleware.CheckAccess(authService))

	v1.GET("/users", getUsers(service))

	v1.GET("/user/:id", getUser(service))
	v1.PUT("/user/:id", putUser(service))

	v1.GET("/friends/:id", getUserFriends(service))

	v1.PUT("/friend/:userID/:friendID", putUserFriend(service))
}

// getUsers - get all users
func getUsers(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := service.GetAll()
		if err != nil {
			response.ErrorMessageJSON(c, http.StatusInternalServerError, err.Error())
			return
		}

		if len(data) == 0 {
			response.SuccessMessageJSON(c, http.StatusNoContent, nil)
			return
		}

		response.SuccessMessageJSON(c, http.StatusOK, data)
		return
	}
}

// getUser - get one user by id
func getUser(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			response.ErrorMessageJSON(c, http.StatusBadRequest, err.Error())
			return
		}

		data, err := service.GetOne(id)
		if err != nil {
			if errors.Is(err, errapp.UserDataNotFound) {
				response.SuccessMessageJSON(c, http.StatusNotFound, nil)
				return
			}
			response.ErrorMessageJSON(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.SuccessMessageJSON(c, http.StatusOK, data)
		return
	}
}

//
func putUser(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			response.ErrorMessageJSON(c, http.StatusBadRequest, err.Error())
			return
		}

		var data user.Data
		if err := c.BindJSON(&data); err != nil {
			response.ErrorMessageJSON(c, http.StatusBadRequest, err.Error())
			return
		}

		err = service.UpdateOne(id, data)
		if err != nil {
			if errors.Is(err, errapp.UserDataNotFound) {
				response.SuccessMessageJSON(c, http.StatusNotFound, nil)
				return
			}
			response.ErrorMessageJSON(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.SuccessMessageJSON(c, http.StatusCreated, nil)
		return
	}
}

func getUserFriends(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			response.ErrorMessageJSON(c, http.StatusBadRequest, err.Error())
			return
		}

		data, err := service.GetUserFriends(id)
		if err != nil {
			if errors.Is(err, errapp.UserDataNotFound) {
				response.SuccessMessageJSON(c, http.StatusNotFound, nil)
				return
			}
			response.ErrorMessageJSON(c, http.StatusInternalServerError, err.Error())
			return
		}

		if len(data) == 0 {
			response.SuccessMessageJSON(c, http.StatusNoContent, data)
			return
		}

		response.SuccessMessageJSON(c, http.StatusOK, data)
		return
	}
}

func putUserFriend(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.Param("userID"))
		if err != nil {
			response.ErrorMessageJSON(c, http.StatusBadRequest, err.Error())
			return
		}

		friendID, err := uuid.Parse(c.Param("friendID"))
		if err != nil {
			response.ErrorMessageJSON(c, http.StatusBadRequest, err.Error())
			return
		}

		err = service.AddUserFriend(userID, friendID)
		if err != nil {
			if errors.Is(err, errapp.UserDataNotFound) {
				response.SuccessMessageJSON(c, http.StatusNotFound, nil)
				return
			}
			response.ErrorMessageJSON(c, http.StatusInternalServerError, err.Error())
			return
		}

		response.SuccessMessageJSON(c, http.StatusOK, nil)
		return
	}
}
