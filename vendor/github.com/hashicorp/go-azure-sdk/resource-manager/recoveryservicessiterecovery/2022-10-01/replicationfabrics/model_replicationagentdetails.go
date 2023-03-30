package replicationfabrics

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReplicationAgentDetails struct {
	BiosId           *string           `json:"biosId,omitempty"`
	FabricObjectId   *string           `json:"fabricObjectId,omitempty"`
	Fqdn             *string           `json:"fqdn,omitempty"`
	Health           *ProtectionHealth `json:"health,omitempty"`
	HealthErrors     *[]HealthError    `json:"healthErrors,omitempty"`
	Id               *string           `json:"id,omitempty"`
	LastHeartbeatUtc *string           `json:"lastHeartbeatUtc,omitempty"`
	Name             *string           `json:"name,omitempty"`
	Version          *string           `json:"version,omitempty"`
}

func (o *ReplicationAgentDetails) GetLastHeartbeatUtcAsTime() (*time.Time, error) {
	if o.LastHeartbeatUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastHeartbeatUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ReplicationAgentDetails) SetLastHeartbeatUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastHeartbeatUtc = &formatted
}
