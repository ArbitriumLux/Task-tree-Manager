package router

import (
	"TasksManager/middleware"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const port = ":9090"

//Router
func Router() {
	router := mux.NewRouter()
	router.HandleFunc("/", middleware.IndexPage).Methods("GET")
	router.HandleFunc("/tree/{TaskName}", middleware.GetTree).Methods("GET")
	router.HandleFunc("/tree/delete/Task", middleware.DeleteTask).Methods("GET")
	router.HandleFunc("/tree/delete/MiniTask", middleware.DeleteMiniTask).Methods("GET")
	router.HandleFunc("/tree/delete/LaborCost", middleware.DeleteLaborCost).Methods("GET")
	router.HandleFunc("/create/task", middleware.CreateTask).Methods("POST")
	router.HandleFunc("/create/miniTask", middleware.CreateMiniTask).Methods("POST")
	router.HandleFunc("/create/laborCost", middleware.CreateLaborCost).Methods("POST")
	router.HandleFunc("/tree/update/miniTask", middleware.UpdateMiniTask).Methods("POST")
	router.HandleFunc("/tree/update/laborCost", middleware.UpdateLaborCost).Methods("POST")

	handler := cors.Default().Handler(router)
	log.Info(http.ListenAndServe(port, handler))
}
