package alertrules

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MLBehaviorAnalyticsAlertRuleProperties struct {
	AlertRuleTemplateName string          `json:"alertRuleTemplateName"`
	Description           *string         `json:"description,omitempty"`
	DisplayName           *string         `json:"displayName,omitempty"`
	Enabled               bool            `json:"enabled"`
	LastModifiedUtc       *string         `json:"lastModifiedUtc,omitempty"`
	Severity              *AlertSeverity  `json:"severity,omitempty"`
	Tactics               *[]AttackTactic `json:"tactics,omitempty"`
	Techniques            *[]string       `json:"techniques,omitempty"`
}

func (o *MLBehaviorAnalyticsAlertRuleProperties) GetLastModifiedUtcAsTime() (*time.Time, error) {
	if o.LastModifiedUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *MLBehaviorAnalyticsAlertRuleProperties) SetLastModifiedUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedUtc = &formatted
}
