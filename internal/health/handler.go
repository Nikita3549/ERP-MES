package health

import (
	"context"
	"log"
	"net/http"
	"time"

	"erp-mes/pkg/db"
	"erp-mes/pkg/res"
)

const (
	StatusOK   = "ok"
	StatusDead = "dead"
)

type DB interface {
	Health(context.Context) error
}

type Handler struct {
	db DB
}

func NewHandler(router *http.ServeMux, DB *db.DB) {
	handler := &Handler{
		db: DB,
	}

	router.HandleFunc("GET /health", handler.Health())
}

func (h *Handler) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dbStatus := StatusOK
		statusCode := http.StatusOK
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		err := h.db.Health(ctx)
		if err != nil {
			dbStatus = StatusDead
			statusCode = http.StatusServiceUnavailable
			log.Printf("database ping: %v", err)
		}

		res.JSON(w, statusCode, HealthRes{
			Status:   StatusOK,
			DBStatus: dbStatus,
		})
	}
}

type HealthRes struct {
	Status   string `json:"status"`
	DBStatus string `json:"databaseStatus"`
}
