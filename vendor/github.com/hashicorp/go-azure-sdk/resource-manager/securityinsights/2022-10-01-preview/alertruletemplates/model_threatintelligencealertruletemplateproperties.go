package alertruletemplates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ThreatIntelligenceAlertRuleTemplateProperties struct {
	Description *string         `json:"description,omitempty"`
	DisplayName *string         `json:"displayName,omitempty"`
	Severity    AlertSeverity   `json:"severity"`
	Tactics     *[]AttackTactic `json:"tactics,omitempty"`
	Techniques  *[]string       `json:"techniques,omitempty"`
}
