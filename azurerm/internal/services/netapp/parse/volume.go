package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type VolumeId struct {
	ResourceGroup     string
	NetAppAccountName string
	CapacityPoolName  string
	Name              string
}

func VolumeID(input string) (*VolumeId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse NetApp Volume ID %q: %+v", input, err)
	}

	service := VolumeId{
		ResourceGroup: id.ResourceGroup,
	}

	if service.NetAppAccountName, err = id.PopSegment("netAppAccounts"); err != nil {
		return nil, err
	}

	if service.CapacityPoolName, err = id.PopSegment("capacityPools"); err != nil {
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
