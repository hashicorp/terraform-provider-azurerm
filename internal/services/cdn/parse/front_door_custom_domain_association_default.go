package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontDoorCustomDomainAssociationDefaultId struct {
	SubscriptionId          string
	ResourceGroup           string
	ProfileName             string
	AfdEndpointName         string
	AssociatioinDefaultName string
	RouteName               string
	CustomDomainName        string
}

func NewFrontDoorCustomDomainAssociationDefaultID(subscriptionId, resourceGroup, profileName, afdEndpointName, associatioinDefaultName, routeName, customDomainName string) FrontDoorCustomDomainAssociationDefaultId {
	return FrontDoorCustomDomainAssociationDefaultId{
		SubscriptionId:          subscriptionId,
		ResourceGroup:           resourceGroup,
		ProfileName:             profileName,
		AfdEndpointName:         afdEndpointName,
		AssociatioinDefaultName: associatioinDefaultName,
		RouteName:               routeName,
		CustomDomainName:        customDomainName,
	}
}

func (id FrontDoorCustomDomainAssociationDefaultId) String() string {
	segments := []string{
		fmt.Sprintf("Custom Domain Name %q", id.CustomDomainName),
		fmt.Sprintf("Route Name %q", id.RouteName),
		fmt.Sprintf("Associatioin Default Name %q", id.AssociatioinDefaultName),
		fmt.Sprintf("Afd Endpoint Name %q", id.AfdEndpointName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Front Door Custom Domain Association Default", segmentsStr)
}

func (id FrontDoorCustomDomainAssociationDefaultId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/afdEndpoints/%s/associatioinDefault/%s/routes/%s/customDomains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociatioinDefaultName, id.RouteName, id.CustomDomainName)
}

// FrontDoorCustomDomainAssociationDefaultID parses a FrontDoorCustomDomainAssociationDefault ID into an FrontDoorCustomDomainAssociationDefaultId struct
func FrontDoorCustomDomainAssociationDefaultID(input string) (*FrontDoorCustomDomainAssociationDefaultId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontDoorCustomDomainAssociationDefaultId{
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
	if resourceId.AssociatioinDefaultName, err = id.PopSegment("associatioinDefault"); err != nil {
		return nil, err
	}
	if resourceId.RouteName, err = id.PopSegment("routes"); err != nil {
		return nil, err
	}
	if resourceId.CustomDomainName, err = id.PopSegment("customDomains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

// FrontDoorCustomDomainAssociationDefaultIDInsensitively parses an FrontDoorCustomDomainAssociationDefault ID into an FrontDoorCustomDomainAssociationDefaultId struct, insensitively
// This should only be used to parse an ID for rewriting, the FrontDoorCustomDomainAssociationDefaultID
// method should be used instead for validation etc.
//
// Whilst this may seem strange, this enables Terraform have consistent casing
// which works around issues in Core, whilst handling broken API responses.
func FrontDoorCustomDomainAssociationDefaultIDInsensitively(input string) (*FrontDoorCustomDomainAssociationDefaultId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontDoorCustomDomainAssociationDefaultId{
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

	// find the correct casing for the 'associatioinDefault' segment
	associatioinDefaultKey := "associatioinDefault"
	for key := range id.Path {
		if strings.EqualFold(key, associatioinDefaultKey) {
			associatioinDefaultKey = key
			break
		}
	}
	if resourceId.AssociatioinDefaultName, err = id.PopSegment(associatioinDefaultKey); err != nil {
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

	// find the correct casing for the 'customDomains' segment
	customDomainsKey := "customDomains"
	for key := range id.Path {
		if strings.EqualFold(key, customDomainsKey) {
			customDomainsKey = key
			break
		}
	}
	if resourceId.CustomDomainName, err = id.PopSegment(customDomainsKey); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
