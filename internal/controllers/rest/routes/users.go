package routes

import (
	"context"
	"errors"
	"log"
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
[get] /users?name=""&surname="" - get users by filters (return user list)
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

// getUsers - get users
func getUsers(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(c)
		defer cancel()

		name := c.DefaultQuery("name", "")
		surname := c.DefaultQuery("surname", "")

		data, err := service.Get(ctx, name, surname)
		if err != nil {
			log.Printf("service.Get error: %v", err)

			response.ErrorMessageJSON(c, http.StatusInternalServerError, response.InternalError)
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
		ctx, cancel := context.WithCancel(c)
		defer cancel()

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			response.ErrorMessageJSON(c, http.StatusBadRequest, err.Error())
			return
		}

		data, err := service.GetByID(ctx, id)
		if err != nil {
			log.Printf("service.GetByID error: %v", err)

			if errors.Is(err, errapp.UserDataNotFound) {
				response.ErrorMessageJSON(c, http.StatusNotFound, errapp.UserDataNotFound.Error())
				return
			}
			response.ErrorMessageJSON(c, http.StatusInternalServerError, response.InternalError)
			return
		}

		response.SuccessMessageJSON(c, http.StatusOK, data)
		return
	}
}

//
func putUser(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(c)
		defer cancel()

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

		err = service.UpdateOne(ctx, id, data)
		if err != nil {
			log.Printf("service.UpdateOne error: %v", err)

			if errors.Is(err, errapp.UserDataNotFound) {
				response.ErrorMessageJSON(c, http.StatusNotFound, errapp.UserDataNotFound.Error())
				return
			}
			response.ErrorMessageJSON(c, http.StatusInternalServerError, response.InternalError)
			return
		}

		response.SuccessMessageJSON(c, http.StatusCreated, nil)
		return
	}
}

func getUserFriends(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(c)
		defer cancel()

		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			response.ErrorMessageJSON(c, http.StatusBadRequest, err.Error())
			return
		}

		data, err := service.GetUserFriends(ctx, id)
		if err != nil {
			log.Printf("service.GetUserFriends error: %v", err)

			if errors.Is(err, errapp.UserDataNotFound) {
				response.ErrorMessageJSON(c, http.StatusNotFound, errapp.UserDataNotFound.Error())
				return
			}

			response.ErrorMessageJSON(c, http.StatusInternalServerError, response.InternalError)
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

func putUserFriend(service user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithCancel(c)
		defer cancel()

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

		err = service.AddUserFriend(ctx, userID, friendID)
		if err != nil {
			log.Printf("service.AddUserFriend error: %v", err)

			if errors.Is(err, errapp.UserDataNotFound) {
				response.ErrorMessageJSON(c, http.StatusNotFound, errapp.UserDataNotFound.Error())
				return
			}

			response.ErrorMessageJSON(c, http.StatusInternalServerError, response.InternalError)
			return
		}

		response.SuccessMessageJSON(c, http.StatusOK, nil)
		return
	}
}
