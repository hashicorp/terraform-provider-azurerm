package resourceproviders

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SourceControlProperties struct {
	ExpirationTime *string `json:"expirationTime,omitempty"`
	RefreshToken   *string `json:"refreshToken,omitempty"`
	Token          *string `json:"token,omitempty"`
	TokenSecret    *string `json:"tokenSecret,omitempty"`
}

func (o *SourceControlProperties) GetExpirationTimeAsTime() (*time.Time, error) {
	if o.ExpirationTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SourceControlProperties) SetExpirationTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationTime = &formatted
}
