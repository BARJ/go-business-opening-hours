package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type StatusHandler struct {
}

func (h StatusHandler) GetRoutes() []Route {
	return []Route{
		{"/status/ping", http.MethodGet, h.ping},
	}
}

func (h StatusHandler) ping(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.Write([]byte("Pong!"))
}
