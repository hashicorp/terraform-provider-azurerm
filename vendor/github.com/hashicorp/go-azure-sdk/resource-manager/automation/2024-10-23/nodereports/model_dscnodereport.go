package nodereports

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DscNodeReport struct {
	ConfigurationVersion *string               `json:"configurationVersion,omitempty"`
	EndTime              *string               `json:"endTime,omitempty"`
	Errors               *[]DscReportError     `json:"errors,omitempty"`
	HostName             *string               `json:"hostName,omitempty"`
	IPV4Addresses        *[]string             `json:"iPV4Addresses,omitempty"`
	IPV6Addresses        *[]string             `json:"iPV6Addresses,omitempty"`
	Id                   *string               `json:"id,omitempty"`
	LastModifiedTime     *string               `json:"lastModifiedTime,omitempty"`
	MetaConfiguration    *DscMetaConfiguration `json:"metaConfiguration,omitempty"`
	NumberOfResources    *int64                `json:"numberOfResources,omitempty"`
	RawErrors            *string               `json:"rawErrors,omitempty"`
	RebootRequested      *string               `json:"rebootRequested,omitempty"`
	RefreshMode          *string               `json:"refreshMode,omitempty"`
	ReportFormatVersion  *string               `json:"reportFormatVersion,omitempty"`
	ReportId             *string               `json:"reportId,omitempty"`
	Resources            *[]DscReportResource  `json:"resources,omitempty"`
	StartTime            *string               `json:"startTime,omitempty"`
	Status               *string               `json:"status,omitempty"`
	Type                 *string               `json:"type,omitempty"`
}

func (o *DscNodeReport) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *DscNodeReport) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *DscNodeReport) GetLastModifiedTimeAsTime() (*time.Time, error) {
	if o.LastModifiedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastModifiedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *DscNodeReport) SetLastModifiedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastModifiedTime = &formatted
}

func (o *DscNodeReport) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *DscNodeReport) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
