package managedinstancekeys

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedInstanceKeyProperties struct {
	AutoRotationEnabled *bool         `json:"autoRotationEnabled,omitempty"`
	CreationDate        *string       `json:"creationDate,omitempty"`
	ServerKeyType       ServerKeyType `json:"serverKeyType"`
	Thumbprint          *string       `json:"thumbprint,omitempty"`
	Uri                 *string       `json:"uri,omitempty"`
}

func (o *ManagedInstanceKeyProperties) GetCreationDateAsTime() (*time.Time, error) {
	if o.CreationDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreationDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ManagedInstanceKeyProperties) SetCreationDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreationDate = &formatted
}
