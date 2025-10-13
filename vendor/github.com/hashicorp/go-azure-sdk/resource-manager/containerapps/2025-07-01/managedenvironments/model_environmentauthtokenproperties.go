package managedenvironments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type EnvironmentAuthTokenProperties struct {
	Expires *string `json:"expires,omitempty"`
	Token   *string `json:"token,omitempty"`
}

func (o *EnvironmentAuthTokenProperties) GetExpiresAsTime() (*time.Time, error) {
	if o.Expires == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Expires, "2006-01-02T15:04:05Z07:00")
}

func (o *EnvironmentAuthTokenProperties) SetExpiresAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Expires = &formatted
}
