package querypackqueries

import (
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/dates"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LogAnalyticsQueryPackQueryProperties struct {
	Author       *string                                      `json:"author,omitempty"`
	Body         string                                       `json:"body"`
	Description  *string                                      `json:"description,omitempty"`
	DisplayName  string                                       `json:"displayName"`
	Id           *string                                      `json:"id,omitempty"`
	Properties   *interface{}                                 `json:"properties,omitempty"`
	Related      *LogAnalyticsQueryPackQueryPropertiesRelated `json:"related,omitempty"`
	Tags         *map[string][]string                         `json:"tags,omitempty"`
	TimeCreated  *string                                      `json:"timeCreated,omitempty"`
	TimeModified *string                                      `json:"timeModified,omitempty"`
}

func (o *LogAnalyticsQueryPackQueryProperties) GetTimeCreatedAsTime() (*time.Time, error) {
	if o.TimeCreated == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeCreated, "2006-01-02T15:04:05Z07:00")
}

func (o *LogAnalyticsQueryPackQueryProperties) SetTimeCreatedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeCreated = &formatted
}

func (o *LogAnalyticsQueryPackQueryProperties) GetTimeModifiedAsTime() (*time.Time, error) {
	if o.TimeModified == nil {
		return nil, nil
	}
	return dates.ParseAsFormat(o.TimeModified, "2006-01-02T15:04:05Z07:00")
}

func (o *LogAnalyticsQueryPackQueryProperties) SetTimeModifiedAsTime(input time.Time) {
	formatted := input.Format("2006-01-02T15:04:05Z07:00")
	o.TimeModified = &formatted
}
