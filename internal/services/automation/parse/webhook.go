package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

type WebhookId struct {
	SubscriptionId        string
	ResourceGroup         string
	AutomationAccountName string
	Name                  string
}

func NewWebhookID(subscriptionId, resourceGroup, automationAccountName, name string) WebhookId {
	return WebhookId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		AutomationAccountName: automationAccountName,
		Name:                  name,
	}
}

func (id WebhookId) String() string {
	segments := []string{
		fmt.Sprintf("Name %q", id.Name),
		fmt.Sprintf("Automation Account Name %q", id.AutomationAccountName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Webhook", segmentsStr)
}

func (id WebhookId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Automation/automationAccounts/%s/webhooks/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.AutomationAccountName, id.Name)
}

// WebhookID parses a Webhook ID into an WebhookId struct
func WebhookID(input string) (*WebhookId, error) {
	id, err := resourceids.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := WebhookId{
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
	if resourceId.Name, err = id.PopSegment("webhooks"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
