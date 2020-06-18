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
		return nil, fmt.Errorf("unable to parse Private DNS Zone Group ID %q: input is empty", input)
	}

	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Private DNS Zone Group ID %q: %+v", input, err)
	}

	privateDnsZoneGroup := NameResourceGroup{
		ResourceGroup: id.ResourceGroup,
	}

	if privateDnsZoneGroup.Name, err = id.PopSegment("privateDnsZoneGroups"); err != nil {
		return nil, err
	}

	if privateDnsZoneGroup.ID = input; err != nil {
		return nil, err
	}

	return &privateDnsZoneGroup, nil
}

func PrivateDnsZoneResourceIDs(input []interface{}) (*[]NameResourceGroup, error) {
	results := make([]NameResourceGroup, 0)

	for _, item := range input {
		v := item.(string)

		if privateDnsZone, err := PrivateDnsZoneResourceID(v); err != nil {
			return nil, fmt.Errorf("unable to parse Private DNS Zone ID %q: %+v", v, err)
		} else {
			results = append(results, *privateDnsZone)
		}
	}

	return &results, nil
}

func PrivateDnsZoneResourceID(input string) (*NameResourceGroup, error) {
	if len(strings.TrimSpace(input)) == 0 {
		return nil, fmt.Errorf("unable to parse Private DNS Zone ID %q: input is empty", input)
	}

	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Private DNS Zone ID %q: %+v", input, err)
	}

	privateDnsZone := NameResourceGroup{
		ResourceGroup: id.ResourceGroup,
	}

	if privateDnsZone.Name, err = id.PopSegment("privateDnsZones"); err != nil {
		return nil, err
	}

	if privateDnsZone.ID = input; err != nil {
		return nil, err
	}

	return &privateDnsZone, nil
}

func PrivateEndpointResourceID(input string) (*NameResourceGroup, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Private Endpoint ID %q: %+v", input, err)
	}

	privateEndpoint := NameResourceGroup{
		ResourceGroup: id.ResourceGroup,
	}

	if privateEndpoint.Name, err = id.PopSegment("privateEndpoints"); err != nil {
		return nil, err
	}

	if privateEndpoint.ID = input; err != nil {
		return nil, err
	}

	return &privateEndpoint, nil
}

func ValidatePrivateDnsZoneResourceID(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return
	}

	if id, err := azure.ParseAzureResourceID(v); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a resource id: %v", k, err))
	} else if _, err = id.PopSegment("privateDnsZones"); err != nil {
		errors = append(errors, fmt.Errorf("Can not parse %q as a private dns zone resource id: %v", k, err))
	}

	return warnings, errors
}
