package reports

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RequestReportRecordContract struct {
	ApiId               *string  `json:"apiId,omitempty"`
	ApiRegion           *string  `json:"apiRegion,omitempty"`
	ApiTime             *float64 `json:"apiTime,omitempty"`
	BackendResponseCode *string  `json:"backendResponseCode,omitempty"`
	Cache               *string  `json:"cache,omitempty"`
	IPAddress           *string  `json:"ipAddress,omitempty"`
	Method              *string  `json:"method,omitempty"`
	OperationId         *string  `json:"operationId,omitempty"`
	ProductId           *string  `json:"productId,omitempty"`
	RequestId           *string  `json:"requestId,omitempty"`
	RequestSize         *int64   `json:"requestSize,omitempty"`
	ResponseCode        *int64   `json:"responseCode,omitempty"`
	ResponseSize        *int64   `json:"responseSize,omitempty"`
	ServiceTime         *float64 `json:"serviceTime,omitempty"`
	SubscriptionId      *string  `json:"subscriptionId,omitempty"`
	Timestamp           *string  `json:"timestamp,omitempty"`
	Url                 *string  `json:"url,omitempty"`
	UserId              *string  `json:"userId,omitempty"`
}

func (o *RequestReportRecordContract) GetTimestampAsTime() (*time.Time, error) {
	if o.Timestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Timestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *RequestReportRecordContract) SetTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Timestamp = &formatted
}
