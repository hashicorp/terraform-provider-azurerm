package webapps

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContainerInfo struct {
	CurrentCPUStats   *ContainerCPUStatistics              `json:"currentCpuStats,omitempty"`
	CurrentTimeStamp  *string                              `json:"currentTimeStamp,omitempty"`
	Eth0              *ContainerNetworkInterfaceStatistics `json:"eth0,omitempty"`
	Id                *string                              `json:"id,omitempty"`
	MemoryStats       *ContainerMemoryStatistics           `json:"memoryStats,omitempty"`
	Name              *string                              `json:"name,omitempty"`
	PreviousCPUStats  *ContainerCPUStatistics              `json:"previousCpuStats,omitempty"`
	PreviousTimeStamp *string                              `json:"previousTimeStamp,omitempty"`
}

func (o *ContainerInfo) GetCurrentTimeStampAsTime() (*time.Time, error) {
	if o.CurrentTimeStamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CurrentTimeStamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ContainerInfo) SetCurrentTimeStampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CurrentTimeStamp = &formatted
}

func (o *ContainerInfo) GetPreviousTimeStampAsTime() (*time.Time, error) {
	if o.PreviousTimeStamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.PreviousTimeStamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ContainerInfo) SetPreviousTimeStampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.PreviousTimeStamp = &formatted
}
