package reports

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReportRecordContract struct {
	ApiId            *string  `json:"apiId,omitempty"`
	ApiRegion        *string  `json:"apiRegion,omitempty"`
	ApiTimeAvg       *float64 `json:"apiTimeAvg,omitempty"`
	ApiTimeMax       *float64 `json:"apiTimeMax,omitempty"`
	ApiTimeMin       *float64 `json:"apiTimeMin,omitempty"`
	Bandwidth        *int64   `json:"bandwidth,omitempty"`
	CacheHitCount    *int64   `json:"cacheHitCount,omitempty"`
	CacheMissCount   *int64   `json:"cacheMissCount,omitempty"`
	CallCountBlocked *int64   `json:"callCountBlocked,omitempty"`
	CallCountFailed  *int64   `json:"callCountFailed,omitempty"`
	CallCountOther   *int64   `json:"callCountOther,omitempty"`
	CallCountSuccess *int64   `json:"callCountSuccess,omitempty"`
	CallCountTotal   *int64   `json:"callCountTotal,omitempty"`
	Country          *string  `json:"country,omitempty"`
	Interval         *string  `json:"interval,omitempty"`
	Name             *string  `json:"name,omitempty"`
	OperationId      *string  `json:"operationId,omitempty"`
	ProductId        *string  `json:"productId,omitempty"`
	Region           *string  `json:"region,omitempty"`
	ServiceTimeAvg   *float64 `json:"serviceTimeAvg,omitempty"`
	ServiceTimeMax   *float64 `json:"serviceTimeMax,omitempty"`
	ServiceTimeMin   *float64 `json:"serviceTimeMin,omitempty"`
	SubscriptionId   *string  `json:"subscriptionId,omitempty"`
	Timestamp        *string  `json:"timestamp,omitempty"`
	UserId           *string  `json:"userId,omitempty"`
	Zip              *string  `json:"zip,omitempty"`
}

func (o *ReportRecordContract) GetTimestampAsTime() (*time.Time, error) {
	if o.Timestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Timestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ReportRecordContract) SetTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Timestamp = &formatted
}
