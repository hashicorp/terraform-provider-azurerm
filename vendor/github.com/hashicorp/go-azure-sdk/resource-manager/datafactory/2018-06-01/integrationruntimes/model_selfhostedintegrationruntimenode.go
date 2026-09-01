package integrationruntimes

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SelfHostedIntegrationRuntimeNode struct {
	Capabilities        *map[string]string                      `json:"capabilities,omitempty"`
	ConcurrentJobsLimit *int64                                  `json:"concurrentJobsLimit,omitempty"`
	ExpiryTime          *string                                 `json:"expiryTime,omitempty"`
	HostServiceUri      *string                                 `json:"hostServiceUri,omitempty"`
	IsActiveDispatcher  *bool                                   `json:"isActiveDispatcher,omitempty"`
	LastConnectTime     *string                                 `json:"lastConnectTime,omitempty"`
	LastEndUpdateTime   *string                                 `json:"lastEndUpdateTime,omitempty"`
	LastStartTime       *string                                 `json:"lastStartTime,omitempty"`
	LastStartUpdateTime *string                                 `json:"lastStartUpdateTime,omitempty"`
	LastStopTime        *string                                 `json:"lastStopTime,omitempty"`
	LastUpdateResult    *IntegrationRuntimeUpdateResult         `json:"lastUpdateResult,omitempty"`
	MachineName         *string                                 `json:"machineName,omitempty"`
	MaxConcurrentJobs   *int64                                  `json:"maxConcurrentJobs,omitempty"`
	NodeName            *string                                 `json:"nodeName,omitempty"`
	RegisterTime        *string                                 `json:"registerTime,omitempty"`
	Status              *SelfHostedIntegrationRuntimeNodeStatus `json:"status,omitempty"`
	Version             *string                                 `json:"version,omitempty"`
	VersionStatus       *string                                 `json:"versionStatus,omitempty"`
}

func (o *SelfHostedIntegrationRuntimeNode) GetExpiryTimeAsTime() (*time.Time, error) {
	if o.ExpiryTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ExpiryTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SelfHostedIntegrationRuntimeNode) SetExpiryTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ExpiryTime = &formatted
}

func (o *SelfHostedIntegrationRuntimeNode) GetLastConnectTimeAsTime() (*time.Time, error) {
	if o.LastConnectTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastConnectTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SelfHostedIntegrationRuntimeNode) SetLastConnectTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastConnectTime = &formatted
}

func (o *SelfHostedIntegrationRuntimeNode) GetLastEndUpdateTimeAsTime() (*time.Time, error) {
	if o.LastEndUpdateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastEndUpdateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SelfHostedIntegrationRuntimeNode) SetLastEndUpdateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastEndUpdateTime = &formatted
}

func (o *SelfHostedIntegrationRuntimeNode) GetLastStartTimeAsTime() (*time.Time, error) {
	if o.LastStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SelfHostedIntegrationRuntimeNode) SetLastStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStartTime = &formatted
}

func (o *SelfHostedIntegrationRuntimeNode) GetLastStartUpdateTimeAsTime() (*time.Time, error) {
	if o.LastStartUpdateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastStartUpdateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SelfHostedIntegrationRuntimeNode) SetLastStartUpdateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStartUpdateTime = &formatted
}

func (o *SelfHostedIntegrationRuntimeNode) GetLastStopTimeAsTime() (*time.Time, error) {
	if o.LastStopTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastStopTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SelfHostedIntegrationRuntimeNode) SetLastStopTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStopTime = &formatted
}

func (o *SelfHostedIntegrationRuntimeNode) GetRegisterTimeAsTime() (*time.Time, error) {
	if o.RegisterTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.RegisterTime, "2006-01-02T15:04:05Z07:00")
}

func (o *SelfHostedIntegrationRuntimeNode) SetRegisterTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.RegisterTime = &formatted
}
