package webapps

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProcessInfoProperties struct {
	Children                   *[]string            `json:"children,omitempty"`
	CommandLine                *string              `json:"command_line,omitempty"`
	DeploymentName             *string              `json:"deployment_name,omitempty"`
	Description                *string              `json:"description,omitempty"`
	EnvironmentVariables       *map[string]string   `json:"environment_variables,omitempty"`
	FileName                   *string              `json:"file_name,omitempty"`
	HandleCount                *int64               `json:"handle_count,omitempty"`
	Href                       *string              `json:"href,omitempty"`
	Identifier                 *int64               `json:"identifier,omitempty"`
	IisProfileTimeoutInSeconds *float64             `json:"iis_profile_timeout_in_seconds,omitempty"`
	IsIisProfileRunning        *bool                `json:"is_iis_profile_running,omitempty"`
	IsProfileRunning           *bool                `json:"is_profile_running,omitempty"`
	IsScmSite                  *bool                `json:"is_scm_site,omitempty"`
	IsWebjob                   *bool                `json:"is_webjob,omitempty"`
	Minidump                   *string              `json:"minidump,omitempty"`
	ModuleCount                *int64               `json:"module_count,omitempty"`
	Modules                    *[]ProcessModuleInfo `json:"modules,omitempty"`
	NonPagedSystemMemory       *int64               `json:"non_paged_system_memory,omitempty"`
	OpenFileHandles            *[]string            `json:"open_file_handles,omitempty"`
	PagedMemory                *int64               `json:"paged_memory,omitempty"`
	PagedSystemMemory          *int64               `json:"paged_system_memory,omitempty"`
	Parent                     *string              `json:"parent,omitempty"`
	PeakPagedMemory            *int64               `json:"peak_paged_memory,omitempty"`
	PeakVirtualMemory          *int64               `json:"peak_virtual_memory,omitempty"`
	PeakWorkingSet             *int64               `json:"peak_working_set,omitempty"`
	PrivateMemory              *int64               `json:"private_memory,omitempty"`
	PrivilegedCpuTime          *string              `json:"privileged_cpu_time,omitempty"`
	StartTime                  *string              `json:"start_time,omitempty"`
	ThreadCount                *int64               `json:"thread_count,omitempty"`
	Threads                    *[]ProcessThreadInfo `json:"threads,omitempty"`
	TimeStamp                  *string              `json:"time_stamp,omitempty"`
	TotalCpuTime               *string              `json:"total_cpu_time,omitempty"`
	UserCpuTime                *string              `json:"user_cpu_time,omitempty"`
	UserName                   *string              `json:"user_name,omitempty"`
	VirtualMemory              *int64               `json:"virtual_memory,omitempty"`
	WorkingSet                 *int64               `json:"working_set,omitempty"`
}

func (o *ProcessInfoProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ProcessInfoProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}

func (o *ProcessInfoProperties) GetTimeStampAsTime() (*time.Time, error) {
	if o.TimeStamp == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeStamp, "2006-01-02T15:04:05Z07:00")
}

func (o *ProcessInfoProperties) SetTimeStampAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeStamp = &formatted
}
