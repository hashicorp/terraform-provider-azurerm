package pipelineruns

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PipelineRunResponse struct {
	CatalogDigest           *string                         `json:"catalogDigest,omitempty"`
	FinishTime              *string                         `json:"finishTime,omitempty"`
	ImportedArtifacts       *[]string                       `json:"importedArtifacts,omitempty"`
	PipelineRunErrorMessage *string                         `json:"pipelineRunErrorMessage,omitempty"`
	Progress                *ProgressProperties             `json:"progress,omitempty"`
	Source                  *ImportPipelineSourceProperties `json:"source,omitempty"`
	StartTime               *string                         `json:"startTime,omitempty"`
	Status                  *string                         `json:"status,omitempty"`
	Target                  *ExportPipelineTargetProperties `json:"target,omitempty"`
	Trigger                 *PipelineTriggerDescriptor      `json:"trigger,omitempty"`
}

func (o *PipelineRunResponse) GetFinishTimeAsTime() (*time.Time, error) {
	if o.FinishTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.FinishTime, "2006-01-02T15:04:05Z07:00")
}

func (o *PipelineRunResponse) SetFinishTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.FinishTime = &formatted
}

func (o *PipelineRunResponse) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *PipelineRunResponse) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
