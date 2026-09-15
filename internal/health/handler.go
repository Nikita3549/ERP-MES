package health

import (
	"context"
	"net/http"
	"time"

	"erp-mes/pkg/db"
	"erp-mes/pkg/res"
)

const (
	StatusOK   = "ok"
	StatusDead = "dead"
)

type Handler struct {
	*db.DB
}

func NewHandler(router *http.ServeMux, DB *db.DB) {
	handler := &Handler{
		DB,
	}

	router.HandleFunc("GET /health", handler.Health())
}

func (h *Handler) Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusMessage := StatusOK
		statusCode := http.StatusOK
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		err := h.DB.Health(ctx)
		if err != nil {
			statusMessage = StatusDead
			statusCode = http.StatusServiceUnavailable
		}

		res.JSON(w, statusCode, HealthRes{
			Status: statusMessage,
		})
	}
}

type HealthRes struct {
	Status string `json:"status"`
}
