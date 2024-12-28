package tokens

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TokenPassword struct {
	CreationTime *string            `json:"creationTime,omitempty"`
	Expiry       *string            `json:"expiry,omitempty"`
	Name         *TokenPasswordName `json:"name,omitempty"`
	Value        *string            `json:"value,omitempty"`
}

func (o *TokenPassword) GetCreationTimeAsTime() (*time.Time, error) {
	if o.CreationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *TokenPassword) SetCreationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationTime = &formatted
}

func (o *TokenPassword) GetExpiryAsTime() (*time.Time, error) {
	if o.Expiry == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Expiry, "2006-01-02T15:04:05Z07:00")
}

func (o *TokenPassword) SetExpiryAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Expiry = &formatted
}
