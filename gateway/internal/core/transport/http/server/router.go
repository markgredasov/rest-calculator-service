package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/middleware"
)

type ApiVersion string

var (
	ApiVersion1 = ApiVersion("v1")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
}

func NewApiVersionRouter(apiVersion ApiVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, core_http_middleware.ChainMiddleware(route.Handler, route.Middleware...))
	}
}
