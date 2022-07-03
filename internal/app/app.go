package app

import (
	"log"

	"social-network/internal/config"
	"social-network/internal/controllers/rest/routes"
	"social-network/internal/repositories/mysql"
	"social-network/internal/service/auth"
	"social-network/internal/service/signup"
	"social-network/internal/service/user"

	"github.com/gin-gonic/gin"
)

func Start() {
	log.Println("start application...")

	// conf := config.New()

	dbConf := config.NewDB()

	log.Println("success get config")

	storage := mysql.New(dbConf)

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

	if err := r.Run(":3004"); err != nil {
		log.Printf("%v", err)
	}
}
