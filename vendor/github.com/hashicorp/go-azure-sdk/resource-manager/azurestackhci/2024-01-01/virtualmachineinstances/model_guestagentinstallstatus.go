package virtualmachineinstances

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GuestAgentInstallStatus struct {
	AgentVersion     *string        `json:"agentVersion,omitempty"`
	ErrorDetails     *[]ErrorDetail `json:"errorDetails,omitempty"`
	LastStatusChange *string        `json:"lastStatusChange,omitempty"`
	Status           *StatusTypes   `json:"status,omitempty"`
	VMUuid           *string        `json:"vmUuid,omitempty"`
}

func (o *GuestAgentInstallStatus) GetLastStatusChangeAsTime() (*time.Time, error) {
	if o.LastStatusChange == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastStatusChange, "2006-01-02T15:04:05Z07:00")
}

func (o *GuestAgentInstallStatus) SetLastStatusChangeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStatusChange = &formatted
}
