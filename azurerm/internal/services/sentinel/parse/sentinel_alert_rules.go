package parse

import (
	"fmt"
	"regexp"

	loganalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
)

type SentinelAlertRuleId struct {
	Subscription  string
	ResourceGroup string
	Workspace     string
	Name          string
}

func (rid SentinelAlertRuleId) String() string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/providers/Microsoft.SecurityInsights/alertRules/%s",
		rid.Subscription, rid.ResourceGroup, rid.Workspace, rid.Name)
}

func SentinelAlertRuleID(input string) (*SentinelAlertRuleId, error) {
	// Example ID: /subscriptions/<sub1>/resourceGroups/<grp1>/providers/Microsoft.OperationalInsights/workspaces/<workspace1>/providers/Microsoft.SecurityInsights/alertRules/<rule1>
	groups := regexp.MustCompile(`^(.+)/providers/Microsoft\.SecurityInsights/alertRules/([^/]+)$`).FindStringSubmatch(input)
	if len(groups) != 3 {
		return nil, fmt.Errorf("failed to parse Sentinel Alert Rule ID: %q", input)
	}

	workspace, name := groups[1], groups[2]

	workspaceId, err := loganalyticsParse.LogAnalyticsWorkspaceID(workspace)
	if err != nil {
		return nil, fmt.Errorf("parsing workspace part of Sentinel Alert Rule ID %q: %+v", input, err)
	}
	return &SentinelAlertRuleId{
		Subscription:  workspaceId.Subscription,
		ResourceGroup: workspaceId.ResourceGroup,
		Workspace:     workspaceId.Name,
		Name:          name,
	}, nil
}

type SentinelAlertRuleActionId struct {
	Subscription  string
	ResourceGroup string
	Workspace     string
	Rule          string
	Name          string
}

func (aid SentinelAlertRuleActionId) FormatSentinelAlertRuleId() SentinelAlertRuleId {
	return SentinelAlertRuleId{
		Subscription:  aid.Subscription,
		ResourceGroup: aid.ResourceGroup,
		Workspace:     aid.Workspace,
		Name:          aid.Rule,
	}
}

func SentinelAlertRuleActionID(input string) (*SentinelAlertRuleActionId, error) {
	// Example ID: /subscriptions/<sub1>/resourceGroups/<grp1>/providers/Microsoft.OperationalInsights/workspaces/<workspace1>/providers/Microsoft.SecurityInsights/alertRules/<rule1>/actions/<action1>
	groups := regexp.MustCompile(`^(.+)/actions/([^/]+)$`).FindStringSubmatch(input)
	if len(groups) != 3 {
		return nil, fmt.Errorf("faield to parse Sentinel Alert Rule Action ID: %q", input)
	}

	rule, name := groups[1], groups[2]

	ruleID, err := SentinelAlertRuleID(rule)
	if err != nil {
		return nil, fmt.Errorf("parsing Alert Rule part of Sentinel Alert Rule Action ID %q: %+v", input, err)
	}
	return &SentinelAlertRuleActionId{
		Subscription:  ruleID.Subscription,
		ResourceGroup: ruleID.ResourceGroup,
		Workspace:     ruleID.Workspace,
		Rule:          ruleID.Name,
		Name:          name,
	}, nil
}
