package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type AfdCustomDomainId struct {
	SubscriptionId   string
	ResourceGroup    string
	ProfileName      string
	CustomDomainName string
}

func NewAfdCustomDomainID(subscriptionId, resourceGroup, profileName, customDomainName string) AfdCustomDomainId {
	return AfdCustomDomainId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		ProfileName:      profileName,
		CustomDomainName: customDomainName,
	}
}

func (id AfdCustomDomainId) String() string {
	segments := []string{
		fmt.Sprintf("Custom Domain Name %q", id.CustomDomainName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Afd Custom Domain", segmentsStr)
}

func (id AfdCustomDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/customDomains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.CustomDomainName)
}

// AfdCustomDomainID parses a AfdCustomDomain ID into an AfdCustomDomainId struct
func AfdCustomDomainID(input string) (*AfdCustomDomainId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := AfdCustomDomainId{
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
	if resourceId.CustomDomainName, err = id.PopSegment("customDomains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
