package networkstatus

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectivityStatusContract struct {
	Error            *string                `json:"error,omitempty"`
	IsOptional       bool                   `json:"isOptional"`
	LastStatusChange string                 `json:"lastStatusChange"`
	LastUpdated      string                 `json:"lastUpdated"`
	Name             string                 `json:"name"`
	ResourceType     string                 `json:"resourceType"`
	Status           ConnectivityStatusType `json:"status"`
}

func (o *ConnectivityStatusContract) GetLastStatusChangeAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.LastStatusChange, "2006-01-02T15:04:05Z07:00")
}

func (o *ConnectivityStatusContract) SetLastStatusChangeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStatusChange = formatted
}

func (o *ConnectivityStatusContract) GetLastUpdatedAsTime() (*time.Time, error) {
	return dates.ParseAsFormat(&o.LastUpdated, "2006-01-02T15:04:05Z07:00")
}

func (o *ConnectivityStatusContract) SetLastUpdatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdated = formatted
}
