// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_alert_processing_rule_action_group

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/alertsmanagement/2021-08-08/alertprocessingrules"
)

func flattenAlertProcessingRuleAddActionGroupId(input []alertprocessingrules.Action) ([]string, error) {
	if len(input) != 1 {
		return make([]string, 0), fmt.Errorf("read add_action_group_ids, the action should contains 1 element, but get %d element", len(input))
	}

	if addActionGroups, ok := input[0].(alertprocessingrules.AddActionGroups); ok {
		return addActionGroups.ActionGroupIds, nil
	}

	return make([]string, 0), fmt.Errorf("read add_action_group_ids, get unsupported action type %v", input[0])
}
