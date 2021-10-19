package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type AuthorizationId struct {
	SubscriptionId          string
	ResourceGroup           string
	ExpressRouteCircuitName string
	Name                    string
}

func NewAuthorizationID(subscriptionId, resourceGroup, expressRouteCircuitName, name string) AuthorizationId {
	return AuthorizationId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ExpressRouteCircuitName: expressRouteCircuitName,
		Name:                    name,
	}
}

func (id AuthorizationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Express Route Circuit Name %q", id.ExpressRouteCircuitName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Authorization", segmentsStr)
}

func (id AuthorizationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/%ss/expressRouteCircuit/authorizations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ExpressRouteCircuitName, id.Name)
}

// AuthorizationID parses a Authorization ID into an AuthorizationId struct
func AuthorizationID(input string) (*AuthorizationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AuthorizationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ExpressRouteCircuitName, err = id.PopSegment("expressRouteCircuits"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("authorizations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
