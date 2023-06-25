package labs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LabAnnouncementProperties struct {
	Enabled           *EnableStatus `json:"enabled,omitempty"`
	ExpirationDate    *string       `json:"expirationDate,omitempty"`
	Expired           *bool         `json:"expired,omitempty"`
	Markdown          *string       `json:"markdown,omitempty"`
	ProvisioningState *string       `json:"provisioningState,omitempty"`
	Title             *string       `json:"title,omitempty"`
	UniqueIdentifier  *string       `json:"uniqueIdentifier,omitempty"`
}

func (o *LabAnnouncementProperties) GetExpirationDateAsTime() (*time.Time, error) {
	if o.ExpirationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpirationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *LabAnnouncementProperties) SetExpirationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpirationDate = &formatted
}
