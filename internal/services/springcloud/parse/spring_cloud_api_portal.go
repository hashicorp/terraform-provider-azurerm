package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudAPIPortalId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
	ApiPortalName  string
}

func NewSpringCloudAPIPortalID(subscriptionId, resourceGroup, springName, apiPortalName string) SpringCloudAPIPortalId {
	return SpringCloudAPIPortalId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
		ApiPortalName:  apiPortalName,
	}
}

func (id SpringCloudAPIPortalId) String() string {
	segments := []string{
		fmt.Sprintf("Api Portal Name %q", id.ApiPortalName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud A P I Portal", segmentsStr)
}

func (id SpringCloudAPIPortalId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/Spring/%s/apiPortals/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.ApiPortalName)
}

// SpringCloudAPIPortalID parses a SpringCloudAPIPortal ID into an SpringCloudAPIPortalId struct
func SpringCloudAPIPortalID(input string) (*SpringCloudAPIPortalId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudAPIPortalId{
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
	if resourceId.ApiPortalName, err = id.PopSegment("apiPortals"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
