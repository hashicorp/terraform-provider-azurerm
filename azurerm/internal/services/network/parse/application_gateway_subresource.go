package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ApplicationGatewayHTTPListenerId struct {
	SubscriptionId         string
	ResourceGroup          string
	ApplicationGatewayName string
	HTTPListener           string
}

// ApplicationGatewayHTTPListenerID parses a ApplicationGateway Path Rule ID into an ApplicationGatewayHTTPListenerId struct
func ApplicationGatewayHTTPListenerID(input string) (*ApplicationGatewayHTTPListenerId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ApplicationGatewayHTTPListenerId{
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

	if resourceId.HTTPListener, err = id.PopSegment("httpListeners"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}

type ApplicationGatewayPathRuleId struct {
	SubscriptionId         string
	ResourceGroup          string
	ApplicationGatewayName string
	URLPathMap             string
	PathRule               string
}

// ApplicationGatewayPathRuleID parses a ApplicationGateway Path Rule ID into an ApplicationGatewayPathRuleId struct
func ApplicationGatewayPathRuleID(input string) (*ApplicationGatewayPathRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ApplicationGatewayPathRuleId{
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

	if resourceId.URLPathMap, err = id.PopSegment("urlPathMaps"); err != nil {
		return nil, err
	}

	if resourceId.PathRule, err = id.PopSegment("pathRules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
