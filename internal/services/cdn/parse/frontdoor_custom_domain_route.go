package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontdoorCustomDomainRouteId struct {
	SubscriptionId  string
	ResourceGroup   string
	ProfileName     string
	AfdEndpointName string
	RouteName       string
	AssociationName string
}

func NewFrontdoorCustomDomainRouteID(subscriptionId, resourceGroup, profileName, afdEndpointName, routeName, associationName string) FrontdoorCustomDomainRouteId {
	return FrontdoorCustomDomainRouteId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		ProfileName:     profileName,
		AfdEndpointName: afdEndpointName,
		RouteName:       routeName,
		AssociationName: associationName,
	}
}

func (id FrontdoorCustomDomainRouteId) String() string {
	segments := []string{
		fmt.Sprintf("Association Name %q", id.AssociationName),
		fmt.Sprintf("Route Name %q", id.RouteName),
		fmt.Sprintf("Afd Endpoint Name %q", id.AfdEndpointName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Frontdoor Custom Domain Route", segmentsStr)
}

func (id FrontdoorCustomDomainRouteId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/afdEndpoints/%s/routes/%s/associations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.RouteName, id.AssociationName)
}

// FrontdoorCustomDomainRouteID parses a FrontdoorCustomDomainRoute ID into an FrontdoorCustomDomainRouteId struct
func FrontdoorCustomDomainRouteID(input string) (*FrontdoorCustomDomainRouteId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorCustomDomainRouteId{
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
	if resourceId.AssociationName, err = id.PopSegment("associations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// FrontdoorCustomDomainRouteIDInsensitively parses an FrontdoorCustomDomainRoute ID into an FrontdoorCustomDomainRouteId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontdoorCustomDomainRouteID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontdoorCustomDomainRouteIDInsensitively(input string) (*FrontdoorCustomDomainRouteId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorCustomDomainRouteId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	// find the correct casing for the 'profiles' segment
	profilesKey := "profiles"
	for key := range id.Path {
		if strings.EqualFold(key, profilesKey) {
			profilesKey = key
			break
		}
	}
	if resourceId.ProfileName, err = id.PopSegment(profilesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'afdEndpoints' segment
	afdEndpointsKey := "afdEndpoints"
	for key := range id.Path {
		if strings.EqualFold(key, afdEndpointsKey) {
			afdEndpointsKey = key
			break
		}
	}
	if resourceId.AfdEndpointName, err = id.PopSegment(afdEndpointsKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'routes' segment
	routesKey := "routes"
	for key := range id.Path {
		if strings.EqualFold(key, routesKey) {
			routesKey = key
			break
		}
	}
	if resourceId.RouteName, err = id.PopSegment(routesKey); err != nil {
		return nil, err
	}

	// find the correct casing for the 'associations' segment
	associationsKey := "associations"
	for key := range id.Path {
		if strings.EqualFold(key, associationsKey) {
			associationsKey = key
			break
		}
	}
	if resourceId.AssociationName, err = id.PopSegment(associationsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
