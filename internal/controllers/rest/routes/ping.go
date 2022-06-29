package routes

import (
	"net/http"

	"social-network/internal/controllers/rest/response"

	"github.com/gin-gonic/gin"
)

func Ping(r *gin.Engine) {
	r.GET("/ping", ping)
}

func ping(c *gin.Context) {
	response.SuccessMessageJSON(c, http.StatusOK, response.AliveMsg)
	return
}
