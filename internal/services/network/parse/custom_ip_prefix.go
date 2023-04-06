package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type CustomIpPrefixId struct {
	SubscriptionId      string
	ResourceGroup       string
	CustomIpPrefixeName string
}

func NewCustomIpPrefixID(subscriptionId, resourceGroup, customIpPrefixeName string) CustomIpPrefixId {
	return CustomIpPrefixId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		CustomIpPrefixeName: customIpPrefixeName,
	}
}

func (id CustomIpPrefixId) String() string {
	segments := []string{
		fmt.Sprintf("Custom Ip Prefixe Name %q", id.CustomIpPrefixeName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Custom Ip Prefix", segmentsStr)
}

func (id CustomIpPrefixId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/customIpPrefixes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CustomIpPrefixeName)
}

// CustomIpPrefixID parses a CustomIpPrefix ID into an CustomIpPrefixId struct
func CustomIpPrefixID(input string) (*CustomIpPrefixId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := CustomIpPrefixId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.CustomIpPrefixeName, err = id.PopSegment("customIpPrefixes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
