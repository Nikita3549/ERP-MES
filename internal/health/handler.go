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
	StatusFail = "fail"
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
		res.JSON(w, http.StatusOK, HealthRes{
			Status: StatusOK,
		})
	}
}

func (h *Handler) Ready() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := StatusOK
		code := http.StatusOK
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := h.db.Health(ctx)
		if err != nil {
			status = StatusFail
			code = http.StatusServiceUnavailable
			log.Printf("database ping error: %v", err)
		}

		res.JSON(w, code, ReadyRes{Status: status})
	}
}

type HealthRes struct {
	Status string `json:"status"`
}

type ReadyRes struct {
	Status string `json:"status"`
}
