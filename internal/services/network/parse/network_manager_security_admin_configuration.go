package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type NetworkManagerSecurityAdminConfigurationId struct {
	SubscriptionId                 string
	ResourceGroup                  string
	NetworkManagerName             string
	SecurityAdminConfigurationName string
}

func NewNetworkManagerSecurityAdminConfigurationID(subscriptionId, resourceGroup, networkManagerName, securityAdminConfigurationName string) NetworkManagerSecurityAdminConfigurationId {
	return NetworkManagerSecurityAdminConfigurationId{
		SubscriptionId:                 subscriptionId,
		ResourceGroup:                  resourceGroup,
		NetworkManagerName:             networkManagerName,
		SecurityAdminConfigurationName: securityAdminConfigurationName,
	}
}

func (id NetworkManagerSecurityAdminConfigurationId) String() string {
	segments := []string{
		fmt.Sprintf("Security Admin Configuration Name %q", id.SecurityAdminConfigurationName),
		fmt.Sprintf("Network Manager Name %q", id.NetworkManagerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Network Manager Security Admin Configuration", segmentsStr)
}

func (id NetworkManagerSecurityAdminConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/networkManagers/%s/securityAdminConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.NetworkManagerName, id.SecurityAdminConfigurationName)
}

// NetworkManagerSecurityAdminConfigurationID parses a NetworkManagerSecurityAdminConfiguration ID into an NetworkManagerSecurityAdminConfigurationId struct
func NetworkManagerSecurityAdminConfigurationID(input string) (*NetworkManagerSecurityAdminConfigurationId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := NetworkManagerSecurityAdminConfigurationId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.NetworkManagerName, err = id.PopSegment("networkManagers"); err != nil {
		return nil, err
	}
	if resourceId.SecurityAdminConfigurationName, err = id.PopSegment("securityAdminConfigurations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
