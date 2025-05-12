package blobcontainers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UpdateHistoryProperty struct {
	AllowProtectedAppendWrites            *bool                         `json:"allowProtectedAppendWrites,omitempty"`
	AllowProtectedAppendWritesAll         *bool                         `json:"allowProtectedAppendWritesAll,omitempty"`
	ImmutabilityPeriodSinceCreationInDays *int64                        `json:"immutabilityPeriodSinceCreationInDays,omitempty"`
	ObjectIdentifier                      *string                       `json:"objectIdentifier,omitempty"`
	TenantId                              *string                       `json:"tenantId,omitempty"`
	Timestamp                             *string                       `json:"timestamp,omitempty"`
	Update                                *ImmutabilityPolicyUpdateType `json:"update,omitempty"`
	Upn                                   *string                       `json:"upn,omitempty"`
}

func (o *UpdateHistoryProperty) GetTimestampAsTime() (*time.Time, error) {
	if o.Timestamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Timestamp, "2006-01-02T15:04:05Z07:00")
}

func (o *UpdateHistoryProperty) SetTimestampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Timestamp = &formatted
}
