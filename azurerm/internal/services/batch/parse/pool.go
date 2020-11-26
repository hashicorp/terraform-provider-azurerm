package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PoolId struct {
	SubscriptionId   string
	ResourceGroup    string
	BatchAccountName string
	Name             string
}

func NewPoolID(subscriptionId, resourceGroup, batchAccountName, name string) PoolId {
	return PoolId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		BatchAccountName: batchAccountName,
		Name:             name,
	}
}

func (id PoolId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Batch/batchAccounts/%s/pools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.BatchAccountName, id.Name)
}

// PoolID parses a Pool ID into an PoolId struct
func PoolID(input string) (*PoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PoolId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.BatchAccountName, err = id.PopSegment("batchAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("pools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
