package alertruletemplates

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NrtAlertRuleTemplateProperties struct {
	AlertDetailsOverride             *AlertDetailsOverride          `json:"alertDetailsOverride,omitempty"`
	AlertRulesCreatedByTemplateCount int64                          `json:"alertRulesCreatedByTemplateCount"`
	CreatedDateUTC                   *string                        `json:"createdDateUTC,omitempty"`
	CustomDetails                    *map[string]string             `json:"customDetails,omitempty"`
	Description                      string                         `json:"description"`
	DisplayName                      string                         `json:"displayName"`
	EntityMappings                   *[]EntityMapping               `json:"entityMappings,omitempty"`
	LastUpdatedDateUTC               *string                        `json:"lastUpdatedDateUTC,omitempty"`
	Query                            string                         `json:"query"`
	RequiredDataConnectors           *[]AlertRuleTemplateDataSource `json:"requiredDataConnectors,omitempty"`
	Severity                         AlertSeverity                  `json:"severity"`
	Status                           TemplateStatus                 `json:"status"`
	Tactics                          *[]AttackTactic                `json:"tactics,omitempty"`
	Version                          string                         `json:"version"`
}

func (o *NrtAlertRuleTemplateProperties) GetCreatedDateUTCAsTime() (*time.Time, error) {
	if o.CreatedDateUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDateUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *NrtAlertRuleTemplateProperties) SetCreatedDateUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDateUTC = &formatted
}

func (o *NrtAlertRuleTemplateProperties) GetLastUpdatedDateUTCAsTime() (*time.Time, error) {
	if o.LastUpdatedDateUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedDateUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *NrtAlertRuleTemplateProperties) SetLastUpdatedDateUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedDateUTC = &formatted
}
