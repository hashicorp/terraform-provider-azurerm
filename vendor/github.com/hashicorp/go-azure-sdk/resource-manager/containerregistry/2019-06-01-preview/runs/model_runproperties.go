package runs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RunProperties struct {
	AgentConfiguration *AgentProperties         `json:"agentConfiguration,omitempty"`
	AgentPoolName      *string                  `json:"agentPoolName,omitempty"`
	CreateTime         *string                  `json:"createTime,omitempty"`
	CustomRegistries   *[]string                `json:"customRegistries,omitempty"`
	FinishTime         *string                  `json:"finishTime,omitempty"`
	ImageUpdateTrigger *ImageUpdateTrigger      `json:"imageUpdateTrigger,omitempty"`
	IsArchiveEnabled   *bool                    `json:"isArchiveEnabled,omitempty"`
	LastUpdatedTime    *string                  `json:"lastUpdatedTime,omitempty"`
	LogArtifact        *ImageDescriptor         `json:"logArtifact,omitempty"`
	OutputImages       *[]ImageDescriptor       `json:"outputImages,omitempty"`
	Platform           *PlatformProperties      `json:"platform,omitempty"`
	ProvisioningState  *ProvisioningState       `json:"provisioningState,omitempty"`
	RunErrorMessage    *string                  `json:"runErrorMessage,omitempty"`
	RunId              *string                  `json:"runId,omitempty"`
	RunType            *RunType                 `json:"runType,omitempty"`
	SourceRegistryAuth *string                  `json:"sourceRegistryAuth,omitempty"`
	SourceTrigger      *SourceTriggerDescriptor `json:"sourceTrigger,omitempty"`
	StartTime          *string                  `json:"startTime,omitempty"`
	Status             *RunStatus               `json:"status,omitempty"`
	Task               *string                  `json:"task,omitempty"`
	TimerTrigger       *TimerTriggerDescriptor  `json:"timerTrigger,omitempty"`
	UpdateTriggerToken *string                  `json:"updateTriggerToken,omitempty"`
}

func (o *RunProperties) GetCreateTimeAsTime() (*time.Time, error) {
	if o.CreateTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreateTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RunProperties) SetCreateTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreateTime = &formatted
}

func (o *RunProperties) GetFinishTimeAsTime() (*time.Time, error) {
	if o.FinishTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.FinishTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RunProperties) SetFinishTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.FinishTime = &formatted
}

func (o *RunProperties) GetLastUpdatedTimeAsTime() (*time.Time, error) {
	if o.LastUpdatedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastUpdatedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RunProperties) SetLastUpdatedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastUpdatedTime = &formatted
}

func (o *RunProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *RunProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
