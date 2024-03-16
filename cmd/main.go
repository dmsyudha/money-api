package main

import (
	"fmt"
	"net/http"

	"github.com/dmsyudha/money-api/internal/service"
)

func main() {
	services := setupDependencies()
	setupRoutes(services)
	healthService := service.NewHealthService()

	http.HandleFunc("/health", healthService.HealthCheckHandler)
	http.HandleFunc("/ping-db", healthService.PingDBHandler)

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
