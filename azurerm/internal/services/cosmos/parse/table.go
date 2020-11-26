package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type TableId struct {
	SubscriptionId      string
	ResourceGroup       string
	DatabaseAccountName string
	Name                string
}

func NewTableID(subscriptionId, resourceGroup, databaseAccountName, name string) TableId {
	return TableId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		DatabaseAccountName: databaseAccountName,
		Name:                name,
	}
}

func (id TableId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/tables/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.DatabaseAccountName, id.Name)
}

func TableID(input string) (*TableId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := TableId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.DatabaseAccountName, err = id.PopSegment("databaseAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("tables"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
