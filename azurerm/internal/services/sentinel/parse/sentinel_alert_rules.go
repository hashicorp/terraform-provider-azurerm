package parse

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SentinelAlertRuleId struct {
	ResourceGroup string
	Workspace     string
	Name          string
}

func NewSentinelAlertRuleID(resourceGroup, workspace, name string) SentinelAlertRuleId {
	return SentinelAlertRuleId{
		ResourceGroup: resourceGroup,
		Workspace:     workspace,
		Name:          name,
	}
}

func (rid SentinelAlertRuleId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/alertRules/%s",
		subscriptionId, rid.ResourceGroup, rid.Workspace, rid.Name)
}

func SentinelAlertRuleID(input string) (*SentinelAlertRuleId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Sentinel Alert Rule ID %q: %+v", input, err)
	}

	rule := SentinelAlertRuleId{
		ResourceGroup: id.ResourceGroup,
	}

	if rule.Workspace, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	if rule.Name, err = id.PopSegment("alertRules"); err != nil {
		return nil, err
	}

	return &rule, nil
}

type SentinelAlertRuleActionId struct {
	ResourceGroup string
	Workspace     string
	Rule          string
	Name          string
}

func NewSentinelAlertRuleActionID(resourceGroup, workspace, rule, name string) SentinelAlertRuleActionId {
	return SentinelAlertRuleActionId{
		ResourceGroup: resourceGroup,
		Workspace:     workspace,
		Rule:          rule,
		Name:          name,
	}
}

func (aid SentinelAlertRuleActionId) ID(subscriptionId string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/alertRules/%s/actions/%s",
		subscriptionId, aid.ResourceGroup, aid.Workspace, aid.Rule, aid.Name)
}

func (aid SentinelAlertRuleActionId) FormatSentinelAlertRuleId() SentinelAlertRuleId {
	return SentinelAlertRuleId{
		ResourceGroup: aid.ResourceGroup,
		Workspace:     aid.Workspace,
		Name:          aid.Rule,
	}
}

func SentinelAlertRuleActionID(input string) (*SentinelAlertRuleActionId, error) {
	id, err := azure.ParseAzureResourceID(input)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Unable to parse Sentinel Alert Rule Action ID %q: %+v", input, err)
	}

	action := SentinelAlertRuleActionId{
		ResourceGroup: id.ResourceGroup,
	}

	if action.Workspace, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	if action.Rule, err = id.PopSegment("alertRules"); err != nil {
		return nil, err
	}

	if action.Name, err = id.PopSegment("actions"); err != nil {
		return nil, err
	}

	return &action, nil
}
