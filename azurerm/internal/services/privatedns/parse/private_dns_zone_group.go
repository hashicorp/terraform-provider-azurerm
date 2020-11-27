package parse

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PrivateDnsZoneGroupId struct {
	ResourceGroup string
	Name          string
}

func PrivateDnsZoneGroupID(input string) (*PrivateDnsZoneGroupId, error) {
	if len(strings.TrimSpace(input)) == 0 {
		return nil, fmt.Errorf("unable to parse Private DNS Zone Group ID %q: input is empty", input)
	}

	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Private DNS Zone Group ID %q: %+v", input, err)
	}

	privateDnsZoneGroup := PrivateDnsZoneGroupId{
		ResourceGroup: id.ResourceGroup,
	}

	if privateDnsZoneGroup.Name, err = id.PopSegment("privateDnsZoneGroups"); err != nil {
		return nil, err
	}

	return &privateDnsZoneGroup, nil
}
