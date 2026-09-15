package indexers

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type IndexerExecutionResult struct {
	EndTime              *string                `json:"endTime,omitempty"`
	ErrorMessage         *string                `json:"errorMessage,omitempty"`
	Errors               []SearchIndexerError   `json:"errors"`
	FinalTrackingState   *string                `json:"finalTrackingState,omitempty"`
	InitialTrackingState *string                `json:"initialTrackingState,omitempty"`
	ItemsFailed          int64                  `json:"itemsFailed"`
	ItemsProcessed       int64                  `json:"itemsProcessed"`
	StartTime            *string                `json:"startTime,omitempty"`
	Status               IndexerExecutionStatus `json:"status"`
	Warnings             []SearchIndexerWarning `json:"warnings"`
}

func (o *IndexerExecutionResult) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *IndexerExecutionResult) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *IndexerExecutionResult) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *IndexerExecutionResult) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
