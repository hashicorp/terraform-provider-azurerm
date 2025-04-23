package consumergroups

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConsumerGroupProperties struct {
	CreatedAt    *string `json:"createdAt,omitempty"`
	UpdatedAt    *string `json:"updatedAt,omitempty"`
	UserMetadata *string `json:"userMetadata,omitempty"`
}

func (o *ConsumerGroupProperties) GetCreatedAtAsTime() (*time.Time, error) {
	if o.CreatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *ConsumerGroupProperties) SetCreatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedAt = &formatted
}

func (o *ConsumerGroupProperties) GetUpdatedAtAsTime() (*time.Time, error) {
	if o.UpdatedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.UpdatedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *ConsumerGroupProperties) SetUpdatedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.UpdatedAt = &formatted
}
