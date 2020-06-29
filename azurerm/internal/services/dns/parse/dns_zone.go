package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DnsZoneId struct {
	ResourceGroup string
	Name          string
}

func DnsZoneID(input string) (*DnsZoneId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse DNS Zone ID %q: %+v", input, err)
	}

	zone := DnsZoneId{
		ResourceGroup: id.ResourceGroup,
	}

	if zone.Name, err = id.PopSegment("dnszones"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &zone, nil
}
