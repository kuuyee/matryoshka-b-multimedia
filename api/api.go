package api

import (
	"github.com/kuuyee/matryoshka-b-multimedia/internal/handlers"
)

// API rest-api接口定义
type API struct {
	handlers map[string]handlers.H
}

// RegisterServiceHandler 注册新服务
func (a *API) RegisterServiceHandler(name string, h handlers.H) {
	a.handlers[name] = h
}

// NewAPI creates a new API instance
func NewAPI() *API {
	return &API{
		handlers: make(map[string]handlers.H),
	}
}
