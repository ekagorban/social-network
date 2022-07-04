package app

import (
	"fmt"
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
	log.Println("start init application...")

	appConf := config.AppNew()
	dbConf := config.DBNew()

	log.Printf("success get app config: %+v; db config: %+v; ", appConf, dbConf)

	storage, err := mysql.New(dbConf)
	if err != nil {
		log.Printf("mysql.New error: %v", err)
		return
	}

	userService := user.NewService(storage)
	signUpService := signup.NewService(storage)
	authService := auth.NewService(appConf, storage)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	routes.Ping(r)

	routes.SignUp(r, signUpService)
	routes.SignIn(r, authService)
	routes.User(r, userService, authService)

	log.Println("success finish init application")

	if err := r.Run(fmt.Sprintf(":%s", appConf.ListenPort)); err != nil {
		log.Printf("%v", err)
	}
}
