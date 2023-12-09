package routers

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"gotham/app"
	"gotham/config"
	"gotham/controllers"
	"gotham/docs"
	GMiddleware "gotham/middlewares"
)

func Route(e *echo.Echo) {
	docs.SwaggerInfo.Title = "SPRINTASIA API"
	docs.SwaggerInfo.Description = "..."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/doc/*", echoSwagger.WrapHandler)

	// server
	e.GET("/status/ping", controllers.ServerController{}.Ping)
	e.GET("/status/version", controllers.ServerController{}.Version)

	v1 := e.Group("/v1")

	// login
	v1.POST("/login", app.Application.Container.GetAuthController().Login)

	//register
	v1.POST("/register", app.Application.Container.GetAuthController().Register)

	r := v1.Group("/restricted")

	c := middleware.JWTConfig{
		Claims:     &config.JwtCustomClaims{},
		SigningKey: []byte(config.Conf.SecretKey),
	}

	r.Use(middleware.JWTWithConfig(c))
	r.Use(app.Application.Container.GetAuthMiddleware().AuthMiddleware)

	// user
	r.GET("/users/:user", app.Application.Container.GetUserController().Show, GMiddleware.Or(app.Application.Container.GetIsAdminMiddleware(), app.Application.Container.GetIsVerifiedMiddleware()))
	r.GET("/users", app.Application.Container.GetUserController().Index)

	//task
	r.GET("/task", app.Application.Container.GetTaskController().GetTaskList)
	r.GET("/task/:id", app.Application.Container.GetTaskController().GetTaskListByID)
	r.POST("/task", app.Application.Container.GetTaskController().CreateTaskList)
	r.PUT("/task/:id", app.Application.Container.GetTaskController().Update)
	r.PUT("/task/status/:id", app.Application.Container.GetTaskController().UpdateTaskList)
	r.DELETE("/task/:id", app.Application.Container.GetTaskController().DeleteTaskList)

	// Start server
	go func() {
		if err := e.Start(":" + config.Conf.Port); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
