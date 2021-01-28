package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type EndpointCustomDomainId struct {
	SubscriptionId   string
	ResourceGroup    string
	ProfileName      string
	EndpointName     string
	CustomdomainName string
}

func NewEndpointCustomDomainID(subscriptionId, resourceGroup, profileName, endpointName, customdomainName string) EndpointCustomDomainId {
	return EndpointCustomDomainId{
		SubscriptionId:   subscriptionId,
		ResourceGroup:    resourceGroup,
		ProfileName:      profileName,
		EndpointName:     endpointName,
		CustomdomainName: customdomainName,
	}
}

func (id EndpointCustomDomainId) String() string {
	segments := []string{
		fmt.Sprintf("Customdomain Name %q", id.CustomdomainName),
		fmt.Sprintf("Endpoint Name %q", id.EndpointName),
		fmt.Sprintf("Profile Name %q", id.ProfileName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Endpoint Custom Domain", segmentsStr)
}

func (id EndpointCustomDomainId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cdn/profiles/%s/endpoints/%s/customdomains/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ProfileName, id.EndpointName, id.CustomdomainName)
}

// EndpointCustomDomainID parses a EndpointCustomDomain ID into an EndpointCustomDomainId struct
func EndpointCustomDomainID(input string) (*EndpointCustomDomainId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := EndpointCustomDomainId{
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
	if resourceId.EndpointName, err = id.PopSegment("endpoints"); err != nil {
		return nil, err
	}
	if resourceId.CustomdomainName, err = id.PopSegment("customdomains"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
