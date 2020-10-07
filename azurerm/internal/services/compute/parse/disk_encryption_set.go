package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DiskEncryptionSetId struct {
	ResourceGroup string
	Name          string
}

func NewDiskEncryptionSetId(resourceGroup, name string) DiskEncryptionSetId {
	return DiskEncryptionSetId{
		ResourceGroup: resourceGroup,
		Name:          name,
	}
}

func (id DiskEncryptionSetId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/diskEncryptionSets/%s", subscriptionId, id.ResourceGroup, id.Name)
}

func DiskEncryptionSetID(input string) (*DiskEncryptionSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Disk Encryption Set ID %q: %+v", input, err)
	}

	encryptionSetId := DiskEncryptionSetId{
		ResourceGroup: id.ResourceGroup,
	}

	if encryptionSetId.Name, err = id.PopSegment("diskEncryptionSets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &encryptionSetId, nil
}
