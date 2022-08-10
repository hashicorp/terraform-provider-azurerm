package topics

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SBTopicProperties struct {
	AccessedAt                          *string              `json:"accessedAt,omitempty"`
	AutoDeleteOnIdle                    *string              `json:"autoDeleteOnIdle,omitempty"`
	CountDetails                        *MessageCountDetails `json:"countDetails,omitempty"`
	CreatedAt                           *string              `json:"createdAt,omitempty"`
	DefaultMessageTimeToLive            *string              `json:"defaultMessageTimeToLive,omitempty"`
	DuplicateDetectionHistoryTimeWindow *string              `json:"duplicateDetectionHistoryTimeWindow,omitempty"`
	EnableBatchedOperations             *bool                `json:"enableBatchedOperations,omitempty"`
	EnableExpress                       *bool                `json:"enableExpress,omitempty"`
	EnablePartitioning                  *bool                `json:"enablePartitioning,omitempty"`
	MaxMessageSizeInKilobytes           *int64               `json:"maxMessageSizeInKilobytes,omitempty"`
	MaxSizeInMegabytes                  *int64               `json:"maxSizeInMegabytes,omitempty"`
	RequiresDuplicateDetection          *bool                `json:"requiresDuplicateDetection,omitempty"`
	SizeInBytes                         *int64               `json:"sizeInBytes,omitempty"`
	Status                              *EntityStatus        `json:"status,omitempty"`
	SubscriptionCount                   *int64               `json:"subscriptionCount,omitempty"`
	SupportOrdering                     *bool                `json:"supportOrdering,omitempty"`
	UpdatedAt                           *string              `json:"updatedAt,omitempty"`
}

func (o *SBTopicProperties) GetAccessedAtAsTime() (*time.Time, error) {
	if o.AccessedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.AccessedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *SBTopicProperties) SetAccessedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.AccessedAt = &formatted
}

func (o *SBTopicProperties) GetCreatedAtAsTime() (*time.Time, error) {
	if o.CreatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *SBTopicProperties) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o *SBTopicProperties) GetUpdatedAtAsTime() (*time.Time, error) {
	if o.UpdatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *SBTopicProperties) SetUpdatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedAt = &formatted
}
