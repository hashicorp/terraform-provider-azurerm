package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type RouteId struct {
	SubscriptionId string
	ResourceGroup  string
	RouteTableName string
	Name           string
}

func NewRouteID(subscriptionId, resourceGroup, routeTableName, name string) RouteId {
	return RouteId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		RouteTableName: routeTableName,
		Name:           name,
	}
}

func (id RouteId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Route Table Name %q", id.RouteTableName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Route", segmentsStr)
}

func (id RouteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/routeTables/%s/routes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.RouteTableName, id.Name)
}

// RouteID parses a Route ID into an RouteId struct
func RouteID(input string) (*RouteId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := RouteId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.RouteTableName, err = id.PopSegment("routeTables"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("routes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
