package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type FrontDoorRuleSetAssociationId struct {
	SubscriptionId  string
	ResourceGroup   string
	ProfileName     string
	AfdEndpointName string
	AssociationName string
}

func NewFrontDoorRuleSetAssociationID(subscriptionId, resourceGroup, profileName, afdEndpointName, associationName string) FrontDoorRuleSetAssociationId {
	return FrontDoorRuleSetAssociationId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		ProfileName:     profileName,
		AfdEndpointName: afdEndpointName,
		AssociationName: associationName,
	}
}

func (id FrontDoorRuleSetAssociationId) String() string {
	segments := []string{
		fmt.Sprintf("Association Name %q", id.AssociationName),
		fmt.Sprintf("Afd Endpoint Name %q", id.AfdEndpointName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Front Door Rule Set Association", segmentsStr)
}

func (id FrontDoorRuleSetAssociationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/afdEndpoints/%s/associations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.AfdEndpointName, id.AssociationName)
}

// FrontDoorRuleSetAssociationID parses a FrontDoorRuleSetAssociation ID into an FrontDoorRuleSetAssociationId struct
func FrontDoorRuleSetAssociationID(input string) (*FrontDoorRuleSetAssociationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := FrontDoorRuleSetAssociationId{
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
	if resourceId.AssociationName, err = id.PopSegment("associations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
