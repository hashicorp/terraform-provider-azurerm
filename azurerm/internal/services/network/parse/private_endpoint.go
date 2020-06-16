package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NameResourceGroup struct {
	ResourceGroup string
	Name          string
	ID            string
}

func PrivateDnsZoneGroupResourceID(input string) (*NameResourceGroup, error) {
	if len(strings.TrimSpace(input)) == 0 {
		return NameResourceGroup{}, fmt.Errorf("input is empty for parse.PrivateDnsZoneGroupResourceID")
	}

	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return NameResourceGroup{}, fmt.Errorf("unable to parse Private DNS Zone Group ID %q: %+v", input, err)
	}

	privateDnsZoneGroup := NameResourceGroup{
		Name:          id.Path["privateDnsZoneGroups"],
		ResourceGroup: id.ResourceGroup,
		ID:            input,
	}

	return privateDnsZoneGroup, nil
}

func PrivateDnsZoneResourceIDs(input []interface{}) ([]NameResourceGroup, error) {
	results := make([]NameResourceGroup, 0)

	for _, item := range input {
		v := item.(string)

		id, err := azure.ParseAzureResourceID(v)
		if err != nil {
			return nil, fmt.Errorf("unable to parse Private DNS Zone ID %q: %+v", input, err)
		}

		privateDnsZone := NameResourceGroup{
			Name:          id.Path["privateDnsZones"],
			ResourceGroup: id.ResourceGroup,
			ID:            v,
		}

		results = append(results, privateDnsZone)
	}

	return results, nil
}

func PrivateEndpointResourceID(input string) (NameResourceGroup, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return NameResourceGroup{}, fmt.Errorf("unable to parse Private Endpoint ID %q: %+v", input, err)
	}

	privateEndpoint := NameResourceGroup{
		Name:          id.Path["privateEndpoints"],
		ResourceGroup: id.ResourceGroup,
		ID:            input,
	}

	return privateEndpoint, nil
}
