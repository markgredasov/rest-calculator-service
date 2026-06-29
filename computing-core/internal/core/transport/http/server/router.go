package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/markgredasov/rest-calculator-service/internal/core/transport/http/middleware"
)

type ApiPrefix string

var (
	ApiPrefixInternal = ApiPrefix("internal")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiPrefix ApiPrefix
}

func NewApiVersionRouter(prefix ApiPrefix) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:  http.NewServeMux(),
		apiPrefix: prefix,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		r.Handle(pattern, core_http_middleware.ChainMiddleware(route.Handler, route.Middleware...))
	}
}
