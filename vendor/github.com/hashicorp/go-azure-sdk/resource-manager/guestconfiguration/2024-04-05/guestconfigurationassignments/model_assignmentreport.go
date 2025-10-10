package guestconfigurationassignments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssignmentReport struct {
	Assignment       *AssignmentInfo             `json:"assignment,omitempty"`
	ComplianceStatus *ComplianceStatus           `json:"complianceStatus,omitempty"`
	EndTime          *string                     `json:"endTime,omitempty"`
	Id               *string                     `json:"id,omitempty"`
	OperationType    *Type                       `json:"operationType,omitempty"`
	ReportId         *string                     `json:"reportId,omitempty"`
	Resources        *[]AssignmentReportResource `json:"resources,omitempty"`
	StartTime        *string                     `json:"startTime,omitempty"`
	VM               *VMInfo                     `json:"vm,omitempty"`
}

func (o *AssignmentReport) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AssignmentReport) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *AssignmentReport) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *AssignmentReport) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
