package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NetAppSnapshotId struct {
	ResourceGroup string
	AccountName   string
	PoolName      string
	VolumeName    string
	Name          string
}

func NetAppSnapshotID(input string) (*NetAppSnapshotId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse NetApp Snapshot ID %q: %+v", input, err)
	}

	service := NetAppSnapshotId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.AccountName, err = id.PopSegment("netAppAccounts"); err != nil {
		return nil, err
	}

	if service.PoolName, err = id.PopSegment("capacityPools"); err != nil {
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
