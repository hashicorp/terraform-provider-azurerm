package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApplicationGatewayURLPathMapPathRuleId struct {
	SubscriptionId         string
	ResourceGroup          string
	ApplicationGatewayName string
	UrlPathMapName         string
	PathRuleName           string
}

func NewApplicationGatewayURLPathMapPathRuleID(subscriptionId, resourceGroup, applicationGatewayName, urlPathMapName, pathRuleName string) ApplicationGatewayURLPathMapPathRuleId {
	return ApplicationGatewayURLPathMapPathRuleId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		ApplicationGatewayName: applicationGatewayName,
		UrlPathMapName:         urlPathMapName,
		PathRuleName:           pathRuleName,
	}
}

func (id ApplicationGatewayURLPathMapPathRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Path Rule Name %q", id.PathRuleName),
		fmt.Sprintf("Url Path Map Name %q", id.UrlPathMapName),
		fmt.Sprintf("Application Gateway Name %q", id.ApplicationGatewayName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Application Gateway U R L Path Map Path Rule", segmentsStr)
}

func (id ApplicationGatewayURLPathMapPathRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/applicationGateways/%s/urlPathMaps/%s/pathRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ApplicationGatewayName, id.UrlPathMapName, id.PathRuleName)
}

// ApplicationGatewayURLPathMapPathRuleID parses a ApplicationGatewayURLPathMapPathRule ID into an ApplicationGatewayURLPathMapPathRuleId struct
func ApplicationGatewayURLPathMapPathRuleID(input string) (*ApplicationGatewayURLPathMapPathRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ApplicationGatewayURLPathMapPathRuleId{
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
	if resourceId.UrlPathMapName, err = id.PopSegment("urlPathMaps"); err != nil {
		return nil, err
	}
	if resourceId.PathRuleName, err = id.PopSegment("pathRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
