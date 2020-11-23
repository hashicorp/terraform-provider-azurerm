package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type DiskEncryptionSetId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewDiskEncryptionSetId(subscriptionId, resourceGroup, name string) DiskEncryptionSetId {
	return DiskEncryptionSetId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id DiskEncryptionSetId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/diskEncryptionSets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func DiskEncryptionSetID(input string) (*DiskEncryptionSetId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("unable to parse Disk Encryption Set ID %q: %+v", input, err)
	}

	encryptionSetId := DiskEncryptionSetId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if encryptionSetId.Name, err = id.PopSegment("diskEncryptionSets"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &encryptionSetId, nil
}
