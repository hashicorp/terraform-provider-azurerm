package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SecurityCenterAutomationId struct {
	SubscriptionId string
	AutomationName string
}

func NewSecurityCenterAutomationID(subscriptionId, automationName string) SecurityCenterAutomationId {
	return SecurityCenterAutomationId{
		SubscriptionId: subscriptionId,
		AutomationName: automationName,
	}
}

func (id SecurityCenterAutomationId) String() string {
	segments := []string{
		fmt.Sprintf("Automation Name %q", id.AutomationName),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Security Center Automation", segmentsStr)
}

func (id SecurityCenterAutomationId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Security/automations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.AutomationName)
}

// SecurityCenterAutomationID parses a SecurityCenterAutomation ID into an SecurityCenterAutomationId struct
func SecurityCenterAutomationID(input string) (*SecurityCenterAutomationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SecurityCenterAutomationId{
		SubscriptionId: id.SubscriptionID,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.AutomationName, err = id.PopSegment("automations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
