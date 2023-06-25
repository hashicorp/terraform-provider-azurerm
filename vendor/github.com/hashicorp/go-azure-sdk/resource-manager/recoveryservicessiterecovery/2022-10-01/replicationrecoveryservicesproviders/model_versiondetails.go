package replicationrecoveryservicesproviders

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VersionDetails struct {
	ExpiryDate *string             `json:"expiryDate,omitempty"`
	Status     *AgentVersionStatus `json:"status,omitempty"`
	Version    *string             `json:"version,omitempty"`
}

func (o *VersionDetails) GetExpiryDateAsTime() (*time.Time, error) {
	if o.ExpiryDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpiryDate, "2006-01-02T15:04:05Z07:00")
}

func (o *VersionDetails) SetExpiryDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpiryDate = &formatted
}
