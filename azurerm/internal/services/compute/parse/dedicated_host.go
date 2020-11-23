package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DedicatedHostId struct {
	SubscriptionId string
	ResourceGroup  string
	HostGroup      string
	Name           string
}

func NewDedicatedHostId(id DedicatedHostGroupId, name string) DedicatedHostId {
	return DedicatedHostId{
		SubscriptionId: id.SubscriptionId,
		ResourceGroup:  id.ResourceGroup,
		HostGroup:      id.Name,
		Name:           name,
	}
}

func (id DedicatedHostId) ID(_ string) string {
	base := NewDedicatedHostGroupId(id.SubscriptionId, id.ResourceGroup, id.HostGroup).ID("")
	return fmt.Sprintf("%s/hosts/%s", base, id.Name)
}

func DedicatedHostID(input string) (*DedicatedHostId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Dedicated Host ID %q: %+v", input, err)
	}

	host := DedicatedHostId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
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
