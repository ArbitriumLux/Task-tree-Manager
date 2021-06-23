package main

import (
	"TasksManager/models"
	"TasksManager/router"
	"TasksManager/server"
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

var (
	logLevel = flag.String("LogLevel", "debug", "Log Level")
)

// Basic execution logic: main > router.Router > handlers > models
func main() {
	flag.Parse()

	//configureLogger
	configureLogger := func() error {
		log.SetFormatter(&log.JSONFormatter{})
		level, err := log.ParseLevel(*logLevel)
		if err != nil {
			return err
		}
		log.SetLevel(level)
		file, err := os.OpenFile("TasksManager.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Info("Failed to log to file")
		}
		return nil
	}
	if err := configureLogger(); err != nil {
		return
	}
	//Database migrate
	server.DataBaseConnection()
	server.Db.AutoMigrate(&models.Tasks{}, &models.MiniTasks{}, &models.LaborCosts{})
	log.Info("Migration Successful!")
	defer server.Db.Close()
	//Start listening to http routes
	router.Router()
}
