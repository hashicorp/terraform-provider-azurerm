package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type SpringCloudGatewayCustomDomainId struct {
	SubscriptionId string
	ResourceGroup  string
	SpringName     string
	GatewayName    string
	DomainName     string
}

func NewSpringCloudGatewayCustomDomainID(subscriptionId, resourceGroup, springName, gatewayName, domainName string) SpringCloudGatewayCustomDomainId {
	return SpringCloudGatewayCustomDomainId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		SpringName:     springName,
		GatewayName:    gatewayName,
		DomainName:     domainName,
	}
}

func (id SpringCloudGatewayCustomDomainId) String() string {
	segments := []string{
		fmt.Sprintf("Domain Name %q", id.DomainName),
		fmt.Sprintf("Gateway Name %q", id.GatewayName),
		fmt.Sprintf("Spring Name %q", id.SpringName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Spring Cloud Gateway Custom Domain", segmentsStr)
}

func (id SpringCloudGatewayCustomDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AppPlatform/Spring/%s/gateways/%s/domains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.SpringName, id.GatewayName, id.DomainName)
}

// SpringCloudGatewayCustomDomainID parses a SpringCloudGatewayCustomDomain ID into an SpringCloudGatewayCustomDomainId struct
func SpringCloudGatewayCustomDomainID(input string) (*SpringCloudGatewayCustomDomainId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SpringCloudGatewayCustomDomainId{
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
	if resourceId.GatewayName, err = id.PopSegment("gateways"); err != nil {
		return nil, err
	}
	if resourceId.DomainName, err = id.PopSegment("domains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
