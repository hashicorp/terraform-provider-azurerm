package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type CapacityPoolId struct {
	SubscriptionId    string
	ResourceGroup     string
	NetAppAccountName string
	Name              string
}

func NewCapacityPoolID(subscriptionId, resourceGroup, netAppAccountName, name string) CapacityPoolId {
	return CapacityPoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroup:     resourceGroup,
		NetAppAccountName: netAppAccountName,
		Name:              name,
	}
}

func (id CapacityPoolId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NetApp/netAppAccounts/%s/capacityPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetAppAccountName, id.Name)
}

// CapacityPoolID parses a CapacityPool ID into an CapacityPoolId struct
func CapacityPoolID(input string) (*CapacityPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CapacityPoolId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.NetAppAccountName, err = id.PopSegment("netAppAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("capacityPools"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
