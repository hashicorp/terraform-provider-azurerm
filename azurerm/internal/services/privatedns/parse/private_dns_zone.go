package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PrivateDnsZoneId struct {
	ResourceGroup string
	Name          string
}

func (id PrivateDnsZoneId) ID(_ string) string {
	// stub implementation until the generator is in
	return ""
}

func PrivateDnsZoneID(input string) (*PrivateDnsZoneId, error) {
	if len(strings.TrimSpace(input)) == 0 {
		return nil, fmt.Errorf("unable to parse Private DNS Zone ID %q: input is empty", input)
	}

	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Private DNS Zone ID %q: %+v", input, err)
	}

	privateDnsZone := PrivateDnsZoneId{
		ResourceGroup: id.ResourceGroup,
	}

	if privateDnsZone.Name, err = id.PopSegment("privateDnsZones"); err != nil {
		return nil, err
	}

	return &privateDnsZone, nil
}
