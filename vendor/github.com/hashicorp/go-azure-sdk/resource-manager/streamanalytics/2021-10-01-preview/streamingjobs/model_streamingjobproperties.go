package streamingjobs

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type StreamingJobProperties struct {
	Cluster                            *ClusterInfo            `json:"cluster,omitempty"`
	CompatibilityLevel                 *CompatibilityLevel     `json:"compatibilityLevel,omitempty"`
	ContentStoragePolicy               *ContentStoragePolicy   `json:"contentStoragePolicy,omitempty"`
	CreatedDate                        *string                 `json:"createdDate,omitempty"`
	DataLocale                         *string                 `json:"dataLocale,omitempty"`
	Etag                               *string                 `json:"etag,omitempty"`
	EventsLateArrivalMaxDelayInSeconds *int64                  `json:"eventsLateArrivalMaxDelayInSeconds,omitempty"`
	EventsOutOfOrderMaxDelayInSeconds  *int64                  `json:"eventsOutOfOrderMaxDelayInSeconds,omitempty"`
	EventsOutOfOrderPolicy             *EventsOutOfOrderPolicy `json:"eventsOutOfOrderPolicy,omitempty"`
	Externals                          *External               `json:"externals,omitempty"`
	Functions                          *[]Function             `json:"functions,omitempty"`
	Inputs                             *[]Input                `json:"inputs,omitempty"`
	JobId                              *string                 `json:"jobId,omitempty"`
	JobState                           *string                 `json:"jobState,omitempty"`
	JobStorageAccount                  *JobStorageAccount      `json:"jobStorageAccount,omitempty"`
	JobType                            *JobType                `json:"jobType,omitempty"`
	LastOutputEventTime                *string                 `json:"lastOutputEventTime,omitempty"`
	OutputErrorPolicy                  *OutputErrorPolicy      `json:"outputErrorPolicy,omitempty"`
	OutputStartMode                    *OutputStartMode        `json:"outputStartMode,omitempty"`
	OutputStartTime                    *string                 `json:"outputStartTime,omitempty"`
	Outputs                            *[]Output               `json:"outputs,omitempty"`
	ProvisioningState                  *string                 `json:"provisioningState,omitempty"`
	Sku                                *Sku                    `json:"sku,omitempty"`
	Transformation                     *Transformation         `json:"transformation,omitempty"`
}

func (o *StreamingJobProperties) GetCreatedDateAsTime() (*time.Time, error) {
	if o.CreatedDate == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.CreatedDate, "2006-01-02T15:04:05Z07:00")
}

func (o *StreamingJobProperties) SetCreatedDateAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.CreatedDate = &formatted
}

func (o *StreamingJobProperties) GetLastOutputEventTimeAsTime() (*time.Time, error) {
	if o.LastOutputEventTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.LastOutputEventTime, "2006-01-02T15:04:05Z07:00")
}

func (o *StreamingJobProperties) SetLastOutputEventTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.LastOutputEventTime = &formatted
}

func (o *StreamingJobProperties) GetOutputStartTimeAsTime() (*time.Time, error) {
	if o.OutputStartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.OutputStartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *StreamingJobProperties) SetOutputStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.OutputStartTime = &formatted
}
