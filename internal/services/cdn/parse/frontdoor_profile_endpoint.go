package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontdoorProfileEndpointId struct {
	SubscriptionId  string
	ResourceGroup   string
	ProfileName     string
	AfdEndpointName string
}

func NewFrontdoorProfileEndpointID(subscriptionId, resourceGroup, profileName, afdEndpointName string) FrontdoorProfileEndpointId {
	return FrontdoorProfileEndpointId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		ProfileName:     profileName,
		AfdEndpointName: afdEndpointName,
	}
}

func (id FrontdoorProfileEndpointId) String() string {
	segments := []string{
		fmt.Sprintf("Afd Endpoint Name %q", id.AfdEndpointName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Frontdoor Profile Endpoint", segmentsStr)
}

func (id FrontdoorProfileEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/afdEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AfdEndpointName)
}

// FrontdoorProfileEndpointID parses a FrontdoorProfileEndpoint ID into an FrontdoorProfileEndpointId struct
func FrontdoorProfileEndpointID(input string) (*FrontdoorProfileEndpointId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontdoorProfileEndpointId{
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

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
