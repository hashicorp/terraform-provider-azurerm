package validate

import "github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"

// IsLogAnalyticsAutomationID parses a resource ID an returns a bool indicating if it is a valid LogAnalyticsAutomationID
func IsLogAnalyticsAutomationID(input string) bool {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return false
	}

	if _, err = id.PopSegment("automationAccounts"); err != nil {
		return false
	}

	return true
}
