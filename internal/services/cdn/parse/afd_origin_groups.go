package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AfdOriginGroupsId struct {
	SubscriptionId  string
	ResourceGroup   string
	ProfileName     string
	OriginGroupName string
}

func NewAfdOriginGroupsID(subscriptionId, resourceGroup, profileName, originGroupName string) AfdOriginGroupsId {
	return AfdOriginGroupsId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		ProfileName:     profileName,
		OriginGroupName: originGroupName,
	}
}

func (id AfdOriginGroupsId) String() string {
	segments := []string{
		fmt.Sprintf("Origin Group Name %q", id.OriginGroupName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Afd Origin Groups", segmentsStr)
}

func (id AfdOriginGroupsId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/originGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.OriginGroupName)
}

// AfdOriginGroupsID parses a AfdOriginGroups ID into an AfdOriginGroupsId struct
func AfdOriginGroupsID(input string) (*AfdOriginGroupsId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AfdOriginGroupsId{
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
	if resourceId.OriginGroupName, err = id.PopSegment("originGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
