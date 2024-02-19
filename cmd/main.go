package main

import (
	todo "github.com/h1lton/mz-course"
	"github.com/h1lton/mz-course/pkg/handler"
	"github.com/h1lton/mz-course/pkg/repository"
	"github.com/h1lton/mz-course/pkg/service"
	"log"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}
}
