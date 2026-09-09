package appserviceplans

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServerFarmRdpDetails struct {
	RdpPassword       *string `json:"rdpPassword,omitempty"`
	RdpPasswordExpiry *string `json:"rdpPasswordExpiry,omitempty"`
}

func (o *ServerFarmRdpDetails) GetRdpPasswordExpiryAsTime() (*time.Time, error) {
	if o.RdpPasswordExpiry == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RdpPasswordExpiry, "2006-01-02T15:04:05Z07:00")
}

func (o *ServerFarmRdpDetails) SetRdpPasswordExpiryAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RdpPasswordExpiry = &formatted
}
