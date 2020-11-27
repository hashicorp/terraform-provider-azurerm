package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SmartDetectorAlertRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	Name           string
}

func NewSmartDetectorAlertRuleID(subscriptionId, resourceGroup, name string) SmartDetectorAlertRuleId {
	return SmartDetectorAlertRuleId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		Name:           name,
	}
}

func (id SmartDetectorAlertRuleId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AlertsManagement/smartdetectoralertrules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.Name)
}

// SmartDetectorAlertRuleID parses a SmartDetectorAlertRule ID into an SmartDetectorAlertRuleId struct
func SmartDetectorAlertRuleID(input string) (*SmartDetectorAlertRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SmartDetectorAlertRuleId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.Name, err = id.PopSegment("smartdetectoralertrules"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
