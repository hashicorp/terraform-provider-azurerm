package activitylogs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EventData struct {
	Authorization        *SenderAuthorization `json:"authorization,omitempty"`
	Caller               *string              `json:"caller,omitempty"`
	Category             *LocalizableString   `json:"category,omitempty"`
	Claims               *map[string]string   `json:"claims,omitempty"`
	CorrelationId        *string              `json:"correlationId,omitempty"`
	Description          *string              `json:"description,omitempty"`
	EventDataId          *string              `json:"eventDataId,omitempty"`
	EventName            *LocalizableString   `json:"eventName,omitempty"`
	EventTimestamp       *string              `json:"eventTimestamp,omitempty"`
	HTTPRequest          *HTTPRequestInfo     `json:"httpRequest,omitempty"`
	Id                   *string              `json:"id,omitempty"`
	Level                *EventLevel          `json:"level,omitempty"`
	OperationId          *string              `json:"operationId,omitempty"`
	OperationName        *LocalizableString   `json:"operationName,omitempty"`
	Properties           *map[string]string   `json:"properties,omitempty"`
	ResourceGroupName    *string              `json:"resourceGroupName,omitempty"`
	ResourceId           *string              `json:"resourceId,omitempty"`
	ResourceProviderName *LocalizableString   `json:"resourceProviderName,omitempty"`
	ResourceType         *LocalizableString   `json:"resourceType,omitempty"`
	Status               *LocalizableString   `json:"status,omitempty"`
	SubStatus            *LocalizableString   `json:"subStatus,omitempty"`
	SubmissionTimestamp  *string              `json:"submissionTimestamp,omitempty"`
	SubscriptionId       *string              `json:"subscriptionId,omitempty"`
	TenantId             *string              `json:"tenantId,omitempty"`
}

func (o *EventData) GetEventTimestampAsTime() (*time.Time, error) {
	if o.EventTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EventTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *EventData) SetEventTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EventTimestamp = &formatted
}

func (o *EventData) GetSubmissionTimestampAsTime() (*time.Time, error) {
	if o.SubmissionTimestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SubmissionTimestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *EventData) SetSubmissionTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SubmissionTimestamp = &formatted
}
