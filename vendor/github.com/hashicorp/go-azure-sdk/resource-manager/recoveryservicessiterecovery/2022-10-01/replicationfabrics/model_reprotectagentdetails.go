package replicationfabrics

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ReprotectAgentDetails struct {
	AccessibleDatastores *[]string         `json:"accessibleDatastores,omitempty"`
	BiosId               *string           `json:"biosId,omitempty"`
	FabricObjectId       *string           `json:"fabricObjectId,omitempty"`
	Fqdn                 *string           `json:"fqdn,omitempty"`
	Health               *ProtectionHealth `json:"health,omitempty"`
	HealthErrors         *[]HealthError    `json:"healthErrors,omitempty"`
	Id                   *string           `json:"id,omitempty"`
	LastDiscoveryInUtc   *string           `json:"lastDiscoveryInUtc,omitempty"`
	LastHeartbeatUtc     *string           `json:"lastHeartbeatUtc,omitempty"`
	Name                 *string           `json:"name,omitempty"`
	ProtectedItemCount   *int64            `json:"protectedItemCount,omitempty"`
	VcenterId            *string           `json:"vcenterId,omitempty"`
	Version              *string           `json:"version,omitempty"`
}

func (o *ReprotectAgentDetails) GetLastDiscoveryInUtcAsTime() (*time.Time, error) {
	if o.LastDiscoveryInUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastDiscoveryInUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ReprotectAgentDetails) SetLastDiscoveryInUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastDiscoveryInUtc = &formatted
}

func (o *ReprotectAgentDetails) GetLastHeartbeatUtcAsTime() (*time.Time, error) {
	if o.LastHeartbeatUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastHeartbeatUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ReprotectAgentDetails) SetLastHeartbeatUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastHeartbeatUtc = &formatted
}
