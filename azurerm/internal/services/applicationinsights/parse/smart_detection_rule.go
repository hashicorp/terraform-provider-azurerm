package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SmartDetectionRuleId struct {
	SubscriptionId         string
	ResourceGroup          string
	ComponentName          string
	SmartDetectionRuleName string
}

func NewSmartDetectionRuleID(subscriptionId, resourceGroup, componentName, smartDetectionRuleName string) SmartDetectionRuleId {
	return SmartDetectionRuleId{
		SubscriptionId:         subscriptionId,
		ResourceGroup:          resourceGroup,
		ComponentName:          componentName,
		SmartDetectionRuleName: smartDetectionRuleName,
	}
}

func (id SmartDetectionRuleId) String() string {
	segments := []string{
		fmt.Sprintf("Smart Detection Rule Name %q", id.SmartDetectionRuleName),
		fmt.Sprintf("Component Name %q", id.ComponentName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Smart Detection Rule", segmentsStr)
}

func (id SmartDetectionRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/microsoft.insights/components/%s/SmartDetectionRule/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.ComponentName, id.SmartDetectionRuleName)
}

// SmartDetectionRuleID parses a SmartDetectionRule ID into an SmartDetectionRuleId struct
func SmartDetectionRuleID(input string) (*SmartDetectionRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SmartDetectionRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.ComponentName, err = id.PopSegment("components"); err != nil {
		return nil, err
	}
	if resourceId.SmartDetectionRuleName, err = id.PopSegment("SmartDetectionRule"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
