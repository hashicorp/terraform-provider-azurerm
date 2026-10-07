package alertruletemplates

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ScheduledAlertRuleTemplateProperties struct {
	AlertDetailsOverride             *AlertDetailsOverride          `json:"alertDetailsOverride,omitempty"`
	AlertRulesCreatedByTemplateCount *int64                         `json:"alertRulesCreatedByTemplateCount,omitempty"`
	CreatedDateUTC                   *string                        `json:"createdDateUTC,omitempty"`
	CustomDetails                    *map[string]string             `json:"customDetails,omitempty"`
	Description                      *string                        `json:"description,omitempty"`
	DisplayName                      *string                        `json:"displayName,omitempty"`
	EntityMappings                   *[]EntityMapping               `json:"entityMappings,omitempty"`
	EventGroupingSettings            *EventGroupingSettings         `json:"eventGroupingSettings,omitempty"`
	LastUpdatedDateUTC               *string                        `json:"lastUpdatedDateUTC,omitempty"`
	Query                            *string                        `json:"query,omitempty"`
	QueryFrequency                   *string                        `json:"queryFrequency,omitempty"`
	QueryPeriod                      *string                        `json:"queryPeriod,omitempty"`
	RequiredDataConnectors           *[]AlertRuleTemplateDataSource `json:"requiredDataConnectors,omitempty"`
	SentinelEntitiesMappings         *[]SentinelEntityMapping       `json:"sentinelEntitiesMappings,omitempty"`
	Severity                         *AlertSeverity                 `json:"severity,omitempty"`
	Status                           *TemplateStatus                `json:"status,omitempty"`
	Tactics                          *[]AttackTactic                `json:"tactics,omitempty"`
	Techniques                       *[]string                      `json:"techniques,omitempty"`
	TriggerOperator                  *TriggerOperator               `json:"triggerOperator,omitempty"`
	TriggerThreshold                 *int64                         `json:"triggerThreshold,omitempty"`
	Version                          *string                        `json:"version,omitempty"`
}

func (o *ScheduledAlertRuleTemplateProperties) GetCreatedDateUTCAsTime() (*time.Time, error) {
	if o.CreatedDateUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDateUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduledAlertRuleTemplateProperties) SetCreatedDateUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDateUTC = &formatted
}

func (o *ScheduledAlertRuleTemplateProperties) GetLastUpdatedDateUTCAsTime() (*time.Time, error) {
	if o.LastUpdatedDateUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedDateUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *ScheduledAlertRuleTemplateProperties) SetLastUpdatedDateUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedDateUTC = &formatted
}
