package exports

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ExportRunProperties struct {
	EndDate             *string                 `json:"endDate,omitempty"`
	Error               *ErrorDetails           `json:"error,omitempty"`
	ExecutionType       *ExecutionType          `json:"executionType,omitempty"`
	FileName            *string                 `json:"fileName,omitempty"`
	ManifestFile        *string                 `json:"manifestFile,omitempty"`
	ProcessingEndTime   *string                 `json:"processingEndTime,omitempty"`
	ProcessingStartTime *string                 `json:"processingStartTime,omitempty"`
	RunSettings         *CommonExportProperties `json:"runSettings,omitempty"`
	StartDate           *string                 `json:"startDate,omitempty"`
	Status              *ExecutionStatus        `json:"status,omitempty"`
	SubmittedBy         *string                 `json:"submittedBy,omitempty"`
	SubmittedTime       *string                 `json:"submittedTime,omitempty"`
}

func (o *ExportRunProperties) GetEndDateAsTime() (*time.Time, error) {
	if o.EndDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ExportRunProperties) SetEndDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndDate = &formatted
}

func (o *ExportRunProperties) GetProcessingEndTimeAsTime() (*time.Time, error) {
	if o.ProcessingEndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ProcessingEndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ExportRunProperties) SetProcessingEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ProcessingEndTime = &formatted
}

func (o *ExportRunProperties) GetProcessingStartTimeAsTime() (*time.Time, error) {
	if o.ProcessingStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ProcessingStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ExportRunProperties) SetProcessingStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ProcessingStartTime = &formatted
}

func (o *ExportRunProperties) GetStartDateAsTime() (*time.Time, error) {
	if o.StartDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartDate, "2006-01-02T15:04:05Z07:00")
}

func (o *ExportRunProperties) SetStartDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartDate = &formatted
}

func (o *ExportRunProperties) GetSubmittedTimeAsTime() (*time.Time, error) {
	if o.SubmittedTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.SubmittedTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ExportRunProperties) SetSubmittedTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.SubmittedTime = &formatted
}
