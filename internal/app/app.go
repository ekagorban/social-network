package app

import (
	"log"

	"social-network/internal/controllers/rest/routes"
	"social-network/internal/repositories/mysql"
	"social-network/internal/service/auth"
	"social-network/internal/service/signup"
	"social-network/internal/service/user"

	"github.com/gin-gonic/gin"
)

func Start() {

	//storage := memory.New()
	storage := mysql.New()

	userService := user.NewService(storage)
	signUpService := signup.NewService(storage)
	authService := auth.NewService(storage)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	routes.Ping(r)

	routes.SignUp(r, signUpService)
	routes.SignIn(r, authService)
	routes.User(r, userService, authService)

	if err := r.Run("127.0.0.1:8081"); err != nil {
		log.Printf("%v", err)
	}
}
