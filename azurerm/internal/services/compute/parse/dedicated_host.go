package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DedicatedHostId struct {
	ResourceGroup string
	HostGroup     string
	Name          string
}

func DedicatedHostID(input string) (*DedicatedHostId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Dedicated Host ID %q: %+v", input, err)
	}

	host := DedicatedHostId{
		ResourceGroup: id.ResourceGroup,
	}

	if host.HostGroup, err = id.PopSegment("hostGroups"); err != nil {
		return nil, err
	}

	if host.Name, err = id.PopSegment("hosts"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &host, nil
}
