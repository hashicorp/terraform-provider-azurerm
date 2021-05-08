package network

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2020-06-01/resources"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

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
	return nil, nil
}
