package api

import (
	"github.com/kuuyee/matryoshka-b-multimedia/internal/handlers"
)

type API struct {
	handlers map[string]handlers.H
}

func (a *API) RegisterServiceHandler(name string, h handlers.H) {
	a.handlers[name] = h
}

func NewAPI() *API {
	return &API{
		handlers: make(map[string]handlers.H),
	}
}
