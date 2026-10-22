// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package monitor_alert_processing_rule_suppression

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/helpers"
)

type AlertProcessingRuleSuppressionModel struct {
	Name              string                                      `tfschema:"name"`
	ResourceGroupName string                                      `tfschema:"resource_group_name"`
	Scopes            []string                                    `tfschema:"scopes"`
	Description       string                                      `tfschema:"description"`
	Enabled           bool                                        `tfschema:"enabled"`
	Condition         []helpers.AlertProcessingRuleConditionModel `tfschema:"condition"`
	Schedule          []helpers.AlertProcessingRuleScheduleModel  `tfschema:"schedule"`
	Tags              map[string]string                           `tfschema:"tags"`
}

func (r Resource) ModelObject() any {
	return &AlertProcessingRuleSuppressionModel{}
}
