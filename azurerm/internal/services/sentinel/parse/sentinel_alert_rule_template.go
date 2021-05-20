package parse

// NOTE: this file is generated via 'go:generate' - manual changes will be overwritten

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SentinelAlertRuleTemplateId struct {
	SubscriptionId        string
	ResourceGroup         string
	WorkspaceName         string
	AlertRuleTemplateName string
}

func NewSentinelAlertRuleTemplateID(subscriptionId, resourceGroup, workspaceName, alertRuleTemplateName string) SentinelAlertRuleTemplateId {
	return SentinelAlertRuleTemplateId{
		SubscriptionId:        subscriptionId,
		ResourceGroup:         resourceGroup,
		WorkspaceName:         workspaceName,
		AlertRuleTemplateName: alertRuleTemplateName,
	}
}

func (id SentinelAlertRuleTemplateId) String() string {
	segments := []string{
		fmt.Sprintf("Alert Rule Template Name %q", id.AlertRuleTemplateName),
		fmt.Sprintf("Workspace Name %q", id.WorkspaceName),
		fmt.Sprintf("Resource Group %q", id.ResourceGroup),
	}
	segmentsStr := strings.Join(segments, " / ")
	return fmt.Sprintf("%s: (%s)", "Sentinel Alert Rule Template", segmentsStr)
}

func (id SentinelAlertRuleTemplateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/AlertRuleTemplates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.AlertRuleTemplateName)
}

// SentinelAlertRuleTemplateID parses a SentinelAlertRuleTemplate ID into an SentinelAlertRuleTemplateId struct
func SentinelAlertRuleTemplateID(input string) (*SentinelAlertRuleTemplateId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, err
	}

	resourceId := SentinelAlertRuleTemplateId{
		SubscriptionId: id.SubscriptionID,
		ResourceGroup:  id.ResourceGroup,
	}

	if resourceId.SubscriptionId == "" {
		return nil, fmt.Errorf("ID was missing the 'subscriptions' element")
	}

	if resourceId.ResourceGroup == "" {
		return nil, fmt.Errorf("ID was missing the 'resourceGroups' element")
	}

	if resourceId.WorkspaceName, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}
	if resourceId.AlertRuleTemplateName, err = id.PopSegment("AlertRuleTemplates"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &resourceId, nil
}
