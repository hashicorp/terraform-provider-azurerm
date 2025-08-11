package imagedefinitions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageDefinitionBuildTaskGroup struct {
	EndTime   *string                     `json:"endTime,omitempty"`
	Name      *string                     `json:"name,omitempty"`
	StartTime *string                     `json:"startTime,omitempty"`
	Status    *ImageDefinitionBuildStatus `json:"status,omitempty"`
	Tasks     *[]ImageDefinitionBuildTask `json:"tasks,omitempty"`
}

func (o *ImageDefinitionBuildTaskGroup) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ImageDefinitionBuildTaskGroup) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *ImageDefinitionBuildTaskGroup) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ImageDefinitionBuildTaskGroup) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
