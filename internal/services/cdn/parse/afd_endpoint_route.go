package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AfdEndpointRouteId struct {
	SubscriptionId  string
	ResourceGroup   string
	ProfileName     string
	AfdEndpointName string
	RouteName       string
}

func NewAfdEndpointRouteID(subscriptionId, resourceGroup, profileName, afdEndpointName, routeName string) AfdEndpointRouteId {
	return AfdEndpointRouteId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		ProfileName:     profileName,
		AfdEndpointName: afdEndpointName,
		RouteName:       routeName,
	}
}

func (id AfdEndpointRouteId) String() string {
	segments := []string{
		fmt.Sprintf("Route Name %q", id.RouteName),
		fmt.Sprintf("Afd Endpoint Name %q", id.AfdEndpointName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Afd Endpoint Route", segmentsStr)
}

func (id AfdEndpointRouteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/afdEndpoints/%s/routes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName)
}

// AfdEndpointRouteID parses a AfdEndpointRoute ID into an AfdEndpointRouteId struct
func AfdEndpointRouteID(input string) (*AfdEndpointRouteId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AfdEndpointRouteId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ProfileName, err = id.PopSegment("profiles"); err != nil {
		return nil, err
	}
	if resourceId.AfdEndpointName, err = id.PopSegment("afdEndpoints"); err != nil {
		return nil, err
	}
	if resourceId.RouteName, err = id.PopSegment("routes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
