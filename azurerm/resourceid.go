package azurerm

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

//PLEASE NOTE: This code has  been moved to terraform-provider-azurerm/azurerm/helpers/azure
//this file is simply a shim to prevent mass code changes

// ResourceID represents a parsed long-form Azure Resource Manager ID
// with the Subscription ID, Resource Group and the Provider as top-
// level fields, and other key-value pairs available via a map in the
// Path field.
type ResourceID struct {
	SubscriptionID string
	ResourceGroup  string
	Provider       string
	Path           map[string]string
}

// parseAzureResourceID converts a long-form Azure Resource Manager ID
// into a ResourceID. We make assumptions about the structure of URLs,
// which is obviously not good, but the best thing available given the
// SDK.
func parseAzureResourceID(id string) (*ResourceID, error) {
	parsed, err := azure.ParseAzureResourceID(id)
	if err != nil {
		return nil, err
	}

	return &ResourceID{
		SubscriptionID: parsed.SubscriptionID,
		ResourceGroup:  parsed.ResourceGroup,
		Provider:       parsed.Provider,
		Path:           parsed.Path,
	}, nil
}

func parseNetworkSecurityGroupName(networkSecurityGroupId string) (string, error) {
	id, err := parseAzureResourceID(networkSecurityGroupId)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Unable to Parse Network Security Group ID '%s': %+v", networkSecurityGroupId, err)
	}

	return id.Path["networkSecurityGroups"], nil
}

func parseRouteTableName(routeTableId string) (string, error) {
	id, err := parseAzureResourceID(routeTableId)
	if err != nil {
		return "", fmt.Errorf("[ERROR] Unable to parse Route Table ID '%s': %+v", routeTableId, err)
	}

	return id.Path["routeTables"], nil
}
