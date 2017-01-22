package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

type Routes []Route

func registerRoute(router *mux.Router, route Route, handler http.Handler) {
	router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	api := router.PathPrefix("/api/").Subrouter()
	for _, route := range routes {
		registerRoute(api, route, route.Handler)
	}
	router.PathPrefix("/app/").Handler(http.HandlerFunc(notFound))
	router.PathPrefix("/").Handler(http.FileServer(assetFS()))
	return router
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

var routes = Routes{
	Route{
		"SourceParentDir",
		"GET",
		"/selectSource/parentDir",
		ParentDir{&sourceCurrentDir},
	},
	Route{
		"GetSourceCurrentDir",
		"GET",
		"/selectSource/currentDir",
		GetCurrentDir{&sourceCurrentDir},
	},
	Route{
		"SetSourceCurrentDir",
		"POST",
		"/selectSource/currentDir",
		SetCurrentDir{&sourceCurrentDir},
	},
	Route{
		"DestParentDir",
		"GET",
		"/selectDest/parentDir",
		ParentDir{&destCurrentDir},
	},
	Route{
		"GetDestCurrentDir",
		"GET",
		"/selectDest/currentDir",
		GetCurrentDir{&destCurrentDir},
	},
	Route{
		"SetDestCurrentDir",
		"POST",
		"/selectDest/currentDir",
		SetCurrentDir{&destCurrentDir},
	},
	Route{
		"SetDestFile",
		"PUT",
		"/convert",
		Convert,
	},
}
