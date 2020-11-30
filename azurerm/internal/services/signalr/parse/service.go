package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ServiceId struct {
	SubscriptionId string
	ResourceGroup  string
	SignalRName    string
}

func NewServiceID(subscriptionId, resourceGroup, signalRName string) ServiceId {
	return ServiceId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SignalRName:    signalRName,
	}
}

func (id ServiceId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SignalRService/SignalR/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SignalRName)
}

// ServiceID parses a Service ID into an ServiceId struct
func ServiceID(input string) (*ServiceId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ServiceId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SignalRName, err = id.PopSegment("SignalR"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
