package replicationprotecteditems

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InMageRcmFailbackMobilityAgentDetails struct {
	AgentVersionExpiryDate               *string                      `json:"agentVersionExpiryDate,omitempty"`
	DriverVersion                        *string                      `json:"driverVersion,omitempty"`
	DriverVersionExpiryDate              *string                      `json:"driverVersionExpiryDate,omitempty"`
	IsUpgradeable                        *string                      `json:"isUpgradeable,omitempty"`
	LastHeartbeatUtc                     *string                      `json:"lastHeartbeatUtc,omitempty"`
	LatestUpgradableVersionWithoutReboot *string                      `json:"latestUpgradableVersionWithoutReboot,omitempty"`
	LatestVersion                        *string                      `json:"latestVersion,omitempty"`
	ReasonsBlockingUpgrade               *[]AgentUpgradeBlockedReason `json:"reasonsBlockingUpgrade,omitempty"`
	Version                              *string                      `json:"version,omitempty"`
}

func (o *InMageRcmFailbackMobilityAgentDetails) GetAgentVersionExpiryDateAsTime() (*time.Time, error) {
	if o.AgentVersionExpiryDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.AgentVersionExpiryDate, "2006-01-02T15:04:05Z07:00")
}

func (o *InMageRcmFailbackMobilityAgentDetails) SetAgentVersionExpiryDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.AgentVersionExpiryDate = &formatted
}

func (o *InMageRcmFailbackMobilityAgentDetails) GetDriverVersionExpiryDateAsTime() (*time.Time, error) {
	if o.DriverVersionExpiryDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.DriverVersionExpiryDate, "2006-01-02T15:04:05Z07:00")
}

func (o *InMageRcmFailbackMobilityAgentDetails) SetDriverVersionExpiryDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.DriverVersionExpiryDate = &formatted
}

func (o *InMageRcmFailbackMobilityAgentDetails) GetLastHeartbeatUtcAsTime() (*time.Time, error) {
	if o.LastHeartbeatUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastHeartbeatUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *InMageRcmFailbackMobilityAgentDetails) SetLastHeartbeatUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastHeartbeatUtc = &formatted
}
