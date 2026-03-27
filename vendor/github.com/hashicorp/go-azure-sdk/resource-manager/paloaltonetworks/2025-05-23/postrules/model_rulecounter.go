package postrules

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RuleCounter struct {
	AppSeen              *AppSeenData `json:"appSeen,omitempty"`
	FirewallName         *string      `json:"firewallName,omitempty"`
	HitCount             *int64       `json:"hitCount,omitempty"`
	LastUpdatedTimestamp *string      `json:"lastUpdatedTimestamp,omitempty"`
	Priority             string       `json:"priority"`
	RequestTimestamp     *string      `json:"requestTimestamp,omitempty"`
	RuleListName         *string      `json:"ruleListName,omitempty"`
	RuleName             string       `json:"ruleName"`
	RuleStackName        *string      `json:"ruleStackName,omitempty"`
	Timestamp            *string      `json:"timestamp,omitempty"`
}

func (o *RuleCounter) GetLastUpdatedTimestampAsTime() (*time.Time, error) {
	if o.LastUpdatedTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *RuleCounter) SetLastUpdatedTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTimestamp = &formatted
}

func (o *RuleCounter) GetRequestTimestampAsTime() (*time.Time, error) {
	if o.RequestTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RequestTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *RuleCounter) SetRequestTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RequestTimestamp = &formatted
}

func (o *RuleCounter) GetTimestampAsTime() (*time.Time, error) {
	if o.Timestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Timestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *RuleCounter) SetTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Timestamp = &formatted
}
