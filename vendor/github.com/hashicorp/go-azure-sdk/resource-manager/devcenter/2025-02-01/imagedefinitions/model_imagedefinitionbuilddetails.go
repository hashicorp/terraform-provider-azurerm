package imagedefinitions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/systemdata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageDefinitionBuildDetails struct {
	EndTime        *string                          `json:"endTime,omitempty"`
	ErrorDetails   *ImageCreationErrorDetails       `json:"errorDetails,omitempty"`
	Id             *string                          `json:"id,omitempty"`
	ImageReference *ImageReference                  `json:"imageReference,omitempty"`
	Name           *string                          `json:"name,omitempty"`
	StartTime      *string                          `json:"startTime,omitempty"`
	Status         *ImageDefinitionBuildStatus      `json:"status,omitempty"`
	SystemData     *systemdata.SystemData           `json:"systemData,omitempty"`
	TaskGroups     *[]ImageDefinitionBuildTaskGroup `json:"taskGroups,omitempty"`
	Type           *string                          `json:"type,omitempty"`
}

func (o *ImageDefinitionBuildDetails) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ImageDefinitionBuildDetails) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *ImageDefinitionBuildDetails) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ImageDefinitionBuildDetails) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
