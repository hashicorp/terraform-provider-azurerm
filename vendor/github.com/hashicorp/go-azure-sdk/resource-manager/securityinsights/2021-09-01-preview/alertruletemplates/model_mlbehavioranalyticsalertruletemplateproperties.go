package alertruletemplates

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MLBehaviorAnalyticsAlertRuleTemplateProperties struct {
	AlertRulesCreatedByTemplateCount int64                          `json:"alertRulesCreatedByTemplateCount"`
	CreatedDateUTC                   *string                        `json:"createdDateUTC,omitempty"`
	Description                      string                         `json:"description"`
	DisplayName                      string                         `json:"displayName"`
	LastUpdatedDateUTC               *string                        `json:"lastUpdatedDateUTC,omitempty"`
	RequiredDataConnectors           *[]AlertRuleTemplateDataSource `json:"requiredDataConnectors,omitempty"`
	Severity                         AlertSeverity                  `json:"severity"`
	Status                           TemplateStatus                 `json:"status"`
	Tactics                          *[]AttackTactic                `json:"tactics,omitempty"`
}

func (o *MLBehaviorAnalyticsAlertRuleTemplateProperties) GetCreatedDateUTCAsTime() (*time.Time, error) {
	if o.CreatedDateUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDateUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *MLBehaviorAnalyticsAlertRuleTemplateProperties) SetCreatedDateUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDateUTC = &formatted
}

func (o *MLBehaviorAnalyticsAlertRuleTemplateProperties) GetLastUpdatedDateUTCAsTime() (*time.Time, error) {
	if o.LastUpdatedDateUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedDateUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *MLBehaviorAnalyticsAlertRuleTemplateProperties) SetLastUpdatedDateUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedDateUTC = &formatted
}
