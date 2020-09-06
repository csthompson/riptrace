package routes

import (
	"net/http"
	"os"

	"github.com/csthompson/riptrace/server/internal/handlers/httphandlers"

	muxhandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Router struct {
	R *mux.Router
}

func (self *Router) GetHandler() http.Handler {
	r := self.R

	s := r.PathPrefix("/agent/v1").Subrouter()
	agents := httphandlers.AgentsHandler{}
	s.HandleFunc("/list", agents.List).Methods("GET")

	//TODO: Get an agent's profile
	//@tags: handler, enhancement

	//Attach logging middleware for each request
	loggedRouter := muxhandlers.LoggingHandler(os.Stdout, self.R)
	return loggedRouter
}
