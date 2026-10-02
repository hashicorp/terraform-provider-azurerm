package hostpool

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegistrationTokenMinimal struct {
	ExpirationTime *string `json:"expirationTime,omitempty"`
	Token          *string `json:"token,omitempty"`
}

func (o *RegistrationTokenMinimal) GetExpirationTimeAsTime() (*time.Time, error) {
	if o.ExpirationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RegistrationTokenMinimal) SetExpirationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationTime = &formatted
}
