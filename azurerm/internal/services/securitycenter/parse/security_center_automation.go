package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SecurityCenterAutomationId struct {
	AutomationName string
	ResourceGroup  string
}

func SecurityCenterAutomationID(input string) (*SecurityCenterAutomationId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse Security Center Automation ID %q: %+v", input, err)
	}

	automation := SecurityCenterAutomationId{
		ResourceGroup: id.ResourceGroup,
	}

	if automation.AutomationName, err = id.PopSegment("automations"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &automation, nil
}
