package alertruletemplates

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NrtAlertRuleTemplateProperties struct {
	AlertDetailsOverride     *AlertDetailsOverride    `json:"alertDetailsOverride,omitempty"`
	CustomDetails            *map[string]string       `json:"customDetails,omitempty"`
	Description              *string                  `json:"description,omitempty"`
	DisplayName              *string                  `json:"displayName,omitempty"`
	EntityMappings           *[]EntityMapping         `json:"entityMappings,omitempty"`
	EventGroupingSettings    *EventGroupingSettings   `json:"eventGroupingSettings,omitempty"`
	Query                    string                   `json:"query"`
	SentinelEntitiesMappings *[]SentinelEntityMapping `json:"sentinelEntitiesMappings,omitempty"`
	Severity                 AlertSeverity            `json:"severity"`
	Tactics                  *[]AttackTactic          `json:"tactics,omitempty"`
	Techniques               *[]string                `json:"techniques,omitempty"`
	Version                  string                   `json:"version"`
}
