package certificateprofiles

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type Revocation struct {
	EffectiveAt   *string           `json:"effectiveAt,omitempty"`
	FailureReason *string           `json:"failureReason,omitempty"`
	Reason        *string           `json:"reason,omitempty"`
	Remarks       *string           `json:"remarks,omitempty"`
	RequestedAt   *string           `json:"requestedAt,omitempty"`
	Status        *RevocationStatus `json:"status,omitempty"`
}

func (o *Revocation) GetEffectiveAtAsTime() (*time.Time, error) {
	if o.EffectiveAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EffectiveAt, "2006-01-02T15:04:05Z07:00")
}

func (o *Revocation) SetEffectiveAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EffectiveAt = &formatted
}

func (o *Revocation) GetRequestedAtAsTime() (*time.Time, error) {
	if o.RequestedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RequestedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *Revocation) SetRequestedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RequestedAt = &formatted
}
