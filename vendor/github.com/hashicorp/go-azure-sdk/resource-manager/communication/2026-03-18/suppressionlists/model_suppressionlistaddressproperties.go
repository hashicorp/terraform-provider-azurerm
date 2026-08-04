package suppressionlists

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SuppressionListAddressProperties struct {
	DataLocation *string `json:"dataLocation,omitempty"`
	Email        string  `json:"email"`
	FirstName    *string `json:"firstName,omitempty"`
	LastModified *string `json:"lastModified,omitempty"`
	LastName     *string `json:"lastName,omitempty"`
	Notes        *string `json:"notes,omitempty"`
}

func (o *SuppressionListAddressProperties) GetLastModifiedAsTime() (*time.Time, error) {
	if o.LastModified == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModified, "2006-01-02T15:04:05Z07:00")
}

func (o *SuppressionListAddressProperties) SetLastModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModified = &formatted
}
