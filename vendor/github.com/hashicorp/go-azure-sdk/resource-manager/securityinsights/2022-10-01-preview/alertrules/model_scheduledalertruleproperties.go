package alertrules

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduledAlertRuleProperties struct {
	AlertDetailsOverride     *AlertDetailsOverride    `json:"alertDetailsOverride,omitempty"`
	AlertRuleTemplateName    *string                  `json:"alertRuleTemplateName,omitempty"`
	CustomDetails            *map[string]string       `json:"customDetails,omitempty"`
	Description              *string                  `json:"description,omitempty"`
	DisplayName              string                   `json:"displayName"`
	Enabled                  bool                     `json:"enabled"`
	EntityMappings           *[]EntityMapping         `json:"entityMappings,omitempty"`
	EventGroupingSettings    *EventGroupingSettings   `json:"eventGroupingSettings,omitempty"`
	IncidentConfiguration    *IncidentConfiguration   `json:"incidentConfiguration,omitempty"`
	LastModifiedUtc          *string                  `json:"lastModifiedUtc,omitempty"`
	Query                    *string                  `json:"query,omitempty"`
	QueryFrequency           *string                  `json:"queryFrequency,omitempty"`
	QueryPeriod              *string                  `json:"queryPeriod,omitempty"`
	SentinelEntitiesMappings *[]SentinelEntityMapping `json:"sentinelEntitiesMappings,omitempty"`
	Severity                 *AlertSeverity           `json:"severity,omitempty"`
	SuppressionDuration      string                   `json:"suppressionDuration"`
	SuppressionEnabled       bool                     `json:"suppressionEnabled"`
	Tactics                  *[]AttackTactic          `json:"tactics,omitempty"`
	Techniques               *[]string                `json:"techniques,omitempty"`
	TemplateVersion          *string                  `json:"templateVersion,omitempty"`
	TriggerOperator          *TriggerOperator         `json:"triggerOperator,omitempty"`
	TriggerThreshold         *int64                   `json:"triggerThreshold,omitempty"`
}

func (o *ScheduledAlertRuleProperties) GetLastModifiedUtcAsTime() (*time.Time, error) {
	if o.LastModifiedUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduledAlertRuleProperties) SetLastModifiedUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedUtc = &formatted
}
