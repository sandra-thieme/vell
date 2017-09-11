package api

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc apiHandler
}

var routes = []Route{
	{
		"createRepo",
		"POST",
		"/repositories",
		CreateRepo,
	},
	{
		"listRepos",
		"GET",
		"/repositories",
		ListRepos,
	},
	{
		"addRPM",
		"POST",
		"/repositories/{name}/packages",
		AddRPM,
	},
	{
		"listRPMs",
		"GET",
		"/repositories/{name}/packages",
		ListPackages,
	},
	{
		"getPackageWithVersion",
		"GET",
		"/repositories/{name}/packages/{packagename}/version/{version}",
		GetPackageWithNameAndVersion,
	},
}
