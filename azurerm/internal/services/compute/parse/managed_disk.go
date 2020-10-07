package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ManagedDiskId struct {
	ResourceGroup string
	Name          string
}

func NewManagedDiskId(resourceGroup, name string) ManagedDiskId {
	return ManagedDiskId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id ManagedDiskId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/disks/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func ManagedDiskID(input string) (*ManagedDiskId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Managed Disk ID %q: %+v", input, err)
	}

	disk := ManagedDiskId{
		ResourceGroup: id.ResourceGroup,
	}

	if disk.Name, err = id.PopSegment("disks"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &disk, nil
}
