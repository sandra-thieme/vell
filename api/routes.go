package api

import (
	"github.com/rkcpi/vell/rpm"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"createRepo",
		"POST",
		"/repositories",
		rpm.CreateRepo,
	},
	Route{
		"listRepos",
		"GET",
		"/repositories",
		rpm.ListRepos,
	},
	Route{
		"addRPM",
		"POST",
		"/repositories/{name}",
		rpm.AddRPM,
	},
}
