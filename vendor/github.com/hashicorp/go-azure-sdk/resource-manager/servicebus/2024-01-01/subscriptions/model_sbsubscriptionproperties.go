package subscriptions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SBSubscriptionProperties struct {
	AccessedAt                                *string                   `json:"accessedAt,omitempty"`
	AutoDeleteOnIdle                          *string                   `json:"autoDeleteOnIdle,omitempty"`
	ClientAffineProperties                    *SBClientAffineProperties `json:"clientAffineProperties,omitempty"`
	CountDetails                              *MessageCountDetails      `json:"countDetails,omitempty"`
	CreatedAt                                 *string                   `json:"createdAt,omitempty"`
	DeadLetteringOnFilterEvaluationExceptions *bool                     `json:"deadLetteringOnFilterEvaluationExceptions,omitempty"`
	DeadLetteringOnMessageExpiration          *bool                     `json:"deadLetteringOnMessageExpiration,omitempty"`
	DefaultMessageTimeToLive                  *string                   `json:"defaultMessageTimeToLive,omitempty"`
	DuplicateDetectionHistoryTimeWindow       *string                   `json:"duplicateDetectionHistoryTimeWindow,omitempty"`
	EnableBatchedOperations                   *bool                     `json:"enableBatchedOperations,omitempty"`
	ForwardDeadLetteredMessagesTo             *string                   `json:"forwardDeadLetteredMessagesTo,omitempty"`
	ForwardTo                                 *string                   `json:"forwardTo,omitempty"`
	IsClientAffine                            *bool                     `json:"isClientAffine,omitempty"`
	LockDuration                              *string                   `json:"lockDuration,omitempty"`
	MaxDeliveryCount                          *int64                    `json:"maxDeliveryCount,omitempty"`
	MessageCount                              *int64                    `json:"messageCount,omitempty"`
	RequiresSession                           *bool                     `json:"requiresSession,omitempty"`
	Status                                    *EntityStatus             `json:"status,omitempty"`
	UpdatedAt                                 *string                   `json:"updatedAt,omitempty"`
}

func (o *SBSubscriptionProperties) GetAccessedAtAsTime() (*time.Time, error) {
	if o.AccessedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.AccessedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *SBSubscriptionProperties) SetAccessedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.AccessedAt = &formatted
}

func (o *SBSubscriptionProperties) GetCreatedAtAsTime() (*time.Time, error) {
	if o.CreatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *SBSubscriptionProperties) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o *SBSubscriptionProperties) GetUpdatedAtAsTime() (*time.Time, error) {
	if o.UpdatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *SBSubscriptionProperties) SetUpdatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedAt = &formatted
}
