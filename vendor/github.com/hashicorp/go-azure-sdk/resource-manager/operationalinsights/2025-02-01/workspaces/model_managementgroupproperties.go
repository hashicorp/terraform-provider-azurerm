package workspaces

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagementGroupProperties struct {
	Created      *string `json:"created,omitempty"`
	DataReceived *string `json:"dataReceived,omitempty"`
	Id           *string `json:"id,omitempty"`
	IsGateway    *bool   `json:"isGateway,omitempty"`
	Name         *string `json:"name,omitempty"`
	ServerCount  *int64  `json:"serverCount,omitempty"`
	Sku          *string `json:"sku,omitempty"`
	Version      *string `json:"version,omitempty"`
}

func (o *ManagementGroupProperties) GetCreatedAsTime() (*time.Time, error) {
	if o.Created == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.Created, "2006-01-02T15:04:05Z07:00")
}

func (o *ManagementGroupProperties) SetCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.Created = &formatted
}

func (o *ManagementGroupProperties) GetDataReceivedAsTime() (*time.Time, error) {
	if o.DataReceived == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DataReceived, "2006-01-02T15:04:05Z07:00")
}

func (o *ManagementGroupProperties) SetDataReceivedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DataReceived = &formatted
}
