package managedclusters

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CommandResultProperties struct {
	ExitCode          *int64  `json:"exitCode,omitempty"`
	FinishedAt        *string `json:"finishedAt,omitempty"`
	Logs              *string `json:"logs,omitempty"`
	ProvisioningState *string `json:"provisioningState,omitempty"`
	Reason            *string `json:"reason,omitempty"`
	StartedAt         *string `json:"startedAt,omitempty"`
}

func (o *CommandResultProperties) GetFinishedAtAsTime() (*time.Time, error) {
	if o.FinishedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.FinishedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *CommandResultProperties) SetFinishedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.FinishedAt = &formatted
}

func (o *CommandResultProperties) GetStartedAtAsTime() (*time.Time, error) {
	if o.StartedAt == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartedAt, "2006-01-02T15:04:05Z07:00")
}

func (o *CommandResultProperties) SetStartedAtAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartedAt = &formatted
}
