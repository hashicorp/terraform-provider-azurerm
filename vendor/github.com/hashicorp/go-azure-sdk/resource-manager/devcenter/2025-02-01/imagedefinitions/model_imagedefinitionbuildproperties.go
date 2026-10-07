package imagedefinitions

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ImageDefinitionBuildProperties struct {
	EndTime        *string                     `json:"endTime,omitempty"`
	ErrorDetails   *ImageCreationErrorDetails  `json:"errorDetails,omitempty"`
	ImageReference *ImageReference             `json:"imageReference,omitempty"`
	StartTime      *string                     `json:"startTime,omitempty"`
	Status         *ImageDefinitionBuildStatus `json:"status,omitempty"`
}

func (o *ImageDefinitionBuildProperties) GetEndTimeAsTime() (*time.Time, error) {
	if o.EndTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.EndTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ImageDefinitionBuildProperties) SetEndTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.EndTime = &formatted
}

func (o *ImageDefinitionBuildProperties) GetStartTimeAsTime() (*time.Time, error) {
	if o.StartTime == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.StartTime, "2006-01-02T15:04:05Z07:00")
}

func (o *ImageDefinitionBuildProperties) SetStartTimeAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.StartTime = &formatted
}
