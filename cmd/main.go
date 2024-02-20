package main

import (
	todo "github.com/h1lton/mz-course"
	"github.com/h1lton/mz-course/pkg/handler"
	"github.com/h1lton/mz-course/pkg/repository"
	"github.com/h1lton/mz-course/pkg/service"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	println(viper.GetString("port"))
	if err := srv.Run(viper.GetString("addr"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occurred while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
