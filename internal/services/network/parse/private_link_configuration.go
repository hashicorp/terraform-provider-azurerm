package parse

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type PrivateLinkConfigurationId struct {
	SubscriptionId         string
	ResourceGroup          string
	ApplicationGatewayName string
	Name                   string
}

func NewPrivateLinkConfigurationID(subscriptionId, resourceGroup, applicationGatewayName, name string) PrivateLinkConfigurationId {
	return PrivateLinkConfigurationId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		ApplicationGatewayName: applicationGatewayName,
		Name:                   name,
	}
}

func (id PrivateLinkConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Application Gateway Name %q", id.ApplicationGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Private Link Configuration", segmentsStr)
}

func (id PrivateLinkConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/privateLinkConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGatewayName, id.Name)
}

// PrivateLinkConfigurationID parses a PrivateLinkConfiguration ID into an PrivateLinkConfigurationId struct
func PrivateLinkConfigurationID(input string) (*PrivateLinkConfigurationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := PrivateLinkConfigurationId{
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
	if resourceId.Name, err = id.PopSegment("privateLinkConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
