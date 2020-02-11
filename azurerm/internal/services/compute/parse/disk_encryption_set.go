package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DiskEncryptionSetId struct {
	ResourceGroup string
	Name          string
}

func DiskEncryptionSetID(input string) (*DiskEncryptionSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Disk Encryption Set ID %q: %+v", input, err)
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
