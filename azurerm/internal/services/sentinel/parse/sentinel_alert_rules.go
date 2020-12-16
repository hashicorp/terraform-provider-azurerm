package parse

import (
	"fmt"
	"regexp"

	loganalyticsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/loganalytics/parse"
)

type SentinelAlertRuleId struct {
	SubscriptionId string
	ResourceGroup  string
	Workspace      string
	Name           string
}

func SentinelAlertRuleID(input string) (*SentinelAlertRuleId, error) {
	// Example ID: /subscriptions/<sub1>/resourceGroups/<grp1>/providers/Microsoft.OperationalInsights/workspaces/<workspace1>/providers/Microsoft.SecurityInsights/alertRules/<rule1>
	groups := regexp.MustCompile(`^(.+)/providers/Microsoft\.SecurityInsights/alertRules/(.+)$`).FindStringSubmatch(input)
	if len(groups) != 3 {
		return nil, fmt.Errorf("failed to parse Sentinel Alert Rule ID: %q", input)
	}

	workspace, name := groups[1], groups[2]

	workspaceId, err := loganalyticsParse.LogAnalyticsWorkspaceID(workspace)
	if err != nil {
		return nil, fmt.Errorf("parsing workspace part of Sentinel Alert Rule ID %q: %+v", input, err)
	}
	return &SentinelAlertRuleId{
		SubscriptionId: workspaceId.SubscriptionId,
		ResourceGroup:  workspaceId.ResourceGroup,
		Workspace:      workspaceId.WorkspaceName,
		Name:           name,
	}, nil
}
