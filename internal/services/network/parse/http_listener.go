package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type HttpListenerId struct {
	SubscriptionId         string
	ResourceGroup          string
	ApplicationGatewayName string
	Name                   string
}

func NewHttpListenerID(subscriptionId, resourceGroup, applicationGatewayName, name string) HttpListenerId {
	return HttpListenerId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		ApplicationGatewayName: applicationGatewayName,
		Name:                   name,
	}
}

func (id HttpListenerId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Application Gateway Name %q", id.ApplicationGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Http Listener", segmentsStr)
}

func (id HttpListenerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/httpListeners/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGatewayName, id.Name)
}

// HttpListenerID parses a HttpListener ID into an HttpListenerId struct
func HttpListenerID(input string) (*HttpListenerId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := HttpListenerId{
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
	if resourceId.Name, err = id.PopSegment("httpListeners"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
