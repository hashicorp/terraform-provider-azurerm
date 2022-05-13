package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudGatewayId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
	GatewayName    string
}

func NewSpringCloudGatewayID(subscriptionId, resourceGroup, springName, gatewayName string) SpringCloudGatewayId {
	return SpringCloudGatewayId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
		GatewayName:    gatewayName,
	}
}

func (id SpringCloudGatewayId) String() string {
	segments := []string{
		fmt.Sprintf("Gateway Name %q", id.GatewayName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Gateway", segmentsStr)
}

func (id SpringCloudGatewayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/Spring/%s/gateways/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.GatewayName)
}

// SpringCloudGatewayID parses a SpringCloudGateway ID into an SpringCloudGatewayId struct
func SpringCloudGatewayID(input string) (*SpringCloudGatewayId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudGatewayId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.SpringName, err = id.PopSegment("Spring"); err != nil {
		return nil, err
	}
	if resourceId.GatewayName, err = id.PopSegment("gateways"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
