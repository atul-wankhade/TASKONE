// @title TaskOne REST API
// @version 1.0
// @description This is the TaskOne API microservice built in Go.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.taskone.io/support
// @contact.email support@taskone.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"TASKONE/config"
	"TASKONE/controller"
	"TASKONE/db"
	"TASKONE/middleware"
	"TASKONE/repository"
	"TASKONE/utils"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	_ "TASKONE/docs" // Swagger generated docs

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {

	//Load Config
	config.LoadConfig()

	utils.InitLogger()
	defer utils.Logger.Sync()

	sqlDB := db.ConnectMySQL()
	db.InitMongo()

	userRepo := repository.NewUserRepository(sqlDB)
	logRepo := repository.NewLogRepository(db.MongoDatabase)

	userController := controller.NewUserController(userRepo, logRepo)
	authController := controller.NewAuthController(userRepo)

	r := chi.NewRouter()
	//use middleware for logging and recover
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)

	//swagger documentation route
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	//create login handler
	r.Post("/login", authController.LoginHandler)
	r.Post("/register", authController.RegisterHandler)
	//create group of protected apis
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/protected", controller.ProtectedHandler)
		r.Get("/user/{id}", userController.GetUserHandler)
		r.Get("/user", userController.GetUserHandlerNew)
	})

	server := &http.Server{
		Addr:    ":" + config.AppConfig.Port,
		Handler: r,
	}

	go func() {
		//start server
		utils.Logger.Info("Starting the Server", zap.String("port", config.AppConfig.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.Logger.Fatal("Failed to start the Server", zap.Error(err))
		}
	}()

	quite := make(chan os.Signal, 1)
	signal.Notify(quite, syscall.SIGTERM, syscall.SIGINT)
	<-quite

	utils.Logger.Info("Shutting down server....")

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		utils.Logger.Fatal("Server Forced to shutdown", zap.Error(err))
	}

	utils.Logger.Info("Server exited properly...")

}
