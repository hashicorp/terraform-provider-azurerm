package staticsites

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StaticSiteUserInvitationResponseResourceProperties struct {
	ExpiresOn     *string `json:"expiresOn,omitempty"`
	InvitationURL *string `json:"invitationUrl,omitempty"`
}

func (o *StaticSiteUserInvitationResponseResourceProperties) GetExpiresOnAsTime() (*time.Time, error) {
	if o.ExpiresOn == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpiresOn, "2006-01-02T15:04:05Z07:00")
}

func (o *StaticSiteUserInvitationResponseResourceProperties) SetExpiresOnAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpiresOn = &formatted
}
