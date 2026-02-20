package guestconfigurationhcrpassignments

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GuestConfigurationAssignmentReportProperties struct {
	Assignment       *AssignmentInfo          `json:"assignment,omitempty"`
	ComplianceStatus *ComplianceStatus        `json:"complianceStatus,omitempty"`
	Details          *AssignmentReportDetails `json:"details,omitempty"`
	EndTime          *string                  `json:"endTime,omitempty"`
	ReportId         *string                  `json:"reportId,omitempty"`
	StartTime        *string                  `json:"startTime,omitempty"`
	VM               *VMInfo                  `json:"vm,omitempty"`
	VMSSResourceId   *string                  `json:"vmssResourceId,omitempty"`
}

func (o *GuestConfigurationAssignmentReportProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *GuestConfigurationAssignmentReportProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *GuestConfigurationAssignmentReportProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *GuestConfigurationAssignmentReportProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
