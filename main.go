package main

import (
	"TASKONE/config"
	"TASKONE/controller"
	"TASKONE/middleware"
	"TASKONE/utils"
	"net/http"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/go-chi/chi"
)

func main() {

	//Load Config
	config.LoadConfig()
	//add logger and sync it
	utils.InitLogger()
	defer utils.Logger.Sync()
	//create new roter
	r := chi.NewRouter()
	//use middleware for logging and recover
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	//create login handler
	r.Post("/Login", controller.LoginHandler)
	//create group of protected apis
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/protected", controller.ProtectedHandler)
	})
	//start server
	utils.Logger.Info("Starting the Server", zap.String("port", config.AppConfig.Port))
	err := http.ListenAndServe(":"+config.AppConfig.Port, r)
	if err != nil {
		utils.Logger.Fatal("Failed to start the Server", zap.Error(err))
	}

}
