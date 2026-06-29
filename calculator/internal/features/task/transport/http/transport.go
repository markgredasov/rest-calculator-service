package tasks_transport

import (
	"context"
	"net/http"

	"github.com/markgredasov/rest-calculator-service/internal/core/domain"
	core_http_server "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
}

type TasksService interface {
	Calculate(ctx context.Context, task domain.Task) (domain.Task, error)
}

func NewTasksHTTPHandler(service TasksService) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: service,
	}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/calculate",
			Handler: h.Calculate,
		},
	}
}
