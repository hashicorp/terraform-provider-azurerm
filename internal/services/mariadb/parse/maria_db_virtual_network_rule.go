package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
)

type MariaDBVirtualNetworkRuleId struct {
	SubscriptionId         string
	ResourceGroup          string
	ServerName             string
	VirtualNetworkRuleName string
}

func NewMariaDBVirtualNetworkRuleID(subscriptionId, resourceGroup, serverName, virtualNetworkRuleName string) MariaDBVirtualNetworkRuleId {
	return MariaDBVirtualNetworkRuleId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		ServerName:             serverName,
		VirtualNetworkRuleName: virtualNetworkRuleName,
	}
}

func (id MariaDBVirtualNetworkRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Virtual Network Rule Name %q", id.VirtualNetworkRuleName),
		fmt.Sprintf("Server Name %q", id.ServerName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Maria D B Virtual Network Rule", segmentsStr)
}

func (id MariaDBVirtualNetworkRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforMariaDB/servers/%s/virtualNetworkRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ServerName, id.VirtualNetworkRuleName)
}

// MariaDBVirtualNetworkRuleID parses a MariaDBVirtualNetworkRule ID into an MariaDBVirtualNetworkRuleId struct
func MariaDBVirtualNetworkRuleID(input string) (*MariaDBVirtualNetworkRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := MariaDBVirtualNetworkRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ServerName, err = id.PopSegment("servers"); err != nil {
		return nil, err
	}
	if resourceId.VirtualNetworkRuleName, err = id.PopSegment("virtualNetworkRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
