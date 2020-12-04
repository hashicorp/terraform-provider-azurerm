package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type LoadBalancerFrontendIpConfigurationId struct {
	SubscriptionId              string
	ResourceGroup               string
	LoadBalancerName            string
	FrontendIPConfigurationName string
}

func NewLoadBalancerFrontendIpConfigurationID(subscriptionId, resourceGroup, loadBalancerName, frontendIPConfigurationName string) LoadBalancerFrontendIpConfigurationId {
	return LoadBalancerFrontendIpConfigurationId{
		SubscriptionId:              subscriptionId,
		ResourceGroup:               resourceGroup,
		LoadBalancerName:            loadBalancerName,
		FrontendIPConfigurationName: frontendIPConfigurationName,
	}
}

func (id LoadBalancerFrontendIpConfigurationId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/loadBalancers/%s/frontendIPConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.LoadBalancerName, id.FrontendIPConfigurationName)
}

// LoadBalancerFrontendIpConfigurationID parses a LoadBalancerFrontendIpConfiguration ID into an LoadBalancerFrontendIpConfigurationId struct
func LoadBalancerFrontendIpConfigurationID(input string) (*LoadBalancerFrontendIpConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := LoadBalancerFrontendIpConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.LoadBalancerName, err = id.PopSegment("loadBalancers"); err != nil {
		return nil, err
	}
	if resourceId.FrontendIPConfigurationName, err = id.PopSegment("frontendIPConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
