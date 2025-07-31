package data

import "fmt"

type ResourceType string

const (
	ResourceTypeData      ResourceType = "Data Source"
	ResourceTypeEphemeral ResourceType = "Ephemeral Resource"
	ResourceTypeResource  ResourceType = "Resource"
)

var (
	resourceFilePathPattern    = "%s/%s_%s.go"
	resourceFileGenPathPattern = "%s/%s_%s_gen.go"

	ResourceTypeToFileSuffix = map[ResourceType]string{
		ResourceTypeData:      "data_source",
		ResourceTypeEphemeral: "ephemeral",
		ResourceTypeResource:  "resource",
	}

	ResourceTypeToDocumentationSubPath = map[ResourceType]string{
		ResourceTypeData:      "d",
		ResourceTypeEphemeral: "ephemeral-resources",
		ResourceTypeResource:  "r",
	}
)

func (r ResourceType) String() string {
	return string(r)
}

func expectedResourceCodePath(pattern string, name string, service Service, resourceType ResourceType) string {
	return fmt.Sprintf(pattern, service.Path, name, ResourceTypeToFileSuffix[resourceType])
}
