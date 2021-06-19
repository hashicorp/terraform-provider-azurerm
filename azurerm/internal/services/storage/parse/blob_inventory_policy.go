package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type BlobInventoryPolicyId struct {
	SubscriptionId      string
	ResourceGroup       string
	StorageAccountName  string
	InventoryPolicyName string
}

func NewBlobInventoryPolicyID(subscriptionId, resourceGroup, storageAccountName, inventoryPolicyName string) BlobInventoryPolicyId {
	return BlobInventoryPolicyId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		StorageAccountName:  storageAccountName,
		InventoryPolicyName: inventoryPolicyName,
	}
}

func (id BlobInventoryPolicyId) String() string {
	segments := []string{
		fmt.Sprintf("Inventory Policy Name %q", id.InventoryPolicyName),
		fmt.Sprintf("Storage Account Name %q", id.StorageAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Blob Inventory Policy", segmentsStr)
}

func (id BlobInventoryPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/inventoryPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.StorageAccountName, id.InventoryPolicyName)
}

// BlobInventoryPolicyID parses a BlobInventoryPolicy ID into an BlobInventoryPolicyId struct
func BlobInventoryPolicyID(input string) (*BlobInventoryPolicyId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := BlobInventoryPolicyId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.StorageAccountName, err = id.PopSegment("storageAccounts"); err != nil {
		return nil, err
	}
	if resourceId.InventoryPolicyName, err = id.PopSegment("inventoryPolicies"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
