package parse

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

type SentinelAlertRuleId struct {
	ResourceGroup string
	Workspace     string
	Name          string
}

func SentinelAlertRuleID(input string) (*SentinelAlertRuleId, error) {
	// Example ID: /subscriptions/<sub1>/resourceGroups/<grp1>/providers/Microsoft.OperationalInsights/workspaces/<workspace1>/providers/Microsoft.SecurityInsights/alertRules/<rule1>
	groups := regexp.MustCompile(`^(.+)/providers/Microsoft.SecurityInsights/alertRules/(.+)$`).FindStringSubmatch(input)
	if len(groups) != 3 {
		return nil, fmt.Errorf("faield to parse Sentinel Alert Rule ID: %q", input)
	}

	scope, name := groups[1], groups[2]

	ruleId := SentinelAlertRuleId{Name: name}

	id, err := azure.ParseAzureResourceID(scope)
	if err != nil {
		return nil, fmt.Errorf("parsing scope of Sentinel Alert Rule ID %q: %+v", input, err)
	}

	ruleId.ResourceGroup = id.ResourceGroup

	if ruleId.Workspace, err = id.PopSegment("workspaces"); err != nil {
		return nil, err
	}

	if err := id.ValidateNoEmptySegments(input); err != nil {
		return nil, err
	}

	return &ruleId, nil
}
