package loadbalancer

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

// This is intentionally private to discourage usage elsewhere (since this behaviour will be removed in 3.0)
// TODO: remove in 3.0
func getZones(ctx context.Context, client *resources.ProvidersClient, resourceType, location string) (*[]string, error) {
	provider, err := client.Get(ctx, "Microsoft.Network", "")
	if err != nil {
		return nil, err
	}
	normalizedLocation := azure.NormalizeLocation(location)
	for _, resource := range *provider.ResourceTypes {
		if resource.ResourceType == nil || *resource.ResourceType != resourceType {
			continue
		}
		if resource.ZoneMappings == nil {
			continue
		}
		for _, zone := range *resource.ZoneMappings {
			if zone.Location != nil && azure.NormalizeLocation(*zone.Location) == normalizedLocation {
				return zone.Zones, nil
			}
		}
	}
	return nil, fmt.Errorf("not found zone mapping for resource %v in location %v", resourceType, location)
}
