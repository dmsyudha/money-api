package service

import (
	"fmt"
	"net/http"

	db "github.com/dmsyudha/money-api/lib/pq"
)

type HealthService interface {
	HealthCheckHandler(w http.ResponseWriter, r *http.Request)
	PingDBHandler(w http.ResponseWriter, r *http.Request)
}

type healthServiceImpl struct{}

func NewHealthService() HealthService {
	return &healthServiceImpl{}
}

// HealthCheckHandler checks if the service is up and running
func (h *healthServiceImpl) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Healthy!")
}

// PingDBHandler checks if the database connection is alive
func (h *healthServiceImpl) PingDBHandler(w http.ResponseWriter, r *http.Request) {
	db, err := db.ConnectToDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error connecting to database: %v", err)
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error getting database handle: %v", err)
		return
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Database connection error: %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Database connection successful!")
}
