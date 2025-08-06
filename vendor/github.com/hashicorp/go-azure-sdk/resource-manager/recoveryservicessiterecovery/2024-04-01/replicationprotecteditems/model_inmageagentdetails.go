package replicationprotecteditems

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageAgentDetails struct {
	AgentExpiryDate        *string `json:"agentExpiryDate,omitempty"`
	AgentUpdateStatus      *string `json:"agentUpdateStatus,omitempty"`
	AgentVersion           *string `json:"agentVersion,omitempty"`
	PostUpdateRebootStatus *string `json:"postUpdateRebootStatus,omitempty"`
}

func (o *InMageAgentDetails) GetAgentExpiryDateAsTime() (*time.Time, error) {
	if o.AgentExpiryDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.AgentExpiryDate, "2006-01-02T15:04:05Z07:00")
}

func (o *InMageAgentDetails) SetAgentExpiryDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.AgentExpiryDate = &formatted
}
