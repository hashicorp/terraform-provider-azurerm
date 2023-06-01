package agents

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AgentProperties struct {
	AgentStatus       *AgentStatus                 `json:"agentStatus,omitempty"`
	AgentVersion      *string                      `json:"agentVersion,omitempty"`
	ArcResourceId     string                       `json:"arcResourceId"`
	ArcVMUuid         string                       `json:"arcVmUuid"`
	Description       *string                      `json:"description,omitempty"`
	ErrorDetails      *AgentPropertiesErrorDetails `json:"errorDetails,omitempty"`
	LastStatusUpdate  *string                      `json:"lastStatusUpdate,omitempty"`
	LocalIPAddress    *string                      `json:"localIPAddress,omitempty"`
	MemoryInMB        *int64                       `json:"memoryInMB,omitempty"`
	NumberOfCores     *int64                       `json:"numberOfCores,omitempty"`
	ProvisioningState *ProvisioningState           `json:"provisioningState,omitempty"`
	UptimeInSeconds   *int64                       `json:"uptimeInSeconds,omitempty"`
}

func (o *AgentProperties) GetLastStatusUpdateAsTime() (*time.Time, error) {
	if o.LastStatusUpdate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastStatusUpdate, "2006-01-02T15:04:05Z07:00")
}

func (o *AgentProperties) SetLastStatusUpdateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStatusUpdate = &formatted
}
