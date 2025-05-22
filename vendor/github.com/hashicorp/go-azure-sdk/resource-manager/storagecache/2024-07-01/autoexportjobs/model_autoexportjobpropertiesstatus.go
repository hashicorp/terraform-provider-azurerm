package autoexportjobs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AutoExportJobPropertiesStatus struct {
	CurrentIterationFilesDiscovered          *int64                `json:"currentIterationFilesDiscovered,omitempty"`
	CurrentIterationFilesExported            *int64                `json:"currentIterationFilesExported,omitempty"`
	CurrentIterationFilesFailed              *int64                `json:"currentIterationFilesFailed,omitempty"`
	CurrentIterationMiBDiscovered            *int64                `json:"currentIterationMiBDiscovered,omitempty"`
	CurrentIterationMiBExported              *int64                `json:"currentIterationMiBExported,omitempty"`
	ExportIterationCount                     *int64                `json:"exportIterationCount,omitempty"`
	LastCompletionTimeUTC                    *string               `json:"lastCompletionTimeUTC,omitempty"`
	LastStartedTimeUTC                       *string               `json:"lastStartedTimeUTC,omitempty"`
	LastSuccessfulIterationCompletionTimeUTC *string               `json:"lastSuccessfulIterationCompletionTimeUTC,omitempty"`
	State                                    *AutoExportStatusType `json:"state,omitempty"`
	StatusCode                               *string               `json:"statusCode,omitempty"`
	StatusMessage                            *string               `json:"statusMessage,omitempty"`
	TotalFilesExported                       *int64                `json:"totalFilesExported,omitempty"`
	TotalFilesFailed                         *int64                `json:"totalFilesFailed,omitempty"`
	TotalMiBExported                         *int64                `json:"totalMiBExported,omitempty"`
}

func (o *AutoExportJobPropertiesStatus) GetLastCompletionTimeUTCAsTime() (*time.Time, error) {
	if o.LastCompletionTimeUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastCompletionTimeUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *AutoExportJobPropertiesStatus) SetLastCompletionTimeUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastCompletionTimeUTC = &formatted
}

func (o *AutoExportJobPropertiesStatus) GetLastStartedTimeUTCAsTime() (*time.Time, error) {
	if o.LastStartedTimeUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastStartedTimeUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *AutoExportJobPropertiesStatus) SetLastStartedTimeUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastStartedTimeUTC = &formatted
}

func (o *AutoExportJobPropertiesStatus) GetLastSuccessfulIterationCompletionTimeUTCAsTime() (*time.Time, error) {
	if o.LastSuccessfulIterationCompletionTimeUTC == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastSuccessfulIterationCompletionTimeUTC, "2006-01-02T15:04:05Z07:00")
}

func (o *AutoExportJobPropertiesStatus) SetLastSuccessfulIterationCompletionTimeUTCAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastSuccessfulIterationCompletionTimeUTC = &formatted
}
