package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type HostPoolId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewHostPoolID(subscriptionId, resourceGroup, name string) HostPoolId {
	return HostPoolId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id HostPoolId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/hostPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

func HostPoolID(input string) (*HostPoolId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := HostPoolId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	// find the correct casing for the `hostPools` segment
	hostPoolsKey := "hostPools"
	for key := range id.Path {
		if strings.EqualFold(key, hostPoolsKey) {
			hostPoolsKey = key
			break
		}
	}
	if resourceId.Name, err = id.PopSegment(hostPoolsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
