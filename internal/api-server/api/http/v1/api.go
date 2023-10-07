package v1

import (
	"github.com/julienschmidt/httprouter"
	"github.com/taaanechka/rest-api-go/internal/api-server/api/http/v1/users"
	userservice "github.com/taaanechka/rest-api-go/internal/api-server/services"
	"github.com/taaanechka/rest-api-go/internal/handlers"
	"github.com/taaanechka/rest-api-go/pkg/logging"
)

type API struct {
	lg      *logging.Logger
	service *userservice.Service
}

func NewHandler(lg *logging.Logger, service *userservice.Service) handlers.Handler {
	return &API{
		service: service,
		lg:      lg,
	}
}

func (a *API) Register(router *httprouter.Router) {
	users := users.NewHandler(a.lg, a.service)
	users.Register(router)
}
