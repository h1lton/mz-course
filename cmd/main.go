package main

import (
	todo "github.com/h1lton/mz-course"
	"github.com/h1lton/mz-course/pkg/handler"
	"github.com/h1lton/mz-course/pkg/repository"
	"github.com/h1lton/mz-course/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	logrus.SetOutput(colorable.NewColorableStdout())

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	println(viper.GetString("port"))
	if err := srv.Run(viper.GetString("addr"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
