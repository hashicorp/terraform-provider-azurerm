package rule

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tools/docthebuilder/data"
)

var (
	// Exceptions contains resources and rules they should skip.
	// Name format: "(d|e|r).<resource name>"
	// TODO: externalize as configurable item?
	Exceptions = map[string]map[string]struct{}{
		"r.azurerm_resource_provider_registration": {
			"S002": struct{}{},
		},
	}
)

func ShouldSkipRule(resourceType data.ResourceType, resourceName string, ruleName string) bool {
	prefix := ""
	switch resourceType {
	case data.ResourceTypeData:
		prefix = "d"
	case data.ResourceTypeEphemeral:
		prefix = "e"
	case data.ResourceTypeResource:
		prefix = "r"
	}

	if v, ok := Exceptions[fmt.Sprintf("%s.%s", prefix, resourceName)]; ok {
		if _, ok := v[ruleName]; ok {
			return true
		}
	}

	return false
}
