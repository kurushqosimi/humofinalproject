package main

import (
	"log"
	"main/internal/configs"
	"main/internal/handlers"
	"main/internal/repositories"
	"main/internal/services"
	"main/pkg/logging"
	"net/http"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal("faced an error while running the project")
	}
}

func run() error {
	logger, err := logging.InitializeLogger()
	if err != nil {
		return err
	}
	config, err := configs.InitConfigs()
	if err != nil {
		return err
	}
	address := config.ServerSetting.Host + config.ServerSetting.Port
	conn, err := repositories.GetConnection(config)
	if err != nil {
		return err
	}
	service := services.NewService(conn, logger)
	handler := handlers.NewHandler(service, logger)
	router := handlers.NewRouter(handler)
	srv := http.Server{
		Addr:    address,
		Handler: router,
	}
	err = srv.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}
