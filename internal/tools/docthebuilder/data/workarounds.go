package data

/*
workarounds.go should only contain code that provides workarounds for data that doesn't follow the expected pattern.
e.g. service folder that we can't get based on `(Registration).Name()` or `(Registration).AssociatedGitHubLabel()`
*/

var (
	// ServiceFolderWorkaround provides a mapping from `(Registration).Name()` to the service directory
	ServiceFolderWorkaround = map[string]string{
		"Cognitive Services":              "cognitive",
		"CosmosDB":                        "cosmos",
		"Trusted Signing":                 "codesigning",
		"PostgreSQL":                      "postgres",
		"Resources":                       "resource",
		"Service Fabric Managed Clusters": "servicefabricmanaged",
		"Container Services":              "containers",
	}
)
