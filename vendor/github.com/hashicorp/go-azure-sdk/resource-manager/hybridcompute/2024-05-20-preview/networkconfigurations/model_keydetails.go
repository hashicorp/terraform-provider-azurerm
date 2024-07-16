package networkconfigurations

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type KeyDetails struct {
	NotAfter   *string `json:"notAfter,omitempty"`
	PublicKey  *string `json:"publicKey,omitempty"`
	RenewAfter *string `json:"renewAfter,omitempty"`
}

func (o *KeyDetails) GetNotAfterAsTime() (*time.Time, error) {
	if o.NotAfter == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.NotAfter, "2006-01-02T15:04:05Z07:00")
}

func (o *KeyDetails) SetNotAfterAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.NotAfter = &formatted
}

func (o *KeyDetails) GetRenewAfterAsTime() (*time.Time, error) {
	if o.RenewAfter == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RenewAfter, "2006-01-02T15:04:05Z07:00")
}

func (o *KeyDetails) SetRenewAfterAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RenewAfter = &formatted
}
