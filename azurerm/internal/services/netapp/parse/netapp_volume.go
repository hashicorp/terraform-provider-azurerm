package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type NetAppVolumeId struct {
	ResourceGroup string
	AccountName   string
	PoolName      string
	Name          string
}

func NetAppVolumeID(input string) (*NetAppVolumeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse NetApp Volume ID %q: %+v", input, err)
	}

	service := NetAppVolumeId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.AccountName, err = id.PopSegment("netAppAccounts"); err != nil {
		return nil, err
	}

	if service.PoolName, err = id.PopSegment("capacityPools"); err != nil {
		return nil, err
	}

	if service.Name, err = id.PopSegment("volumes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &service, nil
}
