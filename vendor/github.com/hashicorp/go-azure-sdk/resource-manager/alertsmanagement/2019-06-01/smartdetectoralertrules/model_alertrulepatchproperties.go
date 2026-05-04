package smartdetectoralertrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRulePatchProperties struct {
	ActionGroups *ActionGroupsInformation `json:"actionGroups,omitempty"`
	Description  *string                  `json:"description,omitempty"`
	Frequency    *string                  `json:"frequency,omitempty"`
	Severity     *Severity                `json:"severity,omitempty"`
	State        *AlertRuleState          `json:"state,omitempty"`
	Throttling   *ThrottlingInformation   `json:"throttling,omitempty"`
}
