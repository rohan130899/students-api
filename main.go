package main

import (
	"students-api/project/initializers"
	"students-api/project/student"

	transportHTTP "students-api/project/transport"

	log "github.com/sirupsen/logrus"
)

func init() {
	initializers.InitLogger()
	initializers.LoadEnvVariables()
}

func Run() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("Setting up Our App")
	db, err := initializers.ConnectToDB()

	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}
	studentService := student.NewService(db)
	handler := transportHTTP.NewHandler(studentService)
	if err := handler.Serve(); err != nil {
		log.Error("failed to gracefully serve our application")
		return err
	}
	return nil
}
func main() {
	if err := Run(); err != nil {
		log.Error(err)
		log.Fatal("Error starting up our REST API")
	}
}
