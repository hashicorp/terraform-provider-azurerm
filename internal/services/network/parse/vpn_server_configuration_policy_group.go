package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type VpnServerConfigurationPolicyGroupId struct {
	SubscriptionId               string
	ResourceGroup                string
	VpnServerConfigurationName   string
	ConfigurationPolicyGroupName string
}

func NewVpnServerConfigurationPolicyGroupID(subscriptionId, resourceGroup, vpnServerConfigurationName, configurationPolicyGroupName string) VpnServerConfigurationPolicyGroupId {
	return VpnServerConfigurationPolicyGroupId{
		SubscriptionId:               subscriptionId,
		ResourceGroup:                resourceGroup,
		VpnServerConfigurationName:   vpnServerConfigurationName,
		ConfigurationPolicyGroupName: configurationPolicyGroupName,
	}
}

func (id VpnServerConfigurationPolicyGroupId) String() string {
	segments := []string{
		fmt.Sprintf("Configuration Policy Group Name %q", id.ConfigurationPolicyGroupName),
		fmt.Sprintf("Vpn Server Configuration Name %q", id.VpnServerConfigurationName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Vpn Server Configuration Policy Group", segmentsStr)
}

func (id VpnServerConfigurationPolicyGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnServerConfigurations/%s/configurationPolicyGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VpnServerConfigurationName, id.ConfigurationPolicyGroupName)
}

// VpnServerConfigurationPolicyGroupID parses a VpnServerConfigurationPolicyGroup ID into an VpnServerConfigurationPolicyGroupId struct
func VpnServerConfigurationPolicyGroupID(input string) (*VpnServerConfigurationPolicyGroupId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := VpnServerConfigurationPolicyGroupId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.VpnServerConfigurationName, err = id.PopSegment("vpnServerConfigurations"); err != nil {
		return nil, err
	}
	if resourceId.ConfigurationPolicyGroupName, err = id.PopSegment("configurationPolicyGroups"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
