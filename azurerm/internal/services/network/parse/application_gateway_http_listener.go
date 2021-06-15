package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApplicationGatewayHTTPListenerId struct {
	SubscriptionId         string
	ResourceGroup          string
	ApplicationGatewayName string
	HttpListenerName       string
}

func NewApplicationGatewayHTTPListenerID(subscriptionId, resourceGroup, applicationGatewayName, httpListenerName string) ApplicationGatewayHTTPListenerId {
	return ApplicationGatewayHTTPListenerId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		ApplicationGatewayName: applicationGatewayName,
		HttpListenerName:       httpListenerName,
	}
}

func (id ApplicationGatewayHTTPListenerId) String() string {
	segments := []string{
		fmt.Sprintf("Http Listener Name %q", id.HttpListenerName),
		fmt.Sprintf("Application Gateway Name %q", id.ApplicationGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Application Gateway H T T P Listener", segmentsStr)
}

func (id ApplicationGatewayHTTPListenerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/httpListeners/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGatewayName, id.HttpListenerName)
}

// ApplicationGatewayHTTPListenerID parses a ApplicationGatewayHTTPListener ID into an ApplicationGatewayHTTPListenerId struct
func ApplicationGatewayHTTPListenerID(input string) (*ApplicationGatewayHTTPListenerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ApplicationGatewayHTTPListenerId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ApplicationGatewayName, err = id.PopSegment("applicationGateways"); err != nil {
		return nil, err
	}
	if resourceId.HttpListenerName, err = id.PopSegment("httpListeners"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
