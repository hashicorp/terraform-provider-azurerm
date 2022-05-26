package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudGatewayRouteConfigId struct {
	SubscriptionId  string
	ResourceGroup   string
	SpringName      string
	GatewayName     string
	RouteConfigName string
}

func NewSpringCloudGatewayRouteConfigID(subscriptionId, resourceGroup, springName, gatewayName, routeConfigName string) SpringCloudGatewayRouteConfigId {
	return SpringCloudGatewayRouteConfigId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		SpringName:      springName,
		GatewayName:     gatewayName,
		RouteConfigName: routeConfigName,
	}
}

func (id SpringCloudGatewayRouteConfigId) String() string {
	segments := []string{
		fmt.Sprintf("Route Config Name %q", id.RouteConfigName),
		fmt.Sprintf("Gateway Name %q", id.GatewayName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Gateway Route Config", segmentsStr)
}

func (id SpringCloudGatewayRouteConfigId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/Spring/%s/gateways/%s/routeConfigs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.GatewayName, id.RouteConfigName)
}

// SpringCloudGatewayRouteConfigID parses a SpringCloudGatewayRouteConfig ID into an SpringCloudGatewayRouteConfigId struct
func SpringCloudGatewayRouteConfigID(input string) (*SpringCloudGatewayRouteConfigId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudGatewayRouteConfigId{
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
	if resourceId.RouteConfigName, err = id.PopSegment("routeConfigs"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
