package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

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

func (id CapacityPoolId) String() string {
	segments := []string{
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
		fmt.Sprintf("Net App Account Name %q", id.NetAppAccountName),
		fmt.Sprintf("Name %q", id.Name),
	}
	return strings.Join(segments, " / ")
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

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
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
