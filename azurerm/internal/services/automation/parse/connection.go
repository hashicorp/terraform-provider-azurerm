package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type ConnectionId struct {
	SubscriptionId        string
	ResourceGroup         string
	AutomationAccountName string
	Name                  string
}

func NewConnectionID(subscriptionId, resourceGroup, automationAccountName, name string) ConnectionId {
	return ConnectionId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		AutomationAccountName: automationAccountName,
		Name:                  name,
	}
}

func (id ConnectionId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Automation Account Name %q", id.AutomationAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Connection", segmentsStr)
}

func (id ConnectionId) ID(_ string) string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/connections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AutomationAccountName, id.Name)
}

// ConnectionID parses a Connection ID into an ConnectionId struct
func ConnectionID(input string) (*ConnectionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := ConnectionId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.AutomationAccountName, err = id.PopSegment("automationAccounts"); err != nil {
		return nil, err
	}
	if resourceId.Name, err = id.PopSegment("connections"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
