package tfratelimiter

import (
	"strings"
)

const (
	Service_unknown         = ""
	Service_azurerm_storage = "azurerm_storage"
)

// Specification of services and rates (operations per second and burst limit).
// Could be moved to a data file at some point.
func getLimit(service Service, operation Operation) (float64, int) {
	switch service {
	case Service_azurerm_storage:
		switch operation {
		case Read:
			return 790.0 / (60 * 5), 3
		case List:
			return 95.0 / (60 * 5), 2
		}
	}

	return -1, -1
}

func serviceForResourceType(resourceType string) Service {
	if strings.HasPrefix(resourceType, "azurerm_storage_account") {
		return Service_azurerm_storage
	}

	return Service_unknown
}

// Indicates the number of read API calls that refreshing a single resource may
// make.  You can determine this by reading the Refresh code in the terraform
// provider for that resource type.
func readsForResourceType(resourceType string) int {
	switch resourceType {
	case "azurerm_storage_account":
		return 2
	}

	return 1
}
