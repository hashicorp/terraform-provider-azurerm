package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AfdOriginsId struct {
	SubscriptionId  string
	ResourceGroup   string
	ProfileName     string
	OriginGroupName string
	OriginName      string
}

func NewAfdOriginsID(subscriptionId, resourceGroup, profileName, originGroupName, originName string) AfdOriginsId {
	return AfdOriginsId{
		SubscriptionId:  subscriptionId,
		ResourceGroup:   resourceGroup,
		ProfileName:     profileName,
		OriginGroupName: originGroupName,
		OriginName:      originName,
	}
}

func (id AfdOriginsId) String() string {
	segments := []string{
		fmt.Sprintf("Origin Name %q", id.OriginName),
		fmt.Sprintf("Origin Group Name %q", id.OriginGroupName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Afd Origins", segmentsStr)
}

func (id AfdOriginsId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/originGroups/%s/origins/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.OriginGroupName, id.OriginName)
}

// AfdOriginsID parses a AfdOrigins ID into an AfdOriginsId struct
func AfdOriginsID(input string) (*AfdOriginsId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AfdOriginsId{
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
	if resourceId.OriginName, err = id.PopSegment("origins"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
