package links

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ValidateResult struct {
	AuthType              *AuthType               `json:"authType,omitempty"`
	IsConnectionAvailable *bool                   `json:"isConnectionAvailable,omitempty"`
	LinkerName            *string                 `json:"linkerName,omitempty"`
	ReportEndTimeUtc      *string                 `json:"reportEndTimeUtc,omitempty"`
	ReportStartTimeUtc    *string                 `json:"reportStartTimeUtc,omitempty"`
	SourceId              *string                 `json:"sourceId,omitempty"`
	TargetId              *string                 `json:"targetId,omitempty"`
	ValidationDetail      *[]ValidationResultItem `json:"validationDetail,omitempty"`
}

func (o *ValidateResult) GetReportEndTimeUtcAsTime() (*time.Time, error) {
	if o.ReportEndTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ReportEndTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ValidateResult) SetReportEndTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ReportEndTimeUtc = &formatted
}

func (o *ValidateResult) GetReportStartTimeUtcAsTime() (*time.Time, error) {
	if o.ReportStartTimeUtc == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.ReportStartTimeUtc, "2006-01-02T15:04:05Z07:00")
}

func (o *ValidateResult) SetReportStartTimeUtcAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.ReportStartTimeUtc = &formatted
}
