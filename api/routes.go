package api

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc apiHandler
}

var routes = []Route{
	Route{
		"createRepo",
		"POST",
		"/repositories",
		CreateRepo,
	},
	Route{
		"listRepos",
		"GET",
		"/repositories",
		ListRepos,
	},
	Route{
		"addRPM",
		"POST",
		"/repositories/{name}/packages",
		AddRPM,
	},
	Route{
		"listRPMs",
		"GET",
		"/repositories/{name}/packages",
		ListPackages,
	},
}
