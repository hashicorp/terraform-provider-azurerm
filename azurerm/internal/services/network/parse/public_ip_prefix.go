package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type PublicIpPrefixId struct {
	SubscriptionId      string
	ResourceGroup       string
	PublicIPPrefixeName string
}

func NewPublicIpPrefixID(subscriptionId, resourceGroup, publicIPPrefixeName string) PublicIpPrefixId {
	return PublicIpPrefixId{
		SubscriptionId:      subscriptionId,
		ResourceGroup:       resourceGroup,
		PublicIPPrefixeName: publicIPPrefixeName,
	}
}

func (id PublicIpPrefixId) String() string {
	segments := []string{
		fmt.Sprintf("Public I P Prefixe Name %q", id.PublicIPPrefixeName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Public Ip Prefix", segmentsStr)
}

func (id PublicIpPrefixId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/publicIPPrefixes/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.PublicIPPrefixeName)
}

// PublicIpPrefixID parses a PublicIpPrefix ID into an PublicIpPrefixId struct
func PublicIpPrefixID(input string) (*PublicIpPrefixId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PublicIpPrefixId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.PublicIPPrefixeName, err = id.PopSegment("publicIPPrefixes"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
