package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type ApplicationGatewayPrivateLinkConfigurationId struct {
	SubscriptionId               string
	ResourceGroup                string
	ApplicationGatewayName       string
	PrivateLinkConfigurationName string
}

func NewApplicationGatewayPrivateLinkConfigurationID(subscriptionId, resourceGroup, applicationGatewayName, privateLinkConfigurationName string) ApplicationGatewayPrivateLinkConfigurationId {
	return ApplicationGatewayPrivateLinkConfigurationId{
		SubscriptionId:               subscriptionId,
		ResourceGroup:                resourceGroup,
		ApplicationGatewayName:       applicationGatewayName,
		PrivateLinkConfigurationName: privateLinkConfigurationName,
	}
}

func (id ApplicationGatewayPrivateLinkConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Private Link Configuration Name %q", id.PrivateLinkConfigurationName),
		fmt.Sprintf("Application Gateway Name %q", id.ApplicationGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Application Gateway Private Link Configuration", segmentsStr)
}

func (id ApplicationGatewayPrivateLinkConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/privateLinkConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGatewayName, id.PrivateLinkConfigurationName)
}

// ApplicationGatewayPrivateLinkConfigurationID parses a ApplicationGatewayPrivateLinkConfiguration ID into an ApplicationGatewayPrivateLinkConfigurationId struct
func ApplicationGatewayPrivateLinkConfigurationID(input string) (*ApplicationGatewayPrivateLinkConfigurationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ApplicationGatewayPrivateLinkConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ApplicationGatewayName, err = id.PopSegment("applicationGateways"); err != nil {
		return nil, err
	}
	if resourceId.PrivateLinkConfigurationName, err = id.PopSegment("privateLinkConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
