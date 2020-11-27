package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SnapshotId struct {
	ResourceGroup     string
	NetAppAccountName string
	CapacityPoolName  string
	VolumeName        string
	Name              string
}

func SnapshotID(input string) (*SnapshotId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse NetApp Snapshot ID %q: %+v", input, err)
	}

	service := SnapshotId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.NetAppAccountName, err = id.PopSegment("netAppAccounts"); err != nil {
		return nil, err
	}

	if service.CapacityPoolName, err = id.PopSegment("capacityPools"); err != nil {
		return nil, err
	}

	if service.VolumeName, err = id.PopSegment("volumes"); err != nil {
		return nil, err
	}

	if service.Name, err = id.PopSegment("snapshots"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
