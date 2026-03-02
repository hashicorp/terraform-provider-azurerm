package imagedefinitions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageDefinitionBuildTask struct {
	DisplayName *string                                      `json:"displayName,omitempty"`
	EndTime     *string                                      `json:"endTime,omitempty"`
	Id          *string                                      `json:"id,omitempty"`
	LogUri      *string                                      `json:"logUri,omitempty"`
	Name        *string                                      `json:"name,omitempty"`
	Parameters  *[]ImageDefinitionBuildTaskParametersInlined `json:"parameters,omitempty"`
	StartTime   *string                                      `json:"startTime,omitempty"`
	Status      *ImageDefinitionBuildStatus                  `json:"status,omitempty"`
}

func (o *ImageDefinitionBuildTask) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ImageDefinitionBuildTask) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *ImageDefinitionBuildTask) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ImageDefinitionBuildTask) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
