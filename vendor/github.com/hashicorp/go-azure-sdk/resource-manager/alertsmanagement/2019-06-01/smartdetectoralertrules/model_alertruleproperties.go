package smartdetectoralertrules

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AlertRuleProperties struct {
	ActionGroups ActionGroupsInformation `json:"actionGroups"`
	Description  *string                 `json:"description,omitempty"`
	Detector     Detector                `json:"detector"`
	Frequency    string                  `json:"frequency"`
	Scope        []string                `json:"scope"`
	Severity     Severity                `json:"severity"`
	State        AlertRuleState          `json:"state"`
	Throttling   *ThrottlingInformation  `json:"throttling,omitempty"`
}
